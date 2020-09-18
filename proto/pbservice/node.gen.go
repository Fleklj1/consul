// Code generated by mog. DO NOT EDIT.

package pbservice

import structs "github.com/hashicorp/consul/agent/structs"

func NodeToStructs(s Node) structs.Node {
	var t structs.Node
	t.ID = s.ID
	t.Node = s.Node
	t.Address = s.Address
	t.Datacenter = s.Datacenter
	t.TaggedAddresses = s.TaggedAddresses
	t.Meta = s.Meta
	t.RaftIndex = RaftIndexToStructs(s.RaftIndex)
	return t
}
func NewNodeFromStructs(t structs.Node) Node {
	var s Node
	s.ID = t.ID
	s.Node = t.Node
	s.Address = t.Address
	s.Datacenter = t.Datacenter
	s.TaggedAddresses = t.TaggedAddresses
	s.Meta = t.Meta
	s.RaftIndex = NewRaftIndexFromStructs(t.RaftIndex)
	return s
}
func NodeServiceToStructs(s NodeService) structs.NodeService {
	var t structs.NodeService
	t.Kind = s.Kind
	t.ID = s.ID
	t.Service = s.Service
	t.Tags = s.Tags
	t.Address = s.Address
	t.TaggedAddresses = MapStringServiceAddressToStructs(s.TaggedAddresses)
	t.Meta = s.Meta
	t.Port = int(s.Port)
	t.Weights = WeightsPtrToStructs(s.Weights)
	t.EnableTagOverride = s.EnableTagOverride
	t.Proxy = ConnectProxyConfigToStructs(s.Proxy)
	t.Connect = ServiceConnectToStructs(s.Connect)
	t.LocallyRegisteredAsSidecar = s.LocallyRegisteredAsSidecar
	t.EnterpriseMeta = EnterpriseMetaToStructs(s.EnterpriseMeta)
	t.RaftIndex = RaftIndexToStructs(s.RaftIndex)
	return t
}
func NewNodeServiceFromStructs(t structs.NodeService) NodeService {
	var s NodeService
	s.Kind = t.Kind
	s.ID = t.ID
	s.Service = t.Service
	s.Tags = t.Tags
	s.Address = t.Address
	s.TaggedAddresses = NewMapStringServiceAddressFromStructs(t.TaggedAddresses)
	s.Meta = t.Meta
	s.Port = int32(t.Port)
	s.Weights = NewWeightsPtrFromStructs(t.Weights)
	s.EnableTagOverride = t.EnableTagOverride
	s.Proxy = NewConnectProxyConfigFromStructs(t.Proxy)
	s.Connect = NewServiceConnectFromStructs(t.Connect)
	s.LocallyRegisteredAsSidecar = t.LocallyRegisteredAsSidecar
	s.EnterpriseMeta = NewEnterpriseMetaFromStructs(t.EnterpriseMeta)
	s.RaftIndex = NewRaftIndexFromStructs(t.RaftIndex)
	return s
}
