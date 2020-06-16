// +build !consulent

package state

import (
	"fmt"

	"github.com/hashicorp/consul/agent/structs"
	memdb "github.com/hashicorp/go-memdb"
)

// servicesTableSchema returns a new table schema used to store information
// about services.
func servicesTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "services",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Node",
							Lowercase: true,
						},
						&memdb.StringFieldIndex{
							Field:     "ServiceID",
							Lowercase: true,
						},
					},
				},
			},
			"node": {
				Name:         "node",
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Node",
					Lowercase: true,
				},
			},
			"service": {
				Name:         "service",
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "ServiceName",
					Lowercase: true,
				},
			},
			"connect": {
				Name:         "connect",
				AllowMissing: true,
				Unique:       false,
				Indexer:      &IndexConnectService{},
			},
			"kind": {
				Name:         "kind",
				AllowMissing: false,
				Unique:       false,
				Indexer:      &IndexServiceKind{},
			},
		},
	}
}

// checksTableSchema returns a new table schema used for storing and indexing
// health check information. Health checks have a number of different attributes
// we want to filter by, so this table is a bit more complex.
func checksTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "checks",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Node",
							Lowercase: true,
						},
						&memdb.StringFieldIndex{
							Field:     "CheckID",
							Lowercase: true,
						},
					},
				},
			},
			"status": {
				Name:         "status",
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Status",
					Lowercase: false,
				},
			},
			"service": {
				Name:         "service",
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "ServiceName",
					Lowercase: true,
				},
			},
			"node": {
				Name:         "node",
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Node",
					Lowercase: true,
				},
			},
			"node_service_check": {
				Name:         "node_service_check",
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Node",
							Lowercase: true,
						},
						&memdb.FieldSetIndex{
							Field: "ServiceID",
						},
					},
				},
			},
			"node_service": {
				Name:         "node_service",
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Node",
							Lowercase: true,
						},
						&memdb.StringFieldIndex{
							Field:     "ServiceID",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

func serviceIndexName(name string, _ *structs.EnterpriseMeta) string {
	return fmt.Sprintf("service.%s", name)
}

func serviceKindIndexName(kind structs.ServiceKind, _ *structs.EnterpriseMeta) string {
	switch kind {
	case structs.ServiceKindTypical:
		// needs a special case here
		return "service_kind.typical"
	default:
		return "service_kind." + string(kind)
	}
}

func (s *Store) catalogUpdateServicesIndexes(tx *txn, idx uint64, _ *structs.EnterpriseMeta) error {
	// overall services index
	if err := indexUpdateMaxTxn(tx, idx, "services"); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}

	return nil
}

func (s *Store) catalogUpdateServiceKindIndexes(tx *txn, kind structs.ServiceKind, idx uint64, _ *structs.EnterpriseMeta) error {
	// service-kind index
	if err := indexUpdateMaxTxn(tx, idx, serviceKindIndexName(kind, nil)); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}

	return nil
}

func (s *Store) catalogUpdateServiceIndexes(tx *txn, serviceName string, idx uint64, _ *structs.EnterpriseMeta) error {
	// per-service index
	if err := indexUpdateMaxTxn(tx, idx, serviceIndexName(serviceName, nil)); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}

	return nil
}

func (s *Store) catalogUpdateServiceExtinctionIndex(tx *txn, idx uint64, _ *structs.EnterpriseMeta) error {
	if err := tx.Insert("index", &IndexEntry{serviceLastExtinctionIndexName, idx}); err != nil {
		return fmt.Errorf("failed updating missing service extinction index: %s", err)
	}
	return nil
}

func (s *Store) catalogInsertService(tx *txn, svc *structs.ServiceNode) error {
	// Insert the service and update the index
	if err := tx.Insert("services", svc); err != nil {
		return fmt.Errorf("failed inserting service: %s", err)
	}

	if err := s.catalogUpdateServicesIndexes(tx, svc.ModifyIndex, &svc.EnterpriseMeta); err != nil {
		return err
	}

	if err := s.catalogUpdateServiceIndexes(tx, svc.ServiceName, svc.ModifyIndex, &svc.EnterpriseMeta); err != nil {
		return err
	}

	if err := s.catalogUpdateServiceKindIndexes(tx, svc.ServiceKind, svc.ModifyIndex, &svc.EnterpriseMeta); err != nil {
		return err
	}

	return nil
}

