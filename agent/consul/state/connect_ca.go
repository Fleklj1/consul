package state

import (
	"fmt"

	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/go-memdb"
)

const (
	caConfigTableName   = "connect-ca-config"
	caRootTableName     = "connect-ca-roots"
	caProviderTableName = "connect-ca-builtin"
)

// caConfigTableSchema returns a new table schema used for storing
// the CA config for Connect.
func caConfigTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: caConfigTableName,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:         "id",
				AllowMissing: true,
				Unique:       true,
				Indexer: &memdb.ConditionalIndex{
					Conditional: func(obj interface{}) (bool, error) { return true, nil },
				},
			},
		},
	}
}

// caRootTableSchema returns a new table schema used for storing
// CA roots for Connect.
func caRootTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: caRootTableName,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.UUIDFieldIndex{
					Field: "ID",
				},
			},
		},
	}
}

// caProviderTableSchema returns a new table schema used for storing
// the built-in CA provider's state for connect. This is only used by
// the internal Consul CA provider.
func caProviderTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: caProviderTableName,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:         "id",
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.ConditionalIndex{
					Conditional: func(obj interface{}) (bool, error) { return true, nil },
				},
			},
		},
	}
}

func init() {
	registerSchema(caConfigTableSchema)
	registerSchema(caRootTableSchema)
	registerSchema(caProviderTableSchema)
}

// CAConfig is used to pull the CA config from the snapshot.
func (s *Snapshot) CAConfig() (*structs.CAConfiguration, error) {
	c, err := s.tx.First(caConfigTableName, "id")
	if err != nil {
		return nil, err
	}

	config, ok := c.(*structs.CAConfiguration)
	if !ok {
		return nil, nil
	}

	return config, nil
}

// CAConfig is used when restoring from a snapshot.
func (s *Restore) CAConfig(config *structs.CAConfiguration) error {
	if err := s.tx.Insert(caConfigTableName, config); err != nil {
		return fmt.Errorf("failed restoring CA config: %s", err)
	}

	return nil
}

// CAConfig is used to get the current Autopilot configuration.
func (s *Store) CAConfig() (uint64, *structs.CAConfiguration, error) {
	tx := s.db.Txn(false)
	defer tx.Abort()

	// Get the autopilot config
	c, err := tx.First(caConfigTableName, "id")
	if err != nil {
		return 0, nil, fmt.Errorf("failed CA config lookup: %s", err)
	}

	config, ok := c.(*structs.CAConfiguration)
	if !ok {
		return 0, nil, nil
	}

	return config.ModifyIndex, config, nil
}

