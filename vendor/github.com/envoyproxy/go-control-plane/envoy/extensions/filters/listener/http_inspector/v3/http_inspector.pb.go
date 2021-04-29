// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/filters/listener/http_inspector/v3/http_inspector.proto

package envoy_extensions_filters_listener_http_inspector_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
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

type HttpInspector struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HttpInspector) Reset()         { *m = HttpInspector{} }
func (m *HttpInspector) String() string { return proto.CompactTextString(m) }
func (*HttpInspector) ProtoMessage()    {}
func (*HttpInspector) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec49bb714f19fafd, []int{0}
}

func (m *HttpInspector) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HttpInspector.Unmarshal(m, b)
}
func (m *HttpInspector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HttpInspector.Marshal(b, m, deterministic)
}
func (m *HttpInspector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HttpInspector.Merge(m, src)
}
func (m *HttpInspector) XXX_Size() int {
	return xxx_messageInfo_HttpInspector.Size(m)
}
func (m *HttpInspector) XXX_DiscardUnknown() {
	xxx_messageInfo_HttpInspector.DiscardUnknown(m)
}

var xxx_messageInfo_HttpInspector proto.InternalMessageInfo

func init() {
	proto.RegisterType((*HttpInspector)(nil), "envoy.extensions.filters.listener.http_inspector.v3.HttpInspector")
}

func init() {
	proto.RegisterFile("envoy/extensions/filters/listener/http_inspector/v3/http_inspector.proto", fileDescriptor_ec49bb714f19fafd)
}

var fileDescriptor_ec49bb714f19fafd = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xf2, 0x48, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0xad, 0x28, 0x49, 0xcd, 0x2b, 0xce, 0xcc, 0xcf, 0x2b, 0xd6, 0x4f, 0xcb, 0xcc,
	0x29, 0x49, 0x2d, 0x2a, 0xd6, 0xcf, 0xc9, 0x2c, 0x2e, 0x49, 0xcd, 0x4b, 0x2d, 0xd2, 0xcf, 0x28,
	0x29, 0x29, 0x88, 0xcf, 0xcc, 0x2b, 0x2e, 0x48, 0x4d, 0x2e, 0xc9, 0x2f, 0xd2, 0x2f, 0x33, 0x46,
	0x13, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x32, 0x06, 0x9b, 0xa4, 0x87, 0x30, 0x49, 0x0f,
	0x6a, 0x92, 0x1e, 0xcc, 0x24, 0x3d, 0x34, 0x7d, 0x65, 0xc6, 0x52, 0xb2, 0xa5, 0x29, 0x05, 0x89,
	0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x60, 0xeb, 0x8b, 0x4b, 0x12, 0x4b, 0x4a, 0x8b,
	0x21, 0x66, 0x4a, 0x29, 0x62, 0x48, 0x97, 0xa5, 0x16, 0x81, 0x0c, 0xcf, 0xcc, 0x4b, 0x87, 0x28,
	0x51, 0x0a, 0xe1, 0xe2, 0xf5, 0x28, 0x29, 0x29, 0xf0, 0x84, 0x99, 0x6a, 0xe5, 0x3c, 0xeb, 0x68,
	0x87, 0x9c, 0x1d, 0x97, 0x0d, 0xc4, 0x39, 0xc9, 0xf9, 0x79, 0x69, 0x99, 0xe9, 0x50, 0xa7, 0xe0,
	0x76, 0x89, 0x91, 0x1e, 0x8a, 0x21, 0x4e, 0x89, 0xbb, 0x1a, 0x4e, 0x5c, 0x64, 0x63, 0x12, 0x60,
	0xe2, 0x72, 0xcc, 0xcc, 0xd7, 0x03, 0x1b, 0x55, 0x50, 0x94, 0x5f, 0x51, 0xa9, 0x47, 0x86, 0x27,
	0x9d, 0x84, 0x50, 0xcc, 0x0e, 0x00, 0x39, 0x3b, 0x80, 0x31, 0x89, 0x0d, 0xec, 0x7e, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xdb, 0x16, 0x99, 0x5b, 0x82, 0x01, 0x00, 0x00,
}
