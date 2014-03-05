package agent

import (
	"github.com/hashicorp/consul/consul"
	"github.com/hashicorp/consul/consul/structs"
	"log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

const (
	syncRetryIntv = 30 * time.Second
)

// syncStatus is used to represent the difference between
// the local and remote state, and if action needs to be taken
type syncStatus struct {
	remoteDelete bool // Should this be deleted from the server
	inSync       bool // Is this in sync with the server
}

// localState is used to represent the node's services,
// and checks. We used it to perform anti-entropy with the
// catalog representation
type localState struct {
	// paused is used to check if we are paused. Must be the first
	// element due to a go bug.
	paused int32

	sync.Mutex
	logger *log.Logger

	// Config is the agent config
	config *Config

	// iface is the consul interface to use for keeping in sync
	iface consul.Interface

	// Services tracks the local services
	services      map[string]*structs.NodeService
	serviceStatus map[string]syncStatus

	// Checks tracks the local checks
	checks      map[string]*structs.HealthCheck
	checkStatus map[string]syncStatus

	// consulCh is used to inform of a change to the known
	// consul nodes. This may be used to retry a sync run
	consulCh chan struct{}

	// triggerCh is used to inform of a change to local state
	// that requires anti-entropy with the server
	triggerCh chan struct{}
}

// Init is used to initialize the local state
func (l *localState) Init(config *Config, logger *log.Logger) {
	l.config = config
	l.logger = logger
	l.services = make(map[string]*structs.NodeService)
	l.serviceStatus = make(map[string]syncStatus)
	l.checks = make(map[string]*structs.HealthCheck)
	l.checkStatus = make(map[string]syncStatus)
	l.consulCh = make(chan struct{}, 1)
	l.triggerCh = make(chan struct{}, 1)
}

// SetIface is used to set the Consul interface. Must be set prior to
// starting anti-entropy
func (l *localState) SetIface(iface consul.Interface) {
	l.iface = iface
}

// changeMade is used to trigger an anti-entropy run
func (l *localState) changeMade() {
	select {
	case l.triggerCh <- struct{}{}:
	default:
	}
}

// ConsulServerUp is used to inform that a new consul server is now
// up. This can be used to speed up the sync process if we are blocking
// waiting to discover a consul server
func (l *localState) ConsulServerUp() {
	select {
	case l.consulCh <- struct{}{}:
	default:
	}
}

// Pause is used to pause state syncronization, this can be
// used to make batch changes
func (l *localState) Pause() {
	atomic.StoreInt32(&l.paused, 1)
}

// Resume is used to resume state syncronization
func (l *localState) Resume() {
	atomic.StoreInt32(&l.paused, 0)
	l.changeMade()
}

// isPaused is used to check if we are paused
func (l *localState) isPaused() bool {
	return atomic.LoadInt32(&l.paused) == 1
}

// AddService is used to add a service entry to the local state.
// This entry is persistent and the agent will make a best effort to
// ensure it is registered
func (l *localState) AddService(service *structs.NodeService) {
	// Assign the ID if none given
	if service.ID == "" && service.Service != "" {
		service.ID = service.Service
	}

	l.Lock()
	defer l.Unlock()

	l.services[service.ID] = service
	l.serviceStatus[service.ID] = syncStatus{}
	l.changeMade()
}

// RemoveService is used to remove a service entry from the local state.
// The agent will make a best effort to ensure it is deregistered
func (l *localState) RemoveService(serviceID string) {
	l.Lock()
	defer l.Unlock()

	delete(l.services, serviceID)
	l.serviceStatus[serviceID] = syncStatus{remoteDelete: true}
	l.changeMade()
}

// Services returns the locally registered services that the
// agent is aware of and are being kept in sync with the server
func (l *localState) Services() map[string]*structs.NodeService {
	services := make(map[string]*structs.NodeService)
	l.Lock()
	defer l.Unlock()

	for name, serv := range l.services {
		services[name] = serv
	}
	return services
}

// AddCheck is used to add a health check to the local state.
// This entry is persistent and the agent will make a best effort to
// ensure it is registered
func (l *localState) AddCheck(check *structs.HealthCheck) {
	// Set the node name
	check.Node = l.config.NodeName

	l.Lock()
	defer l.Unlock()

	l.checks[check.CheckID] = check
	l.checkStatus[check.CheckID] = syncStatus{}
	l.changeMade()
}

// RemoveCheck is used to remove a health check from the local state.
// The agent will make a best effort to ensure it is deregistered
func (l *localState) RemoveCheck(checkID string) {
	l.Lock()
	defer l.Unlock()

	delete(l.checks, checkID)
	l.checkStatus[checkID] = syncStatus{remoteDelete: true}
	l.changeMade()
}

// UpdateCheck is used to update the status of a check
func (l *localState) UpdateCheck(checkID, status, output string) {
	l.Lock()
	defer l.Unlock()

	check, ok := l.checks[checkID]
	if !ok {
		return
	}

	// Do nothing if update is idempotent
	if check.Status == status && check.Notes == output {
		return
	}

	// Update status and mark out of sync
	check.Status = status
	check.Notes = output
	l.checkStatus[checkID] = syncStatus{inSync: false}
	l.changeMade()
}

// Checks returns the locally registered checks that the
// agent is aware of and are being kept in sync with the server
func (l *localState) Checks() map[string]*structs.HealthCheck {
	checks := make(map[string]*structs.HealthCheck)
	l.Lock()
	defer l.Unlock()

	for name, check := range l.checks {
		checks[name] = check
	}
	return checks
}

// antiEntropy is a long running method used to perform anti-entropy
// between local and remote state.
func (l *localState) antiEntropy(shutdownCh chan struct{}) {
SYNC:
	// Sync our state with the servers
	for {
		if err := l.setSyncState(); err != nil {
			l.logger.Printf("[ERR] agent: failed to sync remote state: %v", err)
			select {
			case <-l.consulCh:
			case <-time.After(aeScale(syncRetryIntv, len(l.iface.LANMembers()))):
			case <-shutdownCh:
				return
			}
		}
		break
	}

	// Force-trigger AE to pickup any changes
	l.changeMade()

	// Schedule the next full sync, with a random stagger
	aeIntv := aeScale(l.config.AEInterval, len(l.iface.LANMembers()))
	aeIntv = aeIntv + randomStagger(aeIntv)
	aeTimer := time.After(aeIntv)

	// Wait for sync events
	for {
		select {
		case <-aeTimer:
			goto SYNC
		case <-l.triggerCh:
			// Skip the sync if we are paused
			if l.isPaused() {
				continue
			}
			if err := l.syncChanges(); err != nil {
				l.logger.Printf("[ERR] agent: failed to sync changes: %v", err)
			}
		case <-shutdownCh:
			return
		}
	}
}

// setSyncState does a read of the server state, and updates
// the local syncStatus as appropriate
func (l *localState) setSyncState() error {
	req := structs.NodeSpecificRequest{
		Datacenter: l.config.Datacenter,
		Node:       l.config.NodeName,
	}
	var out1 structs.IndexedNodeServices
	var out2 structs.IndexedHealthChecks
	if e := l.iface.RPC("Catalog.NodeServices", &req, &out1); e != nil {
		return e
	}
	if err := l.iface.RPC("Health.NodeChecks", &req, &out2); err != nil {
		return err
	}
	services := out1.NodeServices
	checks := out2.HealthChecks

	l.Lock()
	defer l.Unlock()

	if services != nil {
		for id, service := range services.Services {
			// If we don't have the service locally, deregister it
			existing, ok := l.services[id]
			if !ok {
				// The Consul service is created automatically, and
				// does not need to be registered
				if id == consul.ConsulServiceID && l.config.Server {
					continue
				}
				l.serviceStatus[id] = syncStatus{remoteDelete: true}
				continue
			}

			// If our definition is different, we need to update it
			equal := reflect.DeepEqual(existing, service)
			l.serviceStatus[id] = syncStatus{inSync: equal}
		}
	}

	for _, check := range checks {
		// If we don't have the check locally, deregister it
		id := check.CheckID
		existing, ok := l.checks[id]
		if !ok {
			// The Serf check is created automatically, and does not
			// need to be registered
			if id == consul.SerfCheckID {
				continue
			}
			l.checkStatus[id] = syncStatus{remoteDelete: true}
			continue
		}

		// If our definition is different, we need to update it
		equal := reflect.DeepEqual(existing, check)
		l.checkStatus[id] = syncStatus{inSync: equal}
	}
	return nil
}

// syncChanges is used to scan the status our local services and checks
// and update any that are out of sync with the server
func (l *localState) syncChanges() error {
	l.Lock()
	defer l.Unlock()

	// Sync the services
	for id, status := range l.serviceStatus {
		if status.remoteDelete {
			if err := l.deleteService(id); err != nil {
				return err
			}
		} else if !status.inSync {
			if err := l.syncService(id); err != nil {
				return err
			}
		}
	}

	// Sync the checks
	for id, status := range l.checkStatus {
		if status.remoteDelete {
			if err := l.deleteCheck(id); err != nil {
				return err
			}
		} else if !status.inSync {
			if err := l.syncCheck(id); err != nil {
				return err
			}
		}
	}
	return nil
}

// deleteService is used to delete a service from the server
func (l *localState) deleteService(id string) error {
	req := structs.DeregisterRequest{
		Datacenter: l.config.Datacenter,
		Node:       l.config.NodeName,
		ServiceID:  id,
	}
	var out struct{}
	err := l.iface.RPC("Catalog.Deregister", &req, &out)
	if err == nil {
		delete(l.serviceStatus, id)
		l.logger.Printf("[INFO] agent: Deregistered service '%s'", id)
	}
	return err
}

// deleteCheck is used to delete a service from the server
func (l *localState) deleteCheck(id string) error {
	req := structs.DeregisterRequest{
		Datacenter: l.config.Datacenter,
		Node:       l.config.NodeName,
		CheckID:    id,
	}
	var out struct{}
	err := l.iface.RPC("Catalog.Deregister", &req, &out)
	if err == nil {
		delete(l.checkStatus, id)
		l.logger.Printf("[INFO] agent: Deregistered check '%s'", id)
	}
	return err
}

// syncService is used to sync a service to the server
func (l *localState) syncService(id string) error {
	req := structs.RegisterRequest{
		Datacenter: l.config.Datacenter,
		Node:       l.config.NodeName,
		Address:    l.config.AdvertiseAddr,
		Service:    l.services[id],
	}
	var out struct{}
	err := l.iface.RPC("Catalog.Register", &req, &out)
	if err == nil {
		l.serviceStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[INFO] agent: Synced service '%s'", id)
	}
	return err
}

// syncCheck is used to sync a service to the server
func (l *localState) syncCheck(id string) error {
	// Pull in the associated service if any
	check := l.checks[id]
	var service *structs.NodeService
	if check.ServiceID != "" {
		if serv, ok := l.services[check.ServiceID]; ok {
			service = serv
		}
	}
	req := structs.RegisterRequest{
		Datacenter: l.config.Datacenter,
		Node:       l.config.NodeName,
		Address:    l.config.AdvertiseAddr,
		Service:    service,
		Check:      l.checks[id],
	}
	var out struct{}
	err := l.iface.RPC("Catalog.Register", &req, &out)
	if err == nil {
		l.checkStatus[id] = syncStatus{inSync: true}
		l.logger.Printf("[INFO] agent: Synced check '%s'", id)
	}
	return err
}
