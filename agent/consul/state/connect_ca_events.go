package state

import (
	"github.com/hashicorp/consul/acl"
	"github.com/hashicorp/consul/agent/consul/stream"
	"github.com/hashicorp/consul/agent/structs"
)

// EventTopicCARoots is the streaming topic to which events will be published
// when the list of active CA Roots changes. Each event payload contains the
// full list of roots.
//
// Note: topics are ordinarily defined in subscribe.proto, but this one isn't
// currently available via the Subscribe endpoint.
const EventTopicCARoots stringer = "CARoots"

// stringer is a convenience type to turn a regular string into a fmt.Stringer
// so that it can be used as a stream.Topic or stream.Subject.
type stringer string

func (s stringer) String() string { return string(s) }

type EventPayloadCARoots struct {
	CARoots structs.CARoots
}

func (e EventPayloadCARoots) Subject() stream.Subject { return stream.SubjectNone }

func (e EventPayloadCARoots) HasReadPermission(authz acl.Authorizer) bool {
	// Require `service:write` on any service in any partition and namespace.
	var authzContext acl.AuthorizerContext
	structs.WildcardEnterpriseMetaInPartition(structs.WildcardSpecifier).
		FillAuthzContext(&authzContext)

	return authz.ServiceWriteAny(&authzContext) == acl.Allow
}

// caRootsChangeEvents returns an event on EventTopicCARoots whenever the list
// of active CA Roots changes.
func caRootsChangeEvents(tx ReadTxn, changes Changes) ([]stream.Event, error) {
	var rootsChanged bool
	for _, c := range changes.Changes {
		if c.Table == tableConnectCARoots {
			rootsChanged = true
			break
		}
	}
	if !rootsChanged {
		return nil, nil
	}

	_, roots, err := caRootsTxn(tx, nil)
	if err != nil {
		return nil, err
	}

	return []stream.Event{
		{
			Topic:   EventTopicCARoots,
			Index:   changes.Index,
			Payload: EventPayloadCARoots{CARoots: roots},
		},
	}, nil
}

// caRootsSnapshot returns a stream.SnapshotFunc that provides a snapshot of
// the current active list of CA Roots.
func caRootsSnapshot(db ReadDB) stream.SnapshotFunc {
	return func(_ stream.SubscribeRequest, buf stream.SnapshotAppender) (uint64, error) {
		tx := db.ReadTxn()
		defer tx.Abort()

		idx, roots, err := caRootsTxn(tx, nil)
		if err != nil {
			return 0, err
		}

		buf.Append([]stream.Event{
			{
				Topic:   EventTopicCARoots,
				Index:   idx,
				Payload: EventPayloadCARoots{CARoots: roots},
			},
		})
		return idx, nil
	}
}
