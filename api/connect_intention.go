package api

import (
	"time"
)

// Intention defines an intention for the Connect Service Graph. This defines
// the allowed or denied behavior of a connection between two services using
// Connect.
type Intention struct {
	// ID is the UUID-based ID for the intention, always generated by Consul.
	ID string

	// Description is a human-friendly description of this intention.
	// It is opaque to Consul and is only stored and transferred in API
	// requests.
	Description string

	// SourceNS, SourceName are the namespace and name, respectively, of
	// the source service. Either of these may be the wildcard "*", but only
	// the full value can be a wildcard. Partial wildcards are not allowed.
	// The source may also be a non-Consul service, as specified by SourceType.
	//
	// DestinationNS, DestinationName is the same, but for the destination
	// service. The same rules apply. The destination is always a Consul
	// service.
	SourceNS, SourceName           string
	DestinationNS, DestinationName string

	// SourceType is the type of the value for the source.
	SourceType IntentionSourceType

	// Action is whether this is a whitelist or blacklist intention.
	Action IntentionAction

	// DefaultAddr, DefaultPort of the local listening proxy (if any) to
	// make this connection.
	DefaultAddr string
	DefaultPort int

	// Meta is arbitrary metadata associated with the intention. This is
	// opaque to Consul but is served in API responses.
	Meta map[string]string

	// CreatedAt and UpdatedAt keep track of when this record was created
	// or modified.
	CreatedAt, UpdatedAt time.Time

	CreateIndex uint64
	ModifyIndex uint64
}

// IntentionAction is the action that the intention represents. This
// can be "allow" or "deny" to whitelist or blacklist intentions.
type IntentionAction string

const (
	IntentionActionAllow IntentionAction = "allow"
	IntentionActionDeny  IntentionAction = "deny"
)

// IntentionSourceType is the type of the source within an intention.
type IntentionSourceType string

const (
	// IntentionSourceConsul is a service within the Consul catalog.
	IntentionSourceConsul IntentionSourceType = "consul"
)

// Intentions returns the list of intentions.
func (h *Connect) Intentions(q *QueryOptions) ([]*Intention, *QueryMeta, error) {
	r := h.c.newRequest("GET", "/v1/connect/intentions")
	r.setQueryOptions(q)
	rtt, resp, err := requireOK(h.c.doRequest(r))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	qm := &QueryMeta{}
	parseQueryMeta(resp, qm)
	qm.RequestTime = rtt

	var out []*Intention
	if err := decodeBody(resp, &out); err != nil {
		return nil, nil, err
	}
	return out, qm, nil
}

// IntentionCreate will create a new intention. The ID in the given
// structure must be empty and a generate ID will be returned on
// success.
func (c *Connect) IntentionCreate(ixn *Intention, q *WriteOptions) (string, *WriteMeta, error) {
	r := c.c.newRequest("POST", "/v1/connect/intentions")
	r.setWriteOptions(q)
	r.obj = ixn
	rtt, resp, err := requireOK(c.c.doRequest(r))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	wm := &WriteMeta{}
	wm.RequestTime = rtt

	var out struct{ ID string }
	if err := decodeBody(resp, &out); err != nil {
		return "", nil, err
	}
	return out.ID, wm, nil
}
