// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/filters/network/thrift_proxy/v3/route.proto

package envoy_extensions_filters_network_thrift_proxy_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RouteConfiguration struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Routes               []*Route `protobuf:"bytes,2,rep,name=routes,proto3" json:"routes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouteConfiguration) Reset()         { *m = RouteConfiguration{} }
func (m *RouteConfiguration) String() string { return proto.CompactTextString(m) }
func (*RouteConfiguration) ProtoMessage()    {}
func (*RouteConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{0}
}

func (m *RouteConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteConfiguration.Unmarshal(m, b)
}
func (m *RouteConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteConfiguration.Marshal(b, m, deterministic)
}
func (m *RouteConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteConfiguration.Merge(m, src)
}
func (m *RouteConfiguration) XXX_Size() int {
	return xxx_messageInfo_RouteConfiguration.Size(m)
}
func (m *RouteConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_RouteConfiguration proto.InternalMessageInfo

func (m *RouteConfiguration) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RouteConfiguration) GetRoutes() []*Route {
	if m != nil {
		return m.Routes
	}
	return nil
}

type Route struct {
	Match                *RouteMatch  `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
	Route                *RouteAction `protobuf:"bytes,2,opt,name=route,proto3" json:"route,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{1}
}

func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (m *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(m, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetMatch() *RouteMatch {
	if m != nil {
		return m.Match
	}
	return nil
}

func (m *Route) GetRoute() *RouteAction {
	if m != nil {
		return m.Route
	}
	return nil
}

type RouteMatch struct {
	// Types that are valid to be assigned to MatchSpecifier:
	//	*RouteMatch_MethodName
	//	*RouteMatch_ServiceName
	MatchSpecifier       isRouteMatch_MatchSpecifier `protobuf_oneof:"match_specifier"`
	Invert               bool                        `protobuf:"varint,3,opt,name=invert,proto3" json:"invert,omitempty"`
	Headers              []*v3.HeaderMatcher         `protobuf:"bytes,4,rep,name=headers,proto3" json:"headers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *RouteMatch) Reset()         { *m = RouteMatch{} }
func (m *RouteMatch) String() string { return proto.CompactTextString(m) }
func (*RouteMatch) ProtoMessage()    {}
func (*RouteMatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{2}
}

func (m *RouteMatch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteMatch.Unmarshal(m, b)
}
func (m *RouteMatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteMatch.Marshal(b, m, deterministic)
}
func (m *RouteMatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteMatch.Merge(m, src)
}
func (m *RouteMatch) XXX_Size() int {
	return xxx_messageInfo_RouteMatch.Size(m)
}
func (m *RouteMatch) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteMatch.DiscardUnknown(m)
}

var xxx_messageInfo_RouteMatch proto.InternalMessageInfo

type isRouteMatch_MatchSpecifier interface {
	isRouteMatch_MatchSpecifier()
}

type RouteMatch_MethodName struct {
	MethodName string `protobuf:"bytes,1,opt,name=method_name,json=methodName,proto3,oneof"`
}

type RouteMatch_ServiceName struct {
	ServiceName string `protobuf:"bytes,2,opt,name=service_name,json=serviceName,proto3,oneof"`
}

func (*RouteMatch_MethodName) isRouteMatch_MatchSpecifier() {}

func (*RouteMatch_ServiceName) isRouteMatch_MatchSpecifier() {}

func (m *RouteMatch) GetMatchSpecifier() isRouteMatch_MatchSpecifier {
	if m != nil {
		return m.MatchSpecifier
	}
	return nil
}

func (m *RouteMatch) GetMethodName() string {
	if x, ok := m.GetMatchSpecifier().(*RouteMatch_MethodName); ok {
		return x.MethodName
	}
	return ""
}

func (m *RouteMatch) GetServiceName() string {
	if x, ok := m.GetMatchSpecifier().(*RouteMatch_ServiceName); ok {
		return x.ServiceName
	}
	return ""
}

func (m *RouteMatch) GetInvert() bool {
	if m != nil {
		return m.Invert
	}
	return false
}

func (m *RouteMatch) GetHeaders() []*v3.HeaderMatcher {
	if m != nil {
		return m.Headers
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*RouteMatch) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RouteMatch_MethodName)(nil),
		(*RouteMatch_ServiceName)(nil),
	}
}

