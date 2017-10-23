package local

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/agent/token"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/consul/types"
)

// permissionDenied is returned when an ACL based rejection happens.
const permissionDenied = "Permission denied"

// syncStatus is used to represent the difference between
// the local and remote state, and if action needs to be taken
type syncStatus struct {
	inSync bool // Is this in sync with the server
}

// Config is the configuration for the State. It is
// populated during NewLocalAgent from the agent configuration to avoid
// race conditions with the agent configuration.
type Config struct {
	AdvertiseAddr       string
	CheckUpdateInterval time.Duration
	Datacenter          string
	DiscardCheckOutput  bool
	NodeID              types.NodeID
	NodeName            string
	TaggedAddresses     map[string]string
}

type delegate interface {
	RPC(method string, args interface{}, reply interface{}) error
}

// State is used to represent the node's services,
// and checks. We used it to perform anti-entropy with the
// catalog representation
type State struct {
	sync.RWMutex
	logger *log.Logger

	// Config is the agent config
	config Config

	// delegate is the consul interface to use for keeping in sync
	delegate delegate

	// nodeInfoInSync tracks whether the server has our correct top-level
	// node information in sync
	nodeInfoInSync bool

	// Services tracks the local services
	services      map[string]*structs.NodeService
	serviceStatus map[string]syncStatus
	serviceTokens map[string]string

	// Checks tracks the local checks
	checks            map[types.CheckID]*structs.HealthCheck
	checkStatus       map[types.CheckID]syncStatus
	checkTokens       map[types.CheckID]string
	checkCriticalTime map[types.CheckID]time.Time

	// Used to track checks that are being deferred
	deferCheck map[types.CheckID]*time.Timer

	// metadata tracks the local metadata fields
	metadata map[string]string

	// triggerCh is used to inform of a change to local state
	// that requires anti-entropy with the server
	triggerCh chan struct{}

	// discardCheckOutput stores whether the output of health checks
	// is stored in the raft log.
	discardCheckOutput atomic.Value // bool

	tokens *token.Store
}

// NewLocalState creates a  is used to initialize the local state
func NewState(c Config, lg *log.Logger, tokens *token.Store, triggerCh chan struct{}) *State {
	l := &State{
		config:            c,
		logger:            lg,
		services:          make(map[string]*structs.NodeService),
		serviceStatus:     make(map[string]syncStatus),
		serviceTokens:     make(map[string]string),
		checks:            make(map[types.CheckID]*structs.HealthCheck),
		checkStatus:       make(map[types.CheckID]syncStatus),
		checkTokens:       make(map[types.CheckID]string),
		checkCriticalTime: make(map[types.CheckID]time.Time),
		deferCheck:        make(map[types.CheckID]*time.Timer),
		metadata:          make(map[string]string),
		triggerCh:         triggerCh,
		tokens:            tokens,
	}
	l.discardCheckOutput.Store(c.DiscardCheckOutput)
	return l
}

func (l *State) SetDelegate(d delegate) {
	l.delegate = d
}

// changeMade is used to trigger an anti-entropy run
func (l *State) changeMade() {
	// todo(fs): IMO, the non-blocking nature of this call should be hidden in the syncer
	select {
	case l.triggerCh <- struct{}{}:
	default:
	}
}

func (l *State) SetDiscardCheckOutput(b bool) {
	l.discardCheckOutput.Store(b)
}

// ServiceToken returns the configured ACL token for the given
// service ID. If none is present, the agent's token is returned.
func (l *State) ServiceToken(id string) string {
	l.RLock()
	defer l.RUnlock()
	return l.serviceToken(id)
}

// serviceToken returns an ACL token associated with a service.
func (l *State) serviceToken(id string) string {
	token := l.serviceTokens[id]
	if token == "" {
		token = l.tokens.UserToken()
	}
	return token
}