// CASetConfig is used to set the current Autopilot configuration.
func (s *Store) CASetConfig(idx uint64, config *structs.CAConfiguration) error {
	tx := s.db.Txn(true)
	defer tx.Abort()

	if err := s.caSetConfigTxn(idx, tx, config); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// CACheckAndSetConfig is used to try updating the CA configuration with a
// given Raft index. If the CAS index specified is not equal to the last observed index
// for the config, then the call is a noop,
func (s *Store) CACheckAndSetConfig(idx, cidx uint64, config *structs.CAConfiguration) (bool, error) {
	tx := s.db.Txn(true)
	defer tx.Abort()

	// Check for an existing config
	existing, err := tx.First(caConfigTableName, "id")
	if err != nil {
		return false, fmt.Errorf("failed CA config lookup: %s", err)
	}

	// If the existing index does not match the provided CAS
	// index arg, then we shouldn't update anything and can safely
	// return early here.
	e, ok := existing.(*structs.CAConfiguration)
	if !ok || e.ModifyIndex != cidx {
		return false, nil
	}

	if err := s.caSetConfigTxn(idx, tx, config); err != nil {
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (s *Store) caSetConfigTxn(idx uint64, tx *memdb.Txn, config *structs.CAConfiguration) error {
	// Check for an existing config
	prev, err := tx.First(caConfigTableName, "id")
	if err != nil {
		return fmt.Errorf("failed CA config lookup: %s", err)
	}

	// Set the indexes, prevent the cluster ID from changing.
	if prev != nil {
		existing := prev.(*structs.CAConfiguration)
		config.CreateIndex = existing.CreateIndex
		config.ClusterSerial = existing.ClusterSerial
	} else {
		config.CreateIndex = idx
	}
	config.ModifyIndex = idx

	if err := tx.Insert(caConfigTableName, config); err != nil {
		return fmt.Errorf("failed updating CA config: %s", err)
	}
	return nil
}

// CARoots is used to pull all the CA roots for the snapshot.
func (s *Snapshot) CARoots() (structs.CARoots, error) {
	ixns, err := s.tx.Get(caRootTableName, "id")
	if err != nil {
		return nil, err
	}

	var ret structs.CARoots
	for wrapped := ixns.Next(); wrapped != nil; wrapped = ixns.Next() {
		ret = append(ret, wrapped.(*structs.CARoot))
	}

	return ret, nil
}

// CARoots is used when restoring from a snapshot.
func (s *Restore) CARoot(r *structs.CARoot) error {
	// Insert
	if err := s.tx.Insert(caRootTableName, r); err != nil {
		return fmt.Errorf("failed restoring CA root: %s", err)
	}
	if err := indexUpdateMaxTxn(s.tx, r.ModifyIndex, caRootTableName); err != nil {
		return fmt.Errorf("failed updating index: %s", err)
	}

	return nil
}

// CARoots returns the list of all CA roots.
func (s *Store) CARoots(ws memdb.WatchSet) (uint64, structs.CARoots, error) {
	tx := s.db.Txn(false)
	defer tx.Abort()

	// Get the index
	idx := maxIndexTxn(tx, caRootTableName)

	// Get all
	iter, err := tx.Get(caRootTableName, "id")
	if err != nil {
		return 0, nil, fmt.Errorf("failed CA root lookup: %s", err)
	}
	ws.Add(iter.WatchCh())

	var results structs.CARoots
	for v := iter.Next(); v != nil; v = iter.Next() {
		results = append(results, v.(*structs.CARoot))
	}
	return idx, results, nil
}

// CARootActive returns the currently active CARoot.
func (s *Store) CARootActive(ws memdb.WatchSet) (uint64, *structs.CARoot, error) {
	// Get all the roots since there should never be that many and just
	// do the filtering in this method.
	var result *structs.CARoot
	idx, roots, err := s.CARoots(ws)
	if err == nil {
		for _, r := range roots {
			if r.Active {
				result = r
				break
			}
		}
	}

	return idx, result, err
}

// CARootSetCAS sets the current CA root state using a check-and-set operation.
// On success, this will replace the previous set of CARoots completely with
// the given set of roots.
//
// The first boolean result returns whether the transaction succeeded or not.
func (s *Store) CARootSetCAS(idx, cidx uint64, rs []*structs.CARoot) (bool, error) {
	tx := s.db.Txn(true)
	defer tx.Abort()

	// There must be exactly one active CA root.
	activeCount := 0
	for _, r := range rs {
		if r.Active {
			activeCount++
		}
	}
	if activeCount != 1 {
		return false, fmt.Errorf("there must be exactly one active CA")
	}

	// Get the current max index
	if midx := maxIndexTxn(tx, caRootTableName); midx != cidx {
		return false, nil
	}

	// Go through and find any existing matching CAs so we can preserve and
	// update their Create/ModifyIndex values.
	for _, r := range rs {
		if r.ID == "" {
			return false, ErrMissingCARootID
		}

		existing, err := tx.First(caRootTableName, "id", r.ID)
		if err != nil {
			return false, fmt.Errorf("failed CA root lookup: %s", err)
		}

		if existing != nil {
			r.CreateIndex = existing.(*structs.CARoot).CreateIndex
		} else {
			r.CreateIndex = idx
		}
		r.ModifyIndex = idx
	}

	// Delete all
	_, err := tx.DeleteAll(caRootTableName, "id")
	if err != nil {
		return false, err
	}

	// Insert all
	for _, r := range rs {
		if err := tx.Insert(caRootTableName, r); err != nil {
			return false, err
		}
	}

	// Update the index
	if err := tx.Insert("index", &IndexEntry{caRootTableName, idx}); err != nil {
		return false, fmt.Errorf("failed updating index: %s", err)
	}

	tx.Commit()
	return true, nil
}

// CAProviderState is used to pull the built-in provider state from the snapshot.
func (s *Snapshot) CAProviderState() (*structs.CAConsulProviderState, error) {
	c, err := s.tx.First(caProviderTableName, "id")
	if err != nil {
		return nil, err
	}

	state, ok := c.(*structs.CAConsulProviderState)
	if !ok {
		return nil, nil
	}

	return state, nil
}

// CAProviderState is used when restoring from a snapshot.
func (s *Restore) CAProviderState(state *structs.CAConsulProviderState) error {
	if err := s.tx.Insert(caProviderTableName, state); err != nil {
		return fmt.Errorf("failed restoring built-in CA state: %s", err)
	}

	return nil
}

// CAProviderState is used to get the current Consul CA provider state.
func (s *Store) CAProviderState() (uint64, *structs.CAConsulProviderState, error) {
	tx := s.db.Txn(false)
	defer tx.Abort()

	// Get the autopilot config
	c, err := tx.First(caProviderTableName, "id")
	if err != nil {
		return 0, nil, fmt.Errorf("failed built-in CA state lookup: %s", err)
	}

	state, ok := c.(*structs.CAConsulProviderState)
	if !ok {
		return 0, nil, nil
	}

	return state.ModifyIndex, state, nil
}

// CASetProviderState is used to set the current built-in CA provider state.
func (s *Store) CASetProviderState(idx uint64, state *structs.CAConsulProviderState) (bool, error) {
	tx := s.db.Txn(true)
	defer tx.Abort()

	// Check for an existing config
	existing, err := tx.First(caProviderTableName, "id")
	if err != nil {
		return false, fmt.Errorf("failed built-in CA state lookup: %s", err)
	}

	// Set the indexes.
	if existing != nil {
		state.CreateIndex = existing.(*structs.CAConfiguration).CreateIndex
	} else {
		state.CreateIndex = idx
	}
	state.ModifyIndex = idx

	if err := tx.Insert(caProviderTableName, state); err != nil {
		return false, fmt.Errorf("failed updating built-in CA state: %s", err)
	}
	tx.Commit()

	return true, nil
}