type RouteAction struct {
	// Types that are valid to be assigned to ClusterSpecifier:
	//	*RouteAction_Cluster
	//	*RouteAction_WeightedClusters
	//	*RouteAction_ClusterHeader
	ClusterSpecifier     isRouteAction_ClusterSpecifier `protobuf_oneof:"cluster_specifier"`
	MetadataMatch        *v31.Metadata                  `protobuf:"bytes,3,opt,name=metadata_match,json=metadataMatch,proto3" json:"metadata_match,omitempty"`
	RateLimits           []*v3.RateLimit                `protobuf:"bytes,4,rep,name=rate_limits,json=rateLimits,proto3" json:"rate_limits,omitempty"`
	StripServiceName     bool                           `protobuf:"varint,5,opt,name=strip_service_name,json=stripServiceName,proto3" json:"strip_service_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *RouteAction) Reset()         { *m = RouteAction{} }
func (m *RouteAction) String() string { return proto.CompactTextString(m) }
func (*RouteAction) ProtoMessage()    {}
func (*RouteAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{3}
}

func (m *RouteAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteAction.Unmarshal(m, b)
}
func (m *RouteAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteAction.Marshal(b, m, deterministic)
}
func (m *RouteAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteAction.Merge(m, src)
}
func (m *RouteAction) XXX_Size() int {
	return xxx_messageInfo_RouteAction.Size(m)
}
func (m *RouteAction) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteAction.DiscardUnknown(m)
}

var xxx_messageInfo_RouteAction proto.InternalMessageInfo

type isRouteAction_ClusterSpecifier interface {
	isRouteAction_ClusterSpecifier()
}

type RouteAction_Cluster struct {
	Cluster string `protobuf:"bytes,1,opt,name=cluster,proto3,oneof"`
}

type RouteAction_WeightedClusters struct {
	WeightedClusters *WeightedCluster `protobuf:"bytes,2,opt,name=weighted_clusters,json=weightedClusters,proto3,oneof"`
}

type RouteAction_ClusterHeader struct {
	ClusterHeader string `protobuf:"bytes,6,opt,name=cluster_header,json=clusterHeader,proto3,oneof"`
}

func (*RouteAction_Cluster) isRouteAction_ClusterSpecifier() {}

func (*RouteAction_WeightedClusters) isRouteAction_ClusterSpecifier() {}

func (*RouteAction_ClusterHeader) isRouteAction_ClusterSpecifier() {}

func (m *RouteAction) GetClusterSpecifier() isRouteAction_ClusterSpecifier {
	if m != nil {
		return m.ClusterSpecifier
	}
	return nil
}

func (m *RouteAction) GetCluster() string {
	if x, ok := m.GetClusterSpecifier().(*RouteAction_Cluster); ok {
		return x.Cluster
	}
	return ""
}

func (m *RouteAction) GetWeightedClusters() *WeightedCluster {
	if x, ok := m.GetClusterSpecifier().(*RouteAction_WeightedClusters); ok {
		return x.WeightedClusters
	}
	return nil
}

func (m *RouteAction) GetClusterHeader() string {
	if x, ok := m.GetClusterSpecifier().(*RouteAction_ClusterHeader); ok {
		return x.ClusterHeader
	}
	return ""
}

func (m *RouteAction) GetMetadataMatch() *v31.Metadata {
	if m != nil {
		return m.MetadataMatch
	}
	return nil
}

func (m *RouteAction) GetRateLimits() []*v3.RateLimit {
	if m != nil {
		return m.RateLimits
	}
	return nil
}

func (m *RouteAction) GetStripServiceName() bool {
	if m != nil {
		return m.StripServiceName
	}
	return false
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*RouteAction) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RouteAction_Cluster)(nil),
		(*RouteAction_WeightedClusters)(nil),
		(*RouteAction_ClusterHeader)(nil),
	}
}

type WeightedCluster struct {
	Clusters             []*WeightedCluster_ClusterWeight `protobuf:"bytes,1,rep,name=clusters,proto3" json:"clusters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *WeightedCluster) Reset()         { *m = WeightedCluster{} }
func (m *WeightedCluster) String() string { return proto.CompactTextString(m) }
func (*WeightedCluster) ProtoMessage()    {}
func (*WeightedCluster) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{4}
}

func (m *WeightedCluster) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeightedCluster.Unmarshal(m, b)
}
func (m *WeightedCluster) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeightedCluster.Marshal(b, m, deterministic)
}
func (m *WeightedCluster) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightedCluster.Merge(m, src)
}
func (m *WeightedCluster) XXX_Size() int {
	return xxx_messageInfo_WeightedCluster.Size(m)
}
func (m *WeightedCluster) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightedCluster.DiscardUnknown(m)
}

var xxx_messageInfo_WeightedCluster proto.InternalMessageInfo

func (m *WeightedCluster) GetClusters() []*WeightedCluster_ClusterWeight {
	if m != nil {
		return m.Clusters
	}
	return nil
}