// AddService is used to add a service entry to the local state.
// This entry is persistent and the agent will make a best effort to
// ensure it is registered
func (l *State) AddService(service *structs.NodeService, token string) {
	// Assign the ID if none given
	if service.ID == "" && service.Service != "" {
		service.ID = service.Service
	}

	l.Lock()
	defer l.Unlock()

	l.services[service.ID] = service
	l.serviceStatus[service.ID] = syncStatus{}
	l.serviceTokens[service.ID] = token
	l.changeMade()
}

// RemoveService is used to remove a service entry from the local state.
// The agent will make a best effort to ensure it is deregistered
func (l *State) RemoveService(serviceID string) error {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.services[serviceID]; ok {
		delete(l.services, serviceID)
		// Leave the service token around, if any, until we successfully
		// delete the service.
		l.serviceStatus[serviceID] = syncStatus{inSync: false}
		l.changeMade()
	} else {
		return fmt.Errorf("Service does not exist")
	}

	return nil
}

// Service returns the locally registered service that the
// agent is aware of and are being kept in sync with the server
func (l *State) Service(id string) *structs.NodeService {
	l.RLock()
	defer l.RUnlock()
	return l.services[id]
}

// Services returns the locally registered services that the
// agent is aware of and are being kept in sync with the server
func (l *State) Services() map[string]*structs.NodeService {
	services := make(map[string]*structs.NodeService)
	l.RLock()
	defer l.RUnlock()

	for name, serv := range l.services {
		services[name] = serv
	}
	return services
}

// CheckToken is used to return the configured health check token for a
// Check, or if none is configured, the default agent ACL token.
func (l *State) CheckToken(checkID types.CheckID) string {
	l.RLock()
	defer l.RUnlock()
	return l.checkToken(checkID)
}

// checkToken returns an ACL token associated with a check.
func (l *State) checkToken(checkID types.CheckID) string {
	token := l.checkTokens[checkID]
	if token == "" {
		token = l.tokens.UserToken()
	}
	return token
}

// AddCheck is used to add a health check to the local state.
// This entry is persistent and the agent will make a best effort to
// ensure it is registered
func (l *State) AddCheck(check *structs.HealthCheck, token string) error {
	l.Lock()
	defer l.Unlock()

	// Set the node name
	check.Node = l.config.NodeName

	if l.discardCheckOutput.Load().(bool) {
		check.Output = ""
	}

	// if there is a serviceID associated with the check, make sure it exists before adding it
	// NOTE - This logic may be moved to be handled within the Agent's Addcheck method after a refactor
	if check.ServiceID != "" && l.services[check.ServiceID] == nil {
		return fmt.Errorf("ServiceID %q does not exist", check.ServiceID)
	}

	l.checks[check.CheckID] = check
	l.checkStatus[check.CheckID] = syncStatus{}
	l.checkTokens[check.CheckID] = token
	delete(l.checkCriticalTime, check.CheckID)
	l.changeMade()
	return nil
}

// RemoveCheck is used to remove a health check from the local state.
// The agent will make a best effort to ensure it is deregistered
func (l *State) RemoveCheck(checkID types.CheckID) {
	l.Lock()
	defer l.Unlock()

	delete(l.checks, checkID)
	// Leave the check token around, if any, until we successfully delete
	// the check.
	delete(l.checkCriticalTime, checkID)
	l.checkStatus[checkID] = syncStatus{inSync: false}
	l.changeMade()
}

