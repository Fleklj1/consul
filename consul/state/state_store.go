package state

import (
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/go-memdb"
)

// StateStore is where we store all of Consul's state, including
// records of node registrations, services, checks, key/value
// pairs and more. The DB is entirely in-memory and is constructed
// from the Raft log through the FSM.
type StateStore struct {
	logger *log.Logger
	db     *memdb.MemDB
}

// IndexEntry keeps a record of the last index per-table.
type IndexEntry struct {
	Key   string
	Value uint64
}

// NewStateStore creates a new in-memory state storage layer.
func NewStateStore(logOutput io.Writer) (*StateStore, error) {
	// Create the in-memory DB
	db, err := memdb.NewMemDB(stateStoreSchema())
	if err != nil {
		return nil, fmt.Errorf("Failed setting up state store: %s", err)
	}

	// Create and return the state store
	s := &StateStore{
		logger: log.New(logOutput, "", log.LstdFlags),
		db:     db,
	}
	return s, nil
}

// EnsureNode is used to upsert node registration or modification.
func (s *StateStore) EnsureNode(idx uint64, node *structs.Node) error {
	tx := s.db.Txn(true)
	defer tx.Abort()

	// Call the node upsert
	if err := s.ensureNodeTxn(idx, node, tx); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// ensureNodeTxn is the inner function called to actually create a node
// registration or modify an existing one in the state store. It allows
// passing in a memdb transaction so it may be part of a larger txn.
func (s *StateStore) ensureNodeTxn(idx uint64, node *structs.Node, tx *memdb.Txn) error {
	// Check for an existing node
	existing, err := tx.First("nodes", "id", node.Node)
	if err != nil {
		return fmt.Errorf("node lookup failed: %s", err)
	}

	// Get the indexes
	if existing != nil {
		node.CreateIndex = existing.(*structs.Node).CreateIndex
		node.ModifyIndex = idx
	} else {
		node.CreateIndex = idx
		node.ModifyIndex = idx
	}

	// Insert the node and update the index
	if err := tx.Insert("nodes", node); err != nil {
		return fmt.Errorf("failed inserting node: %s", err)
	}
	if err := tx.Insert("index", &IndexEntry{"nodes", idx}); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}
	return nil
}

// GetNode is used to retrieve a node registration by node ID.
func (s *StateStore) GetNode(id string) (*structs.Node, error) {
	tx := s.db.Txn(false)
	defer tx.Abort()

	// Retrieve the node from the state store
	node, err := tx.First("nodes", "id", id)
	if err != nil {
		return nil, fmt.Errorf("node lookup failed: %s", err)
	}
	if node != nil {
		return node.(*structs.Node), nil
	}
	return nil, nil
}

// EnsureService is called to upsert creation of a given NodeService.
func (s *StateStore) EnsureService(idx uint64, node string, svc *structs.NodeService) error {
	tx := s.db.Txn(true)
	defer tx.Abort()

	// Call the service registration upsert
	if err := s.ensureServiceTxn(idx, node, svc, tx); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// ensureServiceTxn is used to upsert a service registration within an
// existing memdb transaction.
func (s *StateStore) ensureServiceTxn(idx uint64, node string, svc *structs.NodeService, tx *memdb.Txn) error {
	// Check for existing service
	existing, err := tx.First("services", "id", node, svc.Service)
	if err != nil {
		return fmt.Errorf("failed service lookup: %s", err)
	}

	// Create the service node entry
	entry := &structs.ServiceNode{
		Node:           node,
		ServiceID:      svc.ID,
		ServiceName:    svc.Service,
		ServiceTags:    svc.Tags,
		ServiceAddress: svc.Address,
		ServicePort:    svc.Port,
	}

	// Populate the indexes
	if existing != nil {
		entry.CreateIndex = existing.(*structs.NodeService).CreateIndex
		entry.ModifyIndex = idx
	} else {
		entry.CreateIndex = idx
		entry.ModifyIndex = idx
	}

	// Insert the service and update the index
	if err := tx.Insert("services", entry); err != nil {
		return fmt.Errorf("failed inserting service: %s", err)
	}
	if err := tx.Insert("index", &IndexEntry{"services", idx}); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}
	return nil
}

// NodeServices is used to query service registrations by node ID.
func (s *StateStore) NodeServices(nodeID string) (*structs.NodeServices, error) {
	tx := s.db.Txn(false)
	defer tx.Abort()

	// Query the node
	node, err := tx.First("nodes", "id", nodeID)
	if err != nil {
		return nil, fmt.Errorf("node lookup failed: %s", err)
	}
	if node == nil {
		return nil, nil
	}

	// Read all of the services
	services, err := tx.Get("services", "node", nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed querying services for node %q: %s", nodeID, err)
	}

	// Initialize the node services struct
	ns := &structs.NodeServices{
		Node:     *node.(*structs.Node),
		Services: make(map[string]*structs.NodeService),
	}

	// Add all of the services to the map
	for service := services.Next(); service != nil; service = services.Next() {
		sn := service.(*structs.ServiceNode)
		svc := &structs.NodeService{
			ID:      sn.ServiceID,
			Service: sn.ServiceName,
			Tags:    sn.ServiceTags,
			Address: sn.ServiceAddress,
			Port:    sn.ServicePort,
		}
		ns.Services[svc.ID] = svc
	}

	return ns, nil
}