type WeightedCluster_ClusterWeight struct {
	Name                 string                `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Weight               *wrappers.UInt32Value `protobuf:"bytes,2,opt,name=weight,proto3" json:"weight,omitempty"`
	MetadataMatch        *v31.Metadata         `protobuf:"bytes,3,opt,name=metadata_match,json=metadataMatch,proto3" json:"metadata_match,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *WeightedCluster_ClusterWeight) Reset()         { *m = WeightedCluster_ClusterWeight{} }
func (m *WeightedCluster_ClusterWeight) String() string { return proto.CompactTextString(m) }
func (*WeightedCluster_ClusterWeight) ProtoMessage()    {}
func (*WeightedCluster_ClusterWeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_77b9ee5a656d870d, []int{4, 0}
}

func (m *WeightedCluster_ClusterWeight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeightedCluster_ClusterWeight.Unmarshal(m, b)
}
func (m *WeightedCluster_ClusterWeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeightedCluster_ClusterWeight.Marshal(b, m, deterministic)
}
func (m *WeightedCluster_ClusterWeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeightedCluster_ClusterWeight.Merge(m, src)
}
func (m *WeightedCluster_ClusterWeight) XXX_Size() int {
	return xxx_messageInfo_WeightedCluster_ClusterWeight.Size(m)
}
func (m *WeightedCluster_ClusterWeight) XXX_DiscardUnknown() {
	xxx_messageInfo_WeightedCluster_ClusterWeight.DiscardUnknown(m)
}

var xxx_messageInfo_WeightedCluster_ClusterWeight proto.InternalMessageInfo

func (m *WeightedCluster_ClusterWeight) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *WeightedCluster_ClusterWeight) GetWeight() *wrappers.UInt32Value {
	if m != nil {
		return m.Weight
	}
	return nil
}

func (m *WeightedCluster_ClusterWeight) GetMetadataMatch() *v31.Metadata {
	if m != nil {
		return m.MetadataMatch
	}
	return nil
}