// UpdateCheck is used to update the status of a check
func (l *State) UpdateCheck(checkID types.CheckID, status, output string) {
	l.Lock()
	defer l.Unlock()

	check, ok := l.checks[checkID]
	if !ok {
		return
	}

	if l.discardCheckOutput.Load().(bool) {
		output = ""
	}

	// Update the critical time tracking (this doesn't cause a server updates
	// so we can always keep this up to date).
	if status == api.HealthCritical {
		_, wasCritical := l.checkCriticalTime[checkID]
		if !wasCritical {
			l.checkCriticalTime[checkID] = time.Now()
		}
	} else {
		delete(l.checkCriticalTime, checkID)
	}

	// Do nothing if update is idempotent
	if check.Status == status && check.Output == output {
		return
	}

	// Defer a sync if the output has changed. This is an optimization around
	// frequent updates of output. Instead, we update the output internally,
	// and periodically do a write-back to the servers. If there is a status
	// change we do the write immediately.
	if l.config.CheckUpdateInterval > 0 && check.Status == status {
		check.Output = output
		if _, ok := l.deferCheck[checkID]; !ok {
			intv := time.Duration(uint64(l.config.CheckUpdateInterval)/2) + lib.RandomStagger(l.config.CheckUpdateInterval)
			deferSync := time.AfterFunc(intv, func() {
				l.Lock()
				if _, ok := l.checkStatus[checkID]; ok {
					l.checkStatus[checkID] = syncStatus{inSync: false}
					l.changeMade()
				}
				delete(l.deferCheck, checkID)
				l.Unlock()
			})
			l.deferCheck[checkID] = deferSync
		}
		return
	}

	// Update status and mark out of sync
	check.Status = status
	check.Output = output
	l.checkStatus[checkID] = syncStatus{inSync: false}
	l.changeMade()
}

// Check returns the locally registered check that the
// agent is aware of and are being kept in sync with the server
func (l *State) Check(id types.CheckID) *structs.HealthCheck {
	l.RLock()
	defer l.RUnlock()
	return l.checks[id]
}

// Checks returns the locally registered checks that the
// agent is aware of and are being kept in sync with the server
func (l *State) Checks() map[types.CheckID]*structs.HealthCheck {
	l.RLock()
	defer l.RUnlock()

	checks := make(map[types.CheckID]*structs.HealthCheck)
	for id, c := range l.checks {
		c2 := new(structs.HealthCheck)
		*c2 = *c
		checks[id] = c2
	}
	return checks
}

// CriticalCheck is used to return the duration a check has been critical along
// with its associated health check.
type CriticalCheck struct {
	CriticalFor time.Duration
	Check       *structs.HealthCheck
}

// CriticalChecks returns locally registered health checks that the agent is
// aware of and are being kept in sync with the server, and that are in a
// critical state. This also returns information about how long each check has
// been critical.
func (l *State) CriticalChecks() map[types.CheckID]CriticalCheck {
	checks := make(map[types.CheckID]CriticalCheck)

	l.RLock()
	defer l.RUnlock()

	now := time.Now()
	for checkID, criticalTime := range l.checkCriticalTime {
		checks[checkID] = CriticalCheck{
			CriticalFor: now.Sub(criticalTime),
			Check:       l.checks[checkID],
		}
	}

	return checks
}

// Metadata returns the local node metadata fields that the
// agent is aware of and are being kept in sync with the server
func (l *State) Metadata() map[string]string {
	metadata := make(map[string]string)
	l.RLock()
	defer l.RUnlock()

	for key, value := range l.metadata {
		metadata[key] = value
	}
	return metadata
}

