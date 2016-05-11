package agent

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/consul/structs"
)

// fixupKVOps takes the raw decoded JSON and base64 decodes values in KV ops,
// replacing them with byte arrays.
func fixupKVOps(raw interface{}) error {
	// decodeValue decodes the value member of the given operation.
	decodeValue := func(rawKV interface{}) error {
		rawMap, ok := rawKV.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected raw KV type: %T", rawKV)
		}
		for k, v := range rawMap {
			switch strings.ToLower(k) {
			case "value":
				// Leave the byte slice nil if we have a nil
				// value.
				if v == nil {
					return nil
				}

				// Otherwise, base64 decode it.
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("unexpected value type: %T", v)
				}
				decoded, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					return fmt.Errorf("failed to decode value: %v", err)
				}
				rawMap[k] = decoded
				return nil
			}
		}
		return nil
	}

	// fixupKVOp looks for non-nil KV operations and passes them on for
	// value conversion.
	fixupKVOp := func(rawOp interface{}) error {
		rawMap, ok := rawOp.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected raw op type: %T", rawOp)
		}
		for k, v := range rawMap {
			switch strings.ToLower(k) {
			case "kv":
				if v == nil {
					return nil
				}
				return decodeValue(v)
			}
		}
		return nil
	}

	rawSlice, ok := raw.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected raw type: %t", raw)
	}
	for _, rawOp := range rawSlice {
		if err := fixupKVOp(rawOp); err != nil {
			return err
		}
	}
	return nil
}

// Txn handles requests to apply multiple operations in a single, atomic
// transaction.
func (s *HTTPServer) Txn(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	if req.Method != "PUT" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return nil, nil
	}

	var args structs.TxnRequest
	s.parseDC(req, &args.Datacenter)
	s.parseToken(req, &args.Token)

	// Note the body is in API format, and not the RPC format. If we can't
	// decode it, we will return a 400 since we don't have enough context to
	// associate the error with a given operation.
	var ops api.TxnOps
	if err := decodeBody(req, &ops, fixupKVOps); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte(fmt.Sprintf("Failed to parse body: %v", err)))
		return nil, nil
	}

	// Convert the KV API format into the RPC format. Note that fixupKVOps
	// above will have already converted the base64 encoded strings into
	// byte arrays so we can assign right over.
	for _, in := range ops {
		if in.KV != nil {
			if size := len(in.KV.Value); size > maxKVSize {
				resp.WriteHeader(http.StatusRequestEntityTooLarge)
				resp.Write([]byte(fmt.Sprintf("Value for key %q is too large (%d > %d bytes)",
					in.KV.Key, size, maxKVSize)))
				return nil, nil
			}

			out := &structs.TxnOp{
				KV: &structs.TxnKVOp{
					Verb: structs.KVSOp(in.KV.Verb),
					DirEnt: structs.DirEntry{
						Key:     in.KV.Key,
						Value:   in.KV.Value,
						Flags:   in.KV.Flags,
						Session: in.KV.Session,
						RaftIndex: structs.RaftIndex{
							ModifyIndex: in.KV.Index,
						},
					},
				},
			}
			args.Ops = append(args.Ops, out)
		}
	}

	// Make the request and return a conflict status if there were errors
	// reported from the transaction.
	var reply structs.TxnResponse
	if err := s.agent.RPC("Txn.Apply", &args, &reply); err != nil {
		return nil, err
	}
	if len(reply.Errors) > 0 {
		var buf []byte
		var err error
		buf, err = s.marshalJSON(req, reply)
		if err != nil {
			return nil, err
		}

		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusConflict)
		resp.Write(buf)
		return nil, nil
	}

	// Otherwise, return the results of the successful transaction.
	return reply, nil
}