func (s *Store) catalogServicesMaxIndex(tx *txn, _ *structs.EnterpriseMeta) uint64 {
	return maxIndexTxn(tx, "services")
}

func (s *Store) catalogServiceMaxIndex(tx *txn, serviceName string, _ *structs.EnterpriseMeta) (<-chan struct{}, interface{}, error) {
	return tx.FirstWatch("index", "id", serviceIndexName(serviceName, nil))
}

func (s *Store) catalogServiceKindMaxIndex(tx *txn, ws memdb.WatchSet, kind structs.ServiceKind, entMeta *structs.EnterpriseMeta) uint64 {
	return maxIndexWatchTxn(tx, ws, serviceKindIndexName(kind, nil))
}

func (s *Store) catalogServiceList(tx *txn, _ *structs.EnterpriseMeta, _ bool) (memdb.ResultIterator, error) {
	return tx.Get("services", "id")
}

func (s *Store) catalogServiceListByKind(tx *txn, kind structs.ServiceKind, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("services", "kind", string(kind))
}

func (s *Store) catalogServiceListByNode(tx *txn, node string, _ *structs.EnterpriseMeta, _ bool) (memdb.ResultIterator, error) {
	return tx.Get("services", "node", node)
}

func (s *Store) catalogServiceNodeList(tx *txn, name string, index string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("services", index, name)
}

func (s *Store) catalogServiceLastExtinctionIndex(tx *txn, _ *structs.EnterpriseMeta) (interface{}, error) {
	return tx.First("index", "id", serviceLastExtinctionIndexName)
}

func (s *Store) catalogMaxIndex(tx *txn, _ *structs.EnterpriseMeta, checks bool) uint64 {
	if checks {
		return maxIndexTxn(tx, "nodes", "services", "checks")
	}
	return maxIndexTxn(tx, "nodes", "services")
}

func (s *Store) catalogMaxIndexWatch(tx *txn, ws memdb.WatchSet, _ *structs.EnterpriseMeta, checks bool) uint64 {
	if checks {
		return maxIndexWatchTxn(tx, ws, "nodes", "services", "checks")
	}
	return maxIndexWatchTxn(tx, ws, "nodes", "services")
}

func (s *Store) catalogUpdateCheckIndexes(tx *txn, idx uint64, _ *structs.EnterpriseMeta) error {
	// update the universal index entry
	if err := tx.Insert("index", &IndexEntry{"checks", idx}); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}
	return nil
}

func (s *Store) catalogChecksMaxIndex(tx *txn, _ *structs.EnterpriseMeta) uint64 {
	return maxIndexTxn(tx, "checks")
}

func (s *Store) catalogListChecksByNode(tx *txn, node string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("checks", "node", node)
}

func (s *Store) catalogListChecksByService(tx *txn, service string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("checks", "service", service)
}

func (s *Store) catalogListChecksInState(tx *txn, state string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	// simpler than normal due to the use of the CompoundMultiIndex
	return tx.Get("checks", "status", state)
}

func (s *Store) catalogListChecks(tx *txn, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("checks", "id")
}

func (s *Store) catalogListNodeChecks(tx *txn, node string) (memdb.ResultIterator, error) {
	return tx.Get("checks", "node_service_check", node, false)
}

func (s *Store) catalogListServiceChecks(tx *txn, node string, service string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("checks", "node_service", node, service)
}

func (s *Store) catalogInsertCheck(tx *txn, chk *structs.HealthCheck, idx uint64) error {
	// Insert the check
	if err := tx.Insert("checks", chk); err != nil {
		return fmt.Errorf("failed inserting check: %s", err)
	}

	if err := s.catalogUpdateCheckIndexes(tx, idx, &chk.EnterpriseMeta); err != nil {
		return err
	}

	return nil
}

func (s *Store) catalogChecksForNodeService(tx *txn, node string, service string, entMeta *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get("checks", "node_service", node, service)
}

func (s *Store) validateRegisterRequestTxn(tx *txn, args *structs.RegisterRequest) (*structs.EnterpriseMeta, error) {
	return nil, nil
}

func (s *Store) ValidateRegisterRequest(args *structs.RegisterRequest) (*structs.EnterpriseMeta, error) {
	return nil, nil
}