// UpdateSyncState does a read of the server state, and updates
// the local sync status as appropriate
func (l *State) UpdateSyncState() error {
	req := structs.NodeSpecificRequest{
		Datacenter:   l.config.Datacenter,
		Node:         l.config.NodeName,
		QueryOptions: structs.QueryOptions{Token: l.tokens.AgentToken()},
	}
	var out1 structs.IndexedNodeServices
	var out2 structs.IndexedHealthChecks
	if e := l.delegate.RPC("Catalog.NodeServices", &req, &out1); e != nil {
		return e
	}
	if err := l.delegate.RPC("Health.NodeChecks", &req, &out2); err != nil {
		return err
	}
	checks := out2.HealthChecks

	l.Lock()
	defer l.Unlock()

	// Check the node info
	if out1.NodeServices == nil || out1.NodeServices.Node == nil ||
		out1.NodeServices.Node.ID != l.config.NodeID ||
		!reflect.DeepEqual(out1.NodeServices.Node.TaggedAddresses, l.config.TaggedAddresses) ||
		!reflect.DeepEqual(out1.NodeServices.Node.Meta, l.metadata) {
		l.nodeInfoInSync = false
	}

	// Check all our services
	services := make(map[string]*structs.NodeService)
	if out1.NodeServices != nil {
		services = out1.NodeServices.Services
	}

	for id := range l.services {
		// If the local service doesn't exist remotely, then sync it
		if _, ok := services[id]; !ok {
			l.serviceStatus[id] = syncStatus{inSync: false}
		}
	}

	for id, service := range services {
		// If we don't have the service locally, deregister it
		existing, ok := l.services[id]
		if !ok {
			// The consul service is created automatically, and does
			// not need to be deregistered.
			if id == structs.ConsulServiceID {
				continue
			}
			l.serviceStatus[id] = syncStatus{inSync: false}
			continue
		}

		// If our definition is different, we need to update it. Make a
		// copy so that we don't retain a pointer to any actual state
		// store info for in-memory RPCs.
		if existing.EnableTagOverride {
			existing.Tags = make([]string, len(service.Tags))
			copy(existing.Tags, service.Tags)
		}
		equal := existing.IsSame(service)
		l.serviceStatus[id] = syncStatus{inSync: equal}
	}

	// Index the remote health checks to improve efficiency
	checkIndex := make(map[types.CheckID]*structs.HealthCheck, len(checks))
	for _, check := range checks {
		checkIndex[check.CheckID] = check
	}

	// Sync any check which doesn't exist on the remote side
	for id := range l.checks {
		if _, ok := checkIndex[id]; !ok {
			l.checkStatus[id] = syncStatus{inSync: false}
		}
	}

	for _, check := range checks {
		// If we don't have the check locally, deregister it
		id := check.CheckID
		existing, ok := l.checks[id]
		if !ok {
			// The Serf check is created automatically, and does not
			// need to be deregistered.
			if id == structs.SerfCheckID {
				continue
			}
			l.checkStatus[id] = syncStatus{inSync: false}
			continue
		}

		// If our definition is different, we need to update it
		var equal bool
		if l.config.CheckUpdateInterval == 0 {
			equal = existing.IsSame(check)
		} else {
			// Copy the existing check before potentially modifying
			// it before the compare operation.
			eCopy := existing.Clone()

			// Copy the server's check before modifying, otherwise
			// in-memory RPCs will have side effects.
			cCopy := check.Clone()

			// If there's a defer timer active then we've got a
			// potentially spammy check so we don't sync the output
			// during this sweep since the timer will mark the check
			// out of sync for us. Otherwise, it is safe to sync the
			// output now. This is especially important for checks
			// that don't change state after they are created, in
			// which case we'd never see their output synced back ever.
			if _, ok := l.deferCheck[id]; ok {
				eCopy.Output = ""
				cCopy.Output = ""
			}
			equal = eCopy.IsSame(cCopy)
		}

		// Update the status
		l.checkStatus[id] = syncStatus{inSync: equal}
	}
	return nil
}

