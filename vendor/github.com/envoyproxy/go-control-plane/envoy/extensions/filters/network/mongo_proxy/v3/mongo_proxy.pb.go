// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/filters/network/mongo_proxy/v3/mongo_proxy.proto

package envoy_extensions_filters_network_mongo_proxy_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/common/fault/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
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

type MongoProxy struct {
	StatPrefix           string         `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	AccessLog            string         `protobuf:"bytes,2,opt,name=access_log,json=accessLog,proto3" json:"access_log,omitempty"`
	Delay                *v3.FaultDelay `protobuf:"bytes,3,opt,name=delay,proto3" json:"delay,omitempty"`
	EmitDynamicMetadata  bool           `protobuf:"varint,4,opt,name=emit_dynamic_metadata,json=emitDynamicMetadata,proto3" json:"emit_dynamic_metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *MongoProxy) Reset()         { *m = MongoProxy{} }
func (m *MongoProxy) String() string { return proto.CompactTextString(m) }
func (*MongoProxy) ProtoMessage()    {}
func (*MongoProxy) Descriptor() ([]byte, []int) {
	return fileDescriptor_60ec9b7c3b00e562, []int{0}
}

func (m *MongoProxy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MongoProxy.Unmarshal(m, b)
}
func (m *MongoProxy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MongoProxy.Marshal(b, m, deterministic)
}
func (m *MongoProxy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MongoProxy.Merge(m, src)
}
func (m *MongoProxy) XXX_Size() int {
	return xxx_messageInfo_MongoProxy.Size(m)
}
func (m *MongoProxy) XXX_DiscardUnknown() {
	xxx_messageInfo_MongoProxy.DiscardUnknown(m)
}

var xxx_messageInfo_MongoProxy proto.InternalMessageInfo

func (m *MongoProxy) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *MongoProxy) GetAccessLog() string {
	if m != nil {
		return m.AccessLog
	}
	return ""
}

func (m *MongoProxy) GetDelay() *v3.FaultDelay {
	if m != nil {
		return m.Delay
	}
	return nil
}

func (m *MongoProxy) GetEmitDynamicMetadata() bool {
	if m != nil {
		return m.EmitDynamicMetadata
	}
	return false
}

func init() {
	proto.RegisterType((*MongoProxy)(nil), "envoy.extensions.filters.network.mongo_proxy.v3.MongoProxy")
}

func init() {
	proto.RegisterFile("envoy/extensions/filters/network/mongo_proxy/v3/mongo_proxy.proto", fileDescriptor_60ec9b7c3b00e562)
}

var fileDescriptor_60ec9b7c3b00e562 = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xb1, 0x4b, 0xe4, 0x40,
	0x14, 0xc6, 0x49, 0x6e, 0x6f, 0xef, 0x76, 0xb6, 0xb8, 0x23, 0xc7, 0x71, 0x61, 0x61, 0x8f, 0x68,
	0x95, 0x6a, 0x06, 0x36, 0x2b, 0x82, 0x68, 0x61, 0x58, 0x2c, 0xc4, 0x85, 0x90, 0xd2, 0x26, 0x8c,
	0xc9, 0x24, 0x0c, 0x26, 0xf3, 0x42, 0x32, 0x1b, 0x37, 0x9d, 0xa5, 0xbd, 0x9d, 0x7f, 0x8a, 0xbd,
	0x60, 0xeb, 0xbf, 0x63, 0x25, 0x93, 0x89, 0xec, 0x82, 0x6e, 0x61, 0xf7, 0x92, 0xef, 0xbd, 0xdf,
	0xf7, 0xe6, 0x7b, 0xe8, 0x94, 0x89, 0x06, 0x5a, 0xc2, 0xd6, 0x92, 0x89, 0x9a, 0x83, 0xa8, 0x49,
	0xca, 0x73, 0xc9, 0xaa, 0x9a, 0x08, 0x26, 0x6f, 0xa0, 0xba, 0x26, 0x05, 0x88, 0x0c, 0xa2, 0xb2,
	0x82, 0x75, 0x4b, 0x1a, 0x6f, 0xfb, 0x13, 0x97, 0x15, 0x48, 0xb0, 0x48, 0x87, 0xc0, 0x1b, 0x04,
	0xee, 0x11, 0xb8, 0x47, 0xe0, 0xed, 0x99, 0xc6, 0x9b, 0xcc, 0x77, 0x7a, 0xc6, 0x50, 0x14, 0x20,
	0x48, 0x4a, 0x57, 0xb9, 0x54, 0x66, 0x5d, 0xa1, 0x6d, 0x26, 0xd3, 0x55, 0x52, 0x52, 0x42, 0x85,
	0x00, 0x49, 0x65, 0x37, 0x55, 0x4b, 0x2a, 0x57, 0x75, 0x2f, 0xef, 0x7d, 0x90, 0x1b, 0x56, 0x29,
	0x3a, 0x17, 0x59, 0xdf, 0xf2, 0xaf, 0xa1, 0x39, 0x4f, 0xa8, 0x64, 0xe4, 0xbd, 0xd0, 0xc2, 0xfe,
	0xbd, 0x89, 0xd0, 0x52, 0xed, 0x18, 0xa8, 0x15, 0x2d, 0x17, 0x8d, 0x15, 0x3a, 0x2a, 0x2b, 0x96,
	0xf2, 0xb5, 0x6d, 0x38, 0x86, 0x3b, 0xf2, 0x7f, 0xbc, 0xfa, 0x83, 0xca, 0x74, 0x8c, 0x10, 0x29,
	0x2d, 0xe8, 0x24, 0x6b, 0x8a, 0x10, 0x8d, 0x63, 0x56, 0xd7, 0x51, 0x0e, 0x99, 0x6d, 0xaa, 0xc6,
	0x70, 0xa4, 0xff, 0x5c, 0x40, 0x66, 0x9d, 0xa3, 0xef, 0x09, 0xcb, 0x69, 0x6b, 0x7f, 0x73, 0x0c,
	0x77, 0x3c, 0x9b, 0xe3, 0x9d, 0x49, 0xe9, 0x87, 0x63, 0xfd, 0xde, 0xc6, 0xc3, 0x67, 0xaa, 0x58,
	0xa8, 0xd9, 0x50, 0x23, 0xac, 0x19, 0xfa, 0xcb, 0x0a, 0x2e, 0xa3, 0xa4, 0x15, 0xb4, 0xe0, 0x71,
	0x54, 0x30, 0x49, 0x13, 0x2a, 0xa9, 0x3d, 0x70, 0x0c, 0xf7, 0x67, 0xf8, 0x47, 0x89, 0x0b, 0xad,
	0x2d, 0x7b, 0xe9, 0xe8, 0xf8, 0xe1, 0xe9, 0xee, 0xff, 0x21, 0x3a, 0xd0, 0xb6, 0x31, 0x88, 0x94,
	0x67, 0xbd, 0xe5, 0xe7, 0xb7, 0x99, 0xe1, 0x4d, 0x0c, 0xfe, 0xe5, 0xe3, 0xed, 0xf3, 0xcb, 0xd0,
	0xfc, 0x6d, 0xa2, 0x13, 0x0e, 0x7a, 0x75, 0xdd, 0xf6, 0xc5, 0x7b, 0xfb, 0xbf, 0x36, 0xd0, 0x40,
	0xe5, 0x1d, 0x18, 0x57, 0xc3, 0x2e, 0x78, 0xef, 0x2d, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x39, 0x0c,
	0xaa, 0x7f, 0x02, 0x00, 0x00,
}
