// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/type/matcher/metadata.proto

package envoy_type_matcher

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
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

type MetadataMatcher struct {
	Filter               string                         `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	Path                 []*MetadataMatcher_PathSegment `protobuf:"bytes,2,rep,name=path,proto3" json:"path,omitempty"`
	Value                *ValueMatcher                  `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *MetadataMatcher) Reset()         { *m = MetadataMatcher{} }
func (m *MetadataMatcher) String() string { return proto.CompactTextString(m) }
func (*MetadataMatcher) ProtoMessage()    {}
func (*MetadataMatcher) Descriptor() ([]byte, []int) {
	return fileDescriptor_865eaf6a1e9e266d, []int{0}
}

func (m *MetadataMatcher) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataMatcher.Unmarshal(m, b)
}
func (m *MetadataMatcher) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataMatcher.Marshal(b, m, deterministic)
}
func (m *MetadataMatcher) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataMatcher.Merge(m, src)
}
func (m *MetadataMatcher) XXX_Size() int {
	return xxx_messageInfo_MetadataMatcher.Size(m)
}
func (m *MetadataMatcher) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataMatcher.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataMatcher proto.InternalMessageInfo

func (m *MetadataMatcher) GetFilter() string {
	if m != nil {
		return m.Filter
	}
	return ""
}

func (m *MetadataMatcher) GetPath() []*MetadataMatcher_PathSegment {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *MetadataMatcher) GetValue() *ValueMatcher {
	if m != nil {
		return m.Value
	}
	return nil
}

type MetadataMatcher_PathSegment struct {
	// Types that are valid to be assigned to Segment:
	//	*MetadataMatcher_PathSegment_Key
	Segment              isMetadataMatcher_PathSegment_Segment `protobuf_oneof:"segment"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *MetadataMatcher_PathSegment) Reset()         { *m = MetadataMatcher_PathSegment{} }
func (m *MetadataMatcher_PathSegment) String() string { return proto.CompactTextString(m) }
func (*MetadataMatcher_PathSegment) ProtoMessage()    {}
func (*MetadataMatcher_PathSegment) Descriptor() ([]byte, []int) {
	return fileDescriptor_865eaf6a1e9e266d, []int{0, 0}
}

func (m *MetadataMatcher_PathSegment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataMatcher_PathSegment.Unmarshal(m, b)
}
func (m *MetadataMatcher_PathSegment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataMatcher_PathSegment.Marshal(b, m, deterministic)
}
func (m *MetadataMatcher_PathSegment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataMatcher_PathSegment.Merge(m, src)
}
func (m *MetadataMatcher_PathSegment) XXX_Size() int {
	return xxx_messageInfo_MetadataMatcher_PathSegment.Size(m)
}
func (m *MetadataMatcher_PathSegment) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataMatcher_PathSegment.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataMatcher_PathSegment proto.InternalMessageInfo

type isMetadataMatcher_PathSegment_Segment interface {
	isMetadataMatcher_PathSegment_Segment()
}

type MetadataMatcher_PathSegment_Key struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3,oneof"`
}

func (*MetadataMatcher_PathSegment_Key) isMetadataMatcher_PathSegment_Segment() {}

func (m *MetadataMatcher_PathSegment) GetSegment() isMetadataMatcher_PathSegment_Segment {
	if m != nil {
		return m.Segment
	}
	return nil
}

func (m *MetadataMatcher_PathSegment) GetKey() string {
	if x, ok := m.GetSegment().(*MetadataMatcher_PathSegment_Key); ok {
		return x.Key
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*MetadataMatcher_PathSegment) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*MetadataMatcher_PathSegment_Key)(nil),
	}
}

func init() {
	proto.RegisterType((*MetadataMatcher)(nil), "envoy.type.matcher.MetadataMatcher")
	proto.RegisterType((*MetadataMatcher_PathSegment)(nil), "envoy.type.matcher.MetadataMatcher.PathSegment")
}

func init() { proto.RegisterFile("envoy/type/matcher/metadata.proto", fileDescriptor_865eaf6a1e9e266d) }

