package structs

import (
	"bytes"
	"fmt"
	"github.com/ugorji/go/codec"
)

var (
	ErrNoLeader  = fmt.Errorf("No cluster leader")
	ErrNoDCPath  = fmt.Errorf("No path to datacenter")
	ErrNoServers = fmt.Errorf("No known Consul servers")
)

type MessageType uint8

const (
	RegisterRequestType MessageType = iota
	DeregisterRequestType
)

const (
	HealthUnknown  = "unknown"
	HealthPassing  = "passing"
	HealthWarning  = "warning"
	HealthCritical = "critical"
)

// RegisterRequest is used for the Catalog.Register endpoint
// to register a node as providing a service. If no service
// is provided, the node is registered.
type RegisterRequest struct {
	Datacenter string
	Node       string
	Address    string
	Service    *NodeService
	Check      *HealthCheck
}

// DeregisterRequest is used for the Catalog.Deregister endpoint
// to deregister a node as providing a service. If no service is
// provided the entire node is deregistered.
type DeregisterRequest struct {
	Datacenter string
	Node       string
	ServiceID  string
	CheckID    string
}

// Used to return information about a node
type Node struct {
	Node    string
	Address string
}
type Nodes []Node

// Used to return information about a provided services.
// Maps service name to available tags
type Services map[string][]string

// ServiceNodesRequest is used to query the nodes of a service
type ServiceNodesRequest struct {
	Datacenter  string
	ServiceName string
	ServiceTag  string
	TagFilter   bool // Controls tag filtering
}

// ServiceNode represents a node that is part of a service
type ServiceNode struct {
	Node        string
	Address     string
	ServiceID   string
	ServiceName string
	ServiceTag  string
	ServicePort int
}
type ServiceNodes []ServiceNode

// NodeServiceRequest is used to request the services of a node
type NodeServicesRequest struct {
	Datacenter string
	Node       string
}

// NodeService is a service provided by a node
type NodeService struct {
	ID      string
	Service string
	Tag     string
	Port    int
}
type NodeServices struct {
	Address  string
	Services map[string]*NodeService
}

// HealthCheck represents a single check on a given node
type HealthCheck struct {
	Node        string
	CheckID     string // Unique per-node ID
	Name        string // Check name
	Status      string // The current check status
	Notes       string // Additional notes with the status
	ServiceID   string // optional associated service
	ServiceName string // optional service name
}
type HealthChecks []*HealthCheck

// Decode is used to decode a MsgPack encoded object
func Decode(buf []byte, out interface{}) error {
	var handle codec.MsgpackHandle
	return codec.NewDecoder(bytes.NewReader(buf), &handle).Decode(out)
}

// Encode is used to encode a MsgPack object with type prefix
func Encode(t MessageType, msg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(t))

	handle := codec.MsgpackHandle{}
	encoder := codec.NewEncoder(buf, &handle)
	err := encoder.Encode(msg)
	return buf.Bytes(), err
}
