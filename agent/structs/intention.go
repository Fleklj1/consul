package structs

import (
	"time"
)

// Intention defines an intention for the Connect Service Graph. This defines
// the allowed or denied behavior of a connection between two services using
// Connect.
type Intention struct {
	// ID is the UUID-based ID for the intention, always generated by Consul.
	ID string

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
	CreatedAt, UpdatedAt time.Time `mapstructure:"-"`

	RaftIndex
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

// Intentions is a list of intentions.
type Intentions []*Intention

// IndexedIntentions represents a list of intentions for RPC responses.
type IndexedIntentions struct {
	Intentions Intentions
	QueryMeta
}

// IntentionOp is the operation for a request related to intentions.
type IntentionOp string

const (
	IntentionOpCreate IntentionOp = "create"
	IntentionOpUpdate IntentionOp = "update"
	IntentionOpDelete IntentionOp = "delete"
)

// IntentionRequest is used to create, update, and delete intentions.
type IntentionRequest struct {
	// Datacenter is the target for this request.
	Datacenter string

	// Op is the type of operation being requested.
	Op IntentionOp

	// Intention is the intention.
	Intention *Intention

	// WriteRequest is a common struct containing ACL tokens and other
	// write-related common elements for requests.
	WriteRequest
}

// RequestDatacenter returns the datacenter for a given request.
func (q *IntentionRequest) RequestDatacenter() string {
	return q.Datacenter
}

// IntentionQueryRequest is used to query intentions.
type IntentionQueryRequest struct {
	// Datacenter is the target this request is intended for.
	Datacenter string

	// IntentionID is the ID of a specific intention.
	IntentionID string

	// Options for queries
	QueryOptions
}

// RequestDatacenter returns the datacenter for a given request.
func (q *IntentionQueryRequest) RequestDatacenter() string {
	return q.Datacenter
}