var fileDescriptor_865eaf6a1e9e266d = []byte{
	// 304 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0xb9, 0xa4, 0xff, 0x70, 0x05, 0x54, 0x5e, 0xa8, 0x8a, 0x00, 0xc3, 0xd4, 0xc9, 0x96,
	0xca, 0x06, 0x0b, 0xf2, 0xc4, 0x52, 0xa9, 0x0a, 0x12, 0xfb, 0x41, 0x0d, 0x8d, 0x68, 0xe3, 0x28,
	0xb9, 0x44, 0x64, 0x63, 0x66, 0xe4, 0xe3, 0xf0, 0x09, 0x58, 0xf9, 0x36, 0xa8, 0x0b, 0x28, 0x89,
	0x2b, 0xa1, 0x92, 0xcd, 0xf2, 0xbd, 0xf7, 0x7b, 0xef, 0x8e, 0x9d, 0x99, 0x28, 0xb7, 0x85, 0xa2,
	0x22, 0x36, 0x6a, 0x85, 0xf4, 0xb0, 0x30, 0x89, 0x5a, 0x19, 0xc2, 0x39, 0x12, 0xca, 0x38, 0xb1,
	0x64, 0x39, 0xaf, 0x24, 0xb2, 0x94, 0x48, 0x27, 0x19, 0x9d, 0x34, 0xd8, 0x72, 0x5c, 0x66, 0xa6,
	0xf6, 0x8c, 0x8e, 0xb3, 0x79, 0x8c, 0x0a, 0xa3, 0xc8, 0x12, 0x52, 0x68, 0xa3, 0x54, 0xa5, 0x84,
	0x94, 0xa5, 0x6e, 0x7c, 0x98, 0xe3, 0x32, 0x9c, 0x23, 0x19, 0xb5, 0x79, 0xd4, 0x83, 0xf3, 0x1f,
	0x60, 0x07, 0x53, 0x17, 0x3f, 0xad, 0xb9, 0xfc, 0x94, 0x75, 0x1e, 0xc3, 0x25, 0x99, 0x64, 0x08,
	0x02, 0xc6, 0xbb, 0xba, 0xbb, 0xd6, 0xad, 0xc4, 0x13, 0x10, 0xb8, 0x6f, 0x3e, 0x65, 0xad, 0x18,
	0x69, 0x31, 0xf4, 0x84, 0x3f, 0xee, 0x4f, 0x94, 0xfc, 0xdf, 0x57, 0x6e, 0x31, 0xe5, 0x0c, 0x69,
	0x71, 0x6b, 0x9e, 0x56, 0x26, 0x22, 0xdd, 0x5b, 0xeb, 0xf6, 0x3b, 0x78, 0x3d, 0x08, 0x2a, 0x0c,
	0xbf, 0x66, 0xed, 0x6a, 0x95, 0xa1, 0x2f, 0x60, 0xdc, 0x9f, 0x88, 0x26, 0xde, 0x5d, 0x29, 0x70,
	0xb0, 0x0a, 0xf0, 0x06, 0xde, 0x00, 0x82, 0xda, 0x38, 0xba, 0x64, 0xfd, 0x3f, 0x01, 0xfc, 0x88,
	0xf9, 0xcf, 0xa6, 0xd8, 0x6a, 0x7f, 0xb3, 0x13, 0x94, 0xbf, 0x7a, 0x9f, 0x75, 0x53, 0xa7, 0xf3,
	0xbf, 0x35, 0xe8, 0xab, 0x8f, 0xd7, 0xcf, 0xaf, 0x8e, 0x37, 0x00, 0x26, 0x42, 0x5b, 0x47, 0xc7,
	0x89, 0x7d, 0x29, 0x1a, 0x5a, 0xe8, 0xbd, 0xcd, 0x5a, 0xb3, 0xf2, 0x78, 0x33, 0xb8, 0xef, 0x54,
	0x57, 0xbc, 0xf8, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x76, 0x5a, 0x88, 0xd6, 0x01, 0x00, 0x00,
}