// SyncChanges is used to scan the status our local services and checks
// and update any that are out of sync with the server
func (l *State) SyncChanges() error {
	l.Lock()
	defer l.Unlock()

	// We will do node-level info syncing at the end, since it will get
	// updated by a service or check sync anyway, given how the register
	// API works.

	// Sync the services
	for id, status := range l.serviceStatus {
		if _, ok := l.services[id]; !ok {
			if err := l.deleteService(id); err != nil {
				return err
			}
		} else if !status.inSync {
			if err := l.syncService(id); err != nil {
				return err
			}
		} else {
			l.logger.Printf("[DEBUG] agent: Service '%s' in sync", id)
		}
	}

	// Sync the checks
	for id, status := range l.checkStatus {
		if _, ok := l.checks[id]; !ok {
			if err := l.deleteCheck(id); err != nil {
				return err
			}
		} else if !status.inSync {
			// Cancel a deferred sync
			if timer := l.deferCheck[id]; timer != nil {
				timer.Stop()
				delete(l.deferCheck, id)
			}

			if err := l.syncCheck(id); err != nil {
				return err
			}
		} else {
			l.logger.Printf("[DEBUG] agent: Check '%s' in sync", id)
		}
	}

	// Now sync the node level info if we need to, and didn't do any of
	// the other sync operations.
	if !l.nodeInfoInSync {
		if err := l.syncNodeInfo(); err != nil {
			return err
		}
	} else {
		l.logger.Printf("[DEBUG] agent: Node info in sync")
	}

	return nil
}

// LoadMetadata loads node metadata fields from the agent config and
// updates them on the local agent.
func (l *State) LoadMetadata(data map[string]string) error {
	l.Lock()
	defer l.Unlock()

	for k, v := range data {
		l.metadata[k] = v
	}
	l.changeMade()
	return nil
}

// UnloadMetadata resets the local metadata state
func (l *State) UnloadMetadata() {
	l.Lock()
	defer l.Unlock()
	l.metadata = make(map[string]string)
}

// Stats is used to get various debugging state from the sub-systems
func (l *State) Stats() map[string]string {
	l.RLock()
	defer l.RUnlock()
	return map[string]string{
		"services": strconv.Itoa(len(l.services)),
		"checks":   strconv.Itoa(len(l.checks)),
	}
}

// deleteService is used to delete a service from the server
func (l *State) deleteService(id string) error {
	if id == "" {
		return fmt.Errorf("ServiceID missing")
	}

	req := structs.DeregisterRequest{
		Datacenter:   l.config.Datacenter,
		Node:         l.config.NodeName,
		ServiceID:    id,
		WriteRequest: structs.WriteRequest{Token: l.serviceToken(id)},
	}
	var out struct{}
	err := l.delegate.RPC("Catalog.Deregister", &req, &out)
	if err == nil || strings.Contains(err.Error(), "Unknown service") {
		delete(l.serviceStatus, id)
		delete(l.serviceTokens, id)
		l.logger.Printf("[INFO] agent: Deregistered service '%s'", id)
		return nil
	} else if acl.IsErrPermissionDenied(err) {
		l.serviceStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[WARN] agent: Service '%s' deregistration blocked by ACLs", id)
		return nil
	}
	return err
}

// deleteCheck is used to delete a check from the server
func (l *State) deleteCheck(id types.CheckID) error {
	if id == "" {
		return fmt.Errorf("CheckID missing")
	}

	req := structs.DeregisterRequest{
		Datacenter:   l.config.Datacenter,
		Node:         l.config.NodeName,
		CheckID:      id,
		WriteRequest: structs.WriteRequest{Token: l.checkToken(id)},
	}
	var out struct{}
	err := l.delegate.RPC("Catalog.Deregister", &req, &out)
	if err == nil || strings.Contains(err.Error(), "Unknown check") {
		delete(l.checkStatus, id)
		delete(l.checkTokens, id)
		l.logger.Printf("[INFO] agent: Deregistered check '%s'", id)
		return nil
	} else if acl.IsErrPermissionDenied(err) {
		l.checkStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[WARN] agent: Check '%s' deregistration blocked by ACLs", id)
		return nil
	}
	return err
}

