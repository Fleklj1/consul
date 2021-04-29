// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/filters/http/dynamo/v3/dynamo.proto

package envoy_extensions_filters_http_dynamo_v3

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

type Dynamo struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dynamo) Reset()         { *m = Dynamo{} }
func (m *Dynamo) String() string { return proto.CompactTextString(m) }
func (*Dynamo) ProtoMessage()    {}
func (*Dynamo) Descriptor() ([]byte, []int) {
	return fileDescriptor_79057240c5b18ac4, []int{0}
}

func (m *Dynamo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dynamo.Unmarshal(m, b)
}
func (m *Dynamo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dynamo.Marshal(b, m, deterministic)
}
func (m *Dynamo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dynamo.Merge(m, src)
}
func (m *Dynamo) XXX_Size() int {
	return xxx_messageInfo_Dynamo.Size(m)
}
func (m *Dynamo) XXX_DiscardUnknown() {
	xxx_messageInfo_Dynamo.DiscardUnknown(m)
}

var xxx_messageInfo_Dynamo proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Dynamo)(nil), "envoy.extensions.filters.http.dynamo.v3.Dynamo")
}

func init() {
	proto.RegisterFile("envoy/extensions/filters/http/dynamo/v3/dynamo.proto", fileDescriptor_79057240c5b18ac4)
}

var fileDescriptor_79057240c5b18ac4 = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x49, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0xad, 0x28, 0x49, 0xcd, 0x2b, 0xce, 0xcc, 0xcf, 0x2b, 0xd6, 0x4f, 0xcb, 0xcc,
	0x29, 0x49, 0x2d, 0x2a, 0xd6, 0xcf, 0x28, 0x29, 0x29, 0xd0, 0x4f, 0xa9, 0xcc, 0x4b, 0xcc, 0xcd,
	0xd7, 0x2f, 0x33, 0x86, 0xb2, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xd4, 0xc1, 0xba, 0xf4,
	0x10, 0xba, 0xf4, 0xa0, 0xba, 0xf4, 0x40, 0xba, 0xf4, 0xa0, 0x6a, 0xcb, 0x8c, 0xa5, 0x64, 0x4b,
	0x53, 0x0a, 0x12, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0xc0, 0xc6, 0x17, 0x97, 0x24,
	0x96, 0x94, 0x16, 0x43, 0xcc, 0x91, 0x52, 0xc4, 0x90, 0x2e, 0x4b, 0x2d, 0x02, 0x19, 0x98, 0x99,
	0x97, 0x0e, 0x51, 0xa2, 0x64, 0xc5, 0xc5, 0xe6, 0x02, 0x36, 0xce, 0xca, 0x60, 0xd6, 0xd1, 0x0e,
	0x39, 0x6d, 0x2e, 0x4d, 0x88, 0xdd, 0xc9, 0xf9, 0x79, 0x69, 0x99, 0xe9, 0x50, 0x7b, 0x51, 0xad,
	0x35, 0xd2, 0x83, 0xe8, 0x70, 0xf2, 0xdb, 0xd5, 0x70, 0xe2, 0x22, 0x1b, 0x93, 0x00, 0x13, 0x97,
	0x69, 0x66, 0xbe, 0x1e, 0x58, 0x5f, 0x41, 0x51, 0x7e, 0x45, 0xa5, 0x1e, 0x91, 0xce, 0x77, 0xe2,
	0x86, 0x18, 0x14, 0x00, 0x72, 0x49, 0x00, 0x63, 0x12, 0x1b, 0xd8, 0x49, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x6b, 0xbd, 0x20, 0xc4, 0x35, 0x01, 0x00, 0x00,
}
