// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/extensions/filters/network/local_ratelimit/v3/local_rate_limit.proto

package envoy_extensions_filters_network_local_ratelimit_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	v31 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
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

type LocalRateLimit struct {
	StatPrefix           string                  `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	TokenBucket          *v3.TokenBucket         `protobuf:"bytes,2,opt,name=token_bucket,json=tokenBucket,proto3" json:"token_bucket,omitempty"`
	RuntimeEnabled       *v31.RuntimeFeatureFlag `protobuf:"bytes,3,opt,name=runtime_enabled,json=runtimeEnabled,proto3" json:"runtime_enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *LocalRateLimit) Reset()         { *m = LocalRateLimit{} }
func (m *LocalRateLimit) String() string { return proto.CompactTextString(m) }
func (*LocalRateLimit) ProtoMessage()    {}
func (*LocalRateLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_ffceb0cea724f411, []int{0}
}

func (m *LocalRateLimit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocalRateLimit.Unmarshal(m, b)
}
func (m *LocalRateLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocalRateLimit.Marshal(b, m, deterministic)
}
func (m *LocalRateLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocalRateLimit.Merge(m, src)
}
func (m *LocalRateLimit) XXX_Size() int {
	return xxx_messageInfo_LocalRateLimit.Size(m)
}
func (m *LocalRateLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_LocalRateLimit.DiscardUnknown(m)
}

var xxx_messageInfo_LocalRateLimit proto.InternalMessageInfo

func (m *LocalRateLimit) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *LocalRateLimit) GetTokenBucket() *v3.TokenBucket {
	if m != nil {
		return m.TokenBucket
	}
	return nil
}

func (m *LocalRateLimit) GetRuntimeEnabled() *v31.RuntimeFeatureFlag {
	if m != nil {
		return m.RuntimeEnabled
	}
	return nil
}

func init() {
	proto.RegisterType((*LocalRateLimit)(nil), "envoy.extensions.filters.network.local_ratelimit.v3.LocalRateLimit")
}

func init() {
	proto.RegisterFile("envoy/extensions/filters/network/local_ratelimit/v3/local_rate_limit.proto", fileDescriptor_ffceb0cea724f411)
}

var fileDescriptor_ffceb0cea724f411 = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x3f, 0xaf, 0xd3, 0x30,
	0x14, 0xc5, 0x95, 0x00, 0x0f, 0x70, 0xd1, 0xe3, 0x29, 0x0c, 0x54, 0x91, 0x80, 0xc0, 0x94, 0xc9,
	0x96, 0x5e, 0x36, 0x36, 0x82, 0xde, 0x43, 0x7a, 0xea, 0x50, 0x22, 0xf6, 0xc8, 0x49, 0x6f, 0x8b,
	0x55, 0xd7, 0x8e, 0x9c, 0x9b, 0xd0, 0x6e, 0x8c, 0x88, 0x8f, 0xc0, 0xc8, 0xc7, 0x60, 0x47, 0x62,
	0xe5, 0xeb, 0x74, 0x42, 0xfe, 0x83, 0xda, 0x02, 0x13, 0x5b, 0x74, 0xcf, 0xf1, 0xcf, 0xc7, 0xf7,
	0x84, 0xdc, 0x80, 0x1a, 0xf5, 0x8e, 0xc1, 0x16, 0x41, 0xf5, 0x42, 0xab, 0x9e, 0x2d, 0x85, 0x44,
	0x30, 0x3d, 0x53, 0x80, 0x1f, 0xb4, 0x59, 0x33, 0xa9, 0x5b, 0x2e, 0x6b, 0xc3, 0x11, 0xa4, 0xd8,
	0x08, 0x64, 0x63, 0x71, 0x34, 0xaa, 0xdd, 0x8c, 0x76, 0x46, 0xa3, 0x4e, 0x0a, 0xc7, 0xa2, 0x07,
	0x16, 0x0d, 0x2c, 0x1a, 0x58, 0xf4, 0x0f, 0x16, 0x1d, 0x8b, 0xf4, 0x99, 0x0f, 0xd0, 0x6a, 0xb5,
	0x14, 0x2b, 0xd6, 0x6a, 0x03, 0xf6, 0x86, 0x86, 0xf7, 0xe0, 0xa9, 0x69, 0xe6, 0x0d, 0xb8, 0xeb,
	0x9c, 0x82, 0x7a, 0x0d, 0xaa, 0x6e, 0x86, 0x76, 0x0d, 0xe1, 0xde, 0xf4, 0xc9, 0xb0, 0xe8, 0x38,
	0xe3, 0x4a, 0x69, 0xe4, 0xe8, 0xde, 0xd0, 0x23, 0xc7, 0xa1, 0x0f, 0xf2, 0xf3, 0xbf, 0xe4, 0x11,
	0x8c, 0xcd, 0x27, 0xd4, 0x2a, 0x58, 0x1e, 0x8f, 0x5c, 0x8a, 0x05, 0x47, 0x60, 0xbf, 0x3f, 0xbc,
	0xf0, 0xe2, 0x6b, 0x4c, 0xce, 0x67, 0x36, 0x74, 0xc5, 0x11, 0x66, 0x36, 0x73, 0x92, 0x93, 0x89,
	0xc5, 0xd7, 0x9d, 0x81, 0xa5, 0xd8, 0x4e, 0xa3, 0x2c, 0xca, 0xef, 0x97, 0x77, 0xf7, 0xe5, 0x6d,
	0x13, 0x67, 0x51, 0x45, 0xac, 0x36, 0x77, 0x52, 0xf2, 0x86, 0x3c, 0x38, 0x4e, 0x3b, 0x8d, 0xb3,
	0x28, 0x9f, 0x5c, 0xa6, 0xd4, 0xaf, 0xc9, 0x3e, 0x88, 0x8e, 0x05, 0x7d, 0x67, 0x2d, 0xa5, 0x73,
	0x94, 0xf7, 0xf6, 0xe5, 0x9d, 0xcf, 0x51, 0x7c, 0x11, 0x55, 0x13, 0x3c, 0x8c, 0x93, 0xb7, 0xe4,
	0xa1, 0x19, 0x14, 0x8a, 0x0d, 0xd4, 0xa0, 0x78, 0x23, 0x61, 0x31, 0xbd, 0xe5, 0x58, 0x79, 0x60,
	0xf9, 0xed, 0x51, 0xbb, 0x3d, 0x8b, 0xac, 0xbc, 0xf9, 0x1a, 0x38, 0x0e, 0x06, 0xae, 0x25, 0x5f,
	0x55, 0xe7, 0x01, 0x70, 0xe5, 0xcf, 0xbf, 0xbc, 0xf9, 0xf2, 0xfd, 0xd3, 0xd3, 0x2b, 0xf2, 0xfa,
	0xe4, 0xbc, 0xaf, 0xeb, 0x1f, 0x6d, 0x85, 0x9a, 0xc7, 0x4b, 0x2e, 0xbb, 0xf7, 0x9c, 0x9e, 0x6e,
	0xa4, 0x6c, 0xbe, 0x7d, 0xfc, 0xf1, 0xf3, 0x2c, 0xbe, 0x88, 0xc9, 0x2b, 0xa1, 0x7d, 0xa2, 0xce,
	0xe8, 0xed, 0x8e, 0xfe, 0xc7, 0xff, 0x50, 0x3e, 0x3a, 0x85, 0xcf, 0x6d, 0x0d, 0xf3, 0xa8, 0x39,
	0x73, 0x7d, 0x14, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x19, 0xce, 0x5a, 0x12, 0xb0, 0x02, 0x00,
	0x00,
}