func init() {
	proto.RegisterType((*RouteConfiguration)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.RouteConfiguration")
	proto.RegisterType((*Route)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.Route")
	proto.RegisterType((*RouteMatch)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.RouteMatch")
	proto.RegisterType((*RouteAction)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.RouteAction")
	proto.RegisterType((*WeightedCluster)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.WeightedCluster")
	proto.RegisterType((*WeightedCluster_ClusterWeight)(nil), "envoy.extensions.filters.network.thrift_proxy.v3.WeightedCluster.ClusterWeight")
}

func init() {
	proto.RegisterFile("envoy/extensions/filters/network/thrift_proxy/v3/route.proto", fileDescriptor_77b9ee5a656d870d)
}

var fileDescriptor_77b9ee5a656d870d = []byte{
	// 803 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xcd, 0x6e, 0xf3, 0x44,
	0x14, 0xad, 0x9d, 0x26, 0x0d, 0x37, 0xb4, 0x4d, 0x67, 0x51, 0xa2, 0x02, 0x25, 0x4d, 0x59, 0x44,
	0xa8, 0xb2, 0x4b, 0xb2, 0xa8, 0x88, 0xda, 0xa0, 0x24, 0x05, 0x05, 0x41, 0x69, 0x31, 0xa2, 0x6c,
	0x40, 0xd6, 0xd4, 0x99, 0x24, 0x16, 0x89, 0xc7, 0x1a, 0x8f, 0x9d, 0x76, 0xc7, 0x82, 0x05, 0x74,
	0xc9, 0x92, 0x47, 0x61, 0x0f, 0x62, 0xcb, 0x13, 0xf0, 0x06, 0xac, 0x51, 0x57, 0x68, 0x7e, 0x9c,
	0x9f, 0x56, 0xfd, 0xa4, 0xa4, 0xdf, 0x2a, 0x93, 0x99, 0x7b, 0xcf, 0x3d, 0xf7, 0x9c, 0x3b, 0x63,
	0x38, 0x25, 0x41, 0x42, 0xef, 0x6c, 0x72, 0xcb, 0x49, 0x10, 0xf9, 0x34, 0x88, 0xec, 0xbe, 0x3f,
	0xe2, 0x84, 0x45, 0x76, 0x40, 0xf8, 0x84, 0xb2, 0x1f, 0x6c, 0x3e, 0x64, 0x7e, 0x9f, 0xbb, 0x21,
	0xa3, 0xb7, 0x77, 0x76, 0x52, 0xb7, 0x19, 0x8d, 0x39, 0xb1, 0x42, 0x46, 0x39, 0x45, 0xc7, 0x32,
	0xdb, 0x9a, 0x65, 0x5b, 0x3a, 0xdb, 0xd2, 0xd9, 0xd6, 0x7c, 0xb6, 0x95, 0xd4, 0xf7, 0xde, 0x53,
	0xf5, 0x3c, 0x1a, 0xf4, 0xfd, 0x81, 0xed, 0x51, 0x46, 0x04, 0xe6, 0x0d, 0x8e, 0x34, 0xe4, 0xde,
	0xd1, 0x42, 0x80, 0x2c, 0x36, 0xad, 0xea, 0x7a, 0x74, 0x1c, 0xd2, 0x80, 0x04, 0x3c, 0xd2, 0xd1,
	0xfb, 0x03, 0x4a, 0x07, 0x23, 0x62, 0xcb, 0x7f, 0x37, 0x71, 0xdf, 0x9e, 0x30, 0x1c, 0x86, 0x82,
	0x80, 0x3a, 0x7f, 0x37, 0xee, 0x85, 0xd8, 0xc6, 0x41, 0x40, 0x39, 0xe6, 0xb2, 0xbd, 0x88, 0x63,
	0x1e, 0xa7, 0xc7, 0x07, 0x4f, 0x8e, 0x13, 0xc2, 0x44, 0x23, 0x7e, 0x30, 0xd0, 0x21, 0x6f, 0x25,
	0x78, 0xe4, 0xf7, 0xb0, 0x60, 0xa1, 0x17, 0xea, 0xa0, 0xf2, 0xa7, 0x01, 0xc8, 0x11, 0xac, 0x3a,
	0x92, 0x6a, 0xcc, 0x24, 0x02, 0x42, 0xb0, 0x1e, 0xe0, 0x31, 0x29, 0x19, 0x65, 0xa3, 0xfa, 0x86,
	0x23, 0xd7, 0xe8, 0x12, 0x72, 0x92, 0x7f, 0x54, 0x32, 0xcb, 0x99, 0x6a, 0xa1, 0x76, 0x62, 0x2d,
	0xab, 0x9b, 0x25, 0x2b, 0x39, 0x1a, 0xa6, 0xf1, 0xf9, 0x6f, 0x7f, 0xfc, 0xbc, 0xff, 0x29, 0x9c,
	0x2b, 0x18, 0xa5, 0x95, 0x86, 0x78, 0x06, 0xa1, 0x86, 0x47, 0xe1, 0x10, 0x7f, 0x68, 0x3d, 0x65,
	0x5c, 0xf9, 0xc9, 0x84, 0xac, 0xdc, 0x46, 0xdf, 0x41, 0x76, 0x8c, 0xb9, 0x37, 0x94, 0xe4, 0x0b,
	0xb5, 0xd3, 0x15, 0x69, 0x5e, 0x08, 0x8c, 0x76, 0xfe, 0xa1, 0x9d, 0xbd, 0x37, 0xcc, 0xa2, 0xe1,
	0x28, 0x50, 0xf4, 0x3d, 0x64, 0x25, 0xfd, 0x92, 0x29, 0xd1, 0xcf, 0x56, 0x44, 0x6f, 0x79, 0x82,
	0xf5, 0x3c, 0xbc, 0x44, 0x6d, 0x34, 0x85, 0x26, 0x1f, 0xc1, 0xc9, 0x8a, 0x9a, 0x54, 0x7e, 0x31,
	0x01, 0x66, 0xf4, 0xd1, 0x01, 0x14, 0xc6, 0x84, 0x0f, 0x69, 0xcf, 0x9d, 0xd9, 0xd9, 0x5d, 0x73,
	0x40, 0x6d, 0x7e, 0x29, 0x6c, 0x3d, 0x84, 0x37, 0x23, 0xc2, 0x12, 0xdf, 0x23, 0x2a, 0xc6, 0xd4,
	0x31, 0x05, 0xbd, 0x2b, 0x83, 0x76, 0x21, 0xe7, 0x07, 0x09, 0x61, 0xbc, 0x94, 0x29, 0x1b, 0xd5,
	0xbc, 0xa3, 0xff, 0xa1, 0x26, 0x6c, 0x0c, 0x09, 0xee, 0x11, 0x16, 0x95, 0xd6, 0xe5, 0x50, 0xbc,
	0x6f, 0x2d, 0x30, 0x57, 0xd7, 0x2c, 0xa9, 0x5b, 0x5d, 0x19, 0x25, 0x49, 0x11, 0xe6, 0xa4, 0x49,
	0x8d, 0x8e, 0x68, 0xb7, 0xa9, 0xef, 0xef, 0xf2, 0xed, 0x2a, 0x8f, 0x76, 0x61, 0x5b, 0x7a, 0xe3,
	0x46, 0x21, 0xf1, 0xfc, 0xbe, 0x4f, 0x18, 0xca, 0xfc, 0xd7, 0x36, 0x2a, 0xff, 0x66, 0xa0, 0x30,
	0x27, 0x36, 0x3a, 0x84, 0x0d, 0x6f, 0x14, 0x47, 0x9c, 0x30, 0x25, 0x44, 0x7b, 0xe3, 0xa1, 0xbd,
	0xce, 0xcc, 0xb2, 0xd1, 0x5d, 0x73, 0xd2, 0x13, 0x14, 0xc2, 0xce, 0x84, 0xf8, 0x83, 0x21, 0x27,
	0x3d, 0x57, 0xef, 0x45, 0xda, 0xeb, 0xd6, 0xf2, 0x5e, 0x7f, 0xab, 0xa1, 0x3a, 0x0a, 0xa9, 0xbb,
	0xe6, 0x14, 0x27, 0x8b, 0x5b, 0x11, 0x3a, 0x86, 0x2d, 0x5d, 0xc8, 0x55, 0xb2, 0x94, 0x72, 0x8f,
	0xd9, 0x6d, 0xea, 0x00, 0xa5, 0x22, 0xfa, 0x04, 0xb6, 0xc6, 0x84, 0xe3, 0x1e, 0xe6, 0xd8, 0x55,
	0xa3, 0x9e, 0x91, 0x04, 0xf7, 0x17, 0xc5, 0x17, 0xef, 0x92, 0x20, 0x71, 0xa1, 0x63, 0x9d, 0xcd,
	0x34, 0x4b, 0x0d, 0x47, 0x0b, 0x0a, 0x0c, 0x73, 0xe2, 0x8e, 0xfc, 0xb1, 0xcf, 0x53, 0x03, 0xcb,
	0xcf, 0x18, 0xe8, 0x60, 0x4e, 0xbe, 0x10, 0x81, 0x0e, 0xb0, 0x74, 0x19, 0xa1, 0x23, 0x40, 0x11,
	0x67, 0x7e, 0xe8, 0x2e, 0x8c, 0x50, 0x56, 0xce, 0x48, 0x51, 0x9e, 0x7c, 0x3d, 0x9b, 0xa2, 0xc6,
	0xb9, 0x70, 0xfb, 0x63, 0x38, 0x5b, 0xd1, 0x6d, 0x7d, 0x67, 0x4a, 0xb0, 0x93, 0xea, 0xf5, 0xc8,
	0xf0, 0x7f, 0x32, 0xb0, 0xfd, 0x48, 0x71, 0x14, 0x43, 0x7e, 0x6a, 0xa3, 0x21, 0x3b, 0xbc, 0x7c,
	0xb1, 0x8d, 0x96, 0xfe, 0x55, 0xdb, 0xf2, 0x12, 0xff, 0x6a, 0x98, 0x79, 0xc3, 0x99, 0x96, 0xda,
	0xbb, 0x37, 0x61, 0x73, 0x21, 0x0a, 0xbd, 0x3d, 0xff, 0xa4, 0x4e, 0xcd, 0xd5, 0x6f, 0xeb, 0x19,
	0xe4, 0xd4, 0x5c, 0xe8, 0x51, 0x7b, 0xc7, 0x52, 0x9f, 0x04, 0x2b, 0xfd, 0x24, 0x58, 0xdf, 0x7c,
	0x16, 0xf0, 0x7a, 0xed, 0x1a, 0x8f, 0x62, 0x22, 0x93, 0x3f, 0x30, 0xab, 0x86, 0xa3, 0x93, 0x5e,
	0xd3, 0x40, 0x34, 0xae, 0x85, 0x3f, 0x5f, 0xc1, 0xe5, 0xf2, 0xfe, 0xbc, 0x52, 0xa0, 0x46, 0x57,
	0xe0, 0x76, 0xa0, 0xf5, 0x62, 0xdc, 0xf6, 0xf5, 0xef, 0x3f, 0xfe, 0xf5, 0x77, 0xce, 0x2c, 0x9a,
	0xd0, 0xf4, 0xa9, 0x6a, 0x4e, 0x65, 0x2c, 0x6b, 0x69, 0x5b, 0xbd, 0x92, 0x57, 0x42, 0xde, 0x2b,
	0xe3, 0x26, 0x27, 0x75, 0xae, 0xff, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x64, 0xc8, 0x3c, 0x59, 0x49,
	0x08, 0x00, 0x00,
}
