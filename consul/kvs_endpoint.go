package consul

import (
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/consul/structs"
	"time"
)

// KVS endpoint is used to manipulate the Key-Value store
type KVS struct {
	srv *Server
}

// Apply is used to apply a KVS request to the data store. This should
// only be used for operations that modify the data
func (c *Catalog) Apply(args *structs.KVSRequest, reply *bool) error {
	if done, err := c.srv.forward("KVS.Apply", args.Datacenter, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"consul", "kvs", "apply"}, time.Now())

	// Verify the args
	if args.DirEnt.Key == "" {
		return fmt.Errorf("Must provide key")
	}

	// Apply the update
	resp, err := c.srv.raftApply(structs.KVSRequestType, args)
	if err != nil {
		c.srv.logger.Printf("[ERR] consul.kvs: Apply failed: %v", err)
		return err
	}
	if respErr, ok := resp.(error); ok {
		return respErr
	}

	// Check if the return type is a bool
	if respBool, ok := resp.(bool); ok {
		*reply = respBool
	}
	return nil
}

// Get is used to lookup a single key
func (k *KVS) Get(args *structs.KeyRequest, reply *structs.IndexedDirEntries) error {
	if done, err := k.srv.forward("KVS.Get", args.Datacenter, args, reply); done {
		return err
	}

	// Get the local state
	state := k.srv.fsm.State()
	return k.srv.blockingRPC(&args.BlockingQuery,
		state.QueryTables("KVSGet"),
		func() (uint64, error) {
			index, ent, err := state.KVSGet(args.Key)
			if err != nil {
				return 0, err
			}
			if ent == nil {
				reply.Index = index
				reply.Entries = nil
			} else {
				reply.Index = ent.ModifyIndex
				reply.Entries = structs.DirEntries{ent}
			}
			return reply.Index, nil
		})
}

// List is used to list all keys with a given prefix
func (k *KVS) List(args *structs.KeyRequest, reply *structs.IndexedDirEntries) error {
	if done, err := k.srv.forward("KVS.List", args.Datacenter, args, reply); done {
		return err
	}

	// Get the local state
	state := k.srv.fsm.State()
	return k.srv.blockingRPC(&args.BlockingQuery,
		state.QueryTables("KVSList"),
		func() (uint64, error) {
			index, ent, err := state.KVSList(args.Key)
			if err != nil {
				return 0, err
			}
			if len(ent) == 0 {
				reply.Index = index
				reply.Entries = nil
			} else {
				// Determine the maximum affected index
				var maxIndex uint64
				for _, e := range ent {
					if e.ModifyIndex > maxIndex {
						maxIndex = e.ModifyIndex
					}
				}

				reply.Index = maxIndex
				reply.Entries = ent
			}
			return reply.Index, nil
		})
}