// syncService is used to sync a service to the server
func (l *State) syncService(id string) error {
	req := structs.RegisterRequest{
		Datacenter:      l.config.Datacenter,
		ID:              l.config.NodeID,
		Node:            l.config.NodeName,
		Address:         l.config.AdvertiseAddr,
		TaggedAddresses: l.config.TaggedAddresses,
		NodeMeta:        l.metadata,
		Service:         l.services[id],
		WriteRequest:    structs.WriteRequest{Token: l.serviceToken(id)},
	}

	// If the service has associated checks that are out of sync,
	// piggyback them on the service sync so they are part of the
	// same transaction and are registered atomically. We only let
	// checks ride on service registrations with the same token,
	// otherwise we need to register them separately so they don't
	// pick up privileges from the service token.
	var checks structs.HealthChecks
	for _, check := range l.checks {
		if check.ServiceID == id && (l.serviceToken(id) == l.checkToken(check.CheckID)) {
			if stat, ok := l.checkStatus[check.CheckID]; !ok || !stat.inSync {
				checks = append(checks, check)
			}
		}
	}

	// Backwards-compatibility for Consul < 0.5
	if len(checks) == 1 {
		req.Check = checks[0]
	} else {
		req.Checks = checks
	}

	var out struct{}
	err := l.delegate.RPC("Catalog.Register", &req, &out)
	if err == nil {
		l.serviceStatus[id] = syncStatus{inSync: true}
		// Given how the register API works, this info is also updated
		// every time we sync a service.
		l.nodeInfoInSync = true
		l.logger.Printf("[INFO] agent: Synced service '%s'", id)
		for _, check := range checks {
			l.checkStatus[check.CheckID] = syncStatus{inSync: true}
		}
	} else if acl.IsErrPermissionDenied(err) {
		l.serviceStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[WARN] agent: Service '%s' registration blocked by ACLs", id)
		for _, check := range checks {
			l.checkStatus[check.CheckID] = syncStatus{inSync: true}
		}
		return nil
	}
	return err
}

// syncCheck is used to sync a check to the server
func (l *State) syncCheck(id types.CheckID) error {
	// Pull in the associated service if any
	check := l.checks[id]
	var service *structs.NodeService
	if check.ServiceID != "" {
		if serv, ok := l.services[check.ServiceID]; ok {
			service = serv
		}
	}

	req := structs.RegisterRequest{
		Datacenter:      l.config.Datacenter,
		ID:              l.config.NodeID,
		Node:            l.config.NodeName,
		Address:         l.config.AdvertiseAddr,
		TaggedAddresses: l.config.TaggedAddresses,
		NodeMeta:        l.metadata,
		Service:         service,
		Check:           l.checks[id],
		WriteRequest:    structs.WriteRequest{Token: l.checkToken(id)},
	}
	var out struct{}
	err := l.delegate.RPC("Catalog.Register", &req, &out)
	if err == nil {
		l.checkStatus[id] = syncStatus{inSync: true}
		// Given how the register API works, this info is also updated
		// every time we sync a check.
		l.nodeInfoInSync = true
		l.logger.Printf("[INFO] agent: Synced check '%s'", id)
	} else if acl.IsErrPermissionDenied(err) {
		l.checkStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[WARN] agent: Check '%s' registration blocked by ACLs", id)
		return nil
	}
	return err
}

func (l *State) syncNodeInfo() error {
	req := structs.RegisterRequest{
		Datacenter:      l.config.Datacenter,
		ID:              l.config.NodeID,
		Node:            l.config.NodeName,
		Address:         l.config.AdvertiseAddr,
		TaggedAddresses: l.config.TaggedAddresses,
		NodeMeta:        l.metadata,
		WriteRequest:    structs.WriteRequest{Token: l.tokens.AgentToken()},
	}
	var out struct{}
	err := l.delegate.RPC("Catalog.Register", &req, &out)
	if err == nil {
		l.nodeInfoInSync = true
		l.logger.Printf("[INFO] agent: Synced node info")
	} else if acl.IsErrPermissionDenied(err) {
		l.nodeInfoInSync = true
		l.logger.Printf("[WARN] agent: Node info update blocked by ACLs")
		return nil
	}
	return err
}
