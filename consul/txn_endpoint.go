package consul

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/consul/structs"
)

// Txn endpoint is used to perform multi-object atomic transactions.
type Txn struct {
	srv *Server
}

// Apply is used to apply multiple operations in a single, atomic transaction.
func (t *Txn) Apply(args *structs.KVSAtomicRequest, reply *structs.KVSAtomicResponse) error {
	if done, err := t.srv.forward("Txn.Apply", args, args, reply); done {
		return err
	}
	defer metrics.MeasureSince([]string{"consul", "txn", "apply"}, time.Now())

	// Perform the pre-apply checks on each of the operations.
	acl, err := t.srv.resolveToken(args.Token)
	if err != nil {
		return err
	}
	for i, op := range args.Ops {
		ok, err := kvsPreApply(t.srv, acl, op.Op, &op.DirEnt)
		if err != nil {
			reply.Errors = append(reply.Errors, &structs.KVSAtomicError{i, err.Error()})
		} else if !ok {
			err = fmt.Errorf("failed to lock key %q due to lock delay", op.DirEnt.Key)
			reply.Errors = append(reply.Errors, &structs.KVSAtomicError{i, err.Error()})
		}
	}
	if len(reply.Errors) > 0 {
		return nil
	}

	// Apply the update.
	resp, err := t.srv.raftApply(structs.KVSAtomicRequestType, args)
	if err != nil {
		t.srv.logger.Printf("[ERR] consul.kvs: ApplyAtomic failed: %v", err)
		return err
	}
	if respErr, ok := resp.(error); ok {
		return respErr
	}

	// Convert the return type. This should be a cheap copy since we are
	// just taking the two slices.
	if respAtomic, ok := resp.(structs.KVSAtomicResponse); ok {
		*reply = respAtomic
	} else {
		return fmt.Errorf("unexpected return type %T", resp)
	}
	return nil
}
