// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/pbautoconf/auto_config.proto

package pbautoconf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	pbconfig "github.com/hashicorp/consul/proto/pbconfig"
	pbconnect "github.com/hashicorp/consul/proto/pbconnect"
	io "io"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// AutoConfigRequest is the data structure to be sent along with the
// AutoConfig.InitialConfiguration RPC
type AutoConfigRequest struct {
	// Datacenter is the local datacenter name. This wont actually be set by clients
	// but rather will be set by the servers to allow for forwarding to
	// the leader. If it ever happens to be set and differs from the local datacenters
	// name then an error should be returned.
	Datacenter string `protobuf:"bytes,1,opt,name=Datacenter,proto3" json:"Datacenter,omitempty"`
	// Node is the node name that the requester would like to assume
	// the identity of.
	Node string `protobuf:"bytes,2,opt,name=Node,proto3" json:"Node,omitempty"`
	// Segment is the network segment that the requester would like to join
	Segment string `protobuf:"bytes,4,opt,name=Segment,proto3" json:"Segment,omitempty"`
	// JWT is a signed JSON Web Token used to authorize the request
	JWT string `protobuf:"bytes,5,opt,name=JWT,proto3" json:"JWT,omitempty"`
	// ConsulToken is a Consul ACL token that the agent requesting the
	// configuration already has.
	ConsulToken string `protobuf:"bytes,6,opt,name=ConsulToken,proto3" json:"ConsulToken,omitempty"`
	// CSR is a certificate signing request to be used when generating the
	// agents TLS certificate
	CSR                  string   `protobuf:"bytes,7,opt,name=CSR,proto3" json:"CSR,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AutoConfigRequest) Reset()         { *m = AutoConfigRequest{} }
func (m *AutoConfigRequest) String() string { return proto.CompactTextString(m) }
func (*AutoConfigRequest) ProtoMessage()    {}
func (*AutoConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccc5af992e5daf69, []int{0}
}
func (m *AutoConfigRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AutoConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AutoConfigRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AutoConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AutoConfigRequest.Merge(m, src)
}
func (m *AutoConfigRequest) XXX_Size() int {
	return m.Size()
}
func (m *AutoConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AutoConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AutoConfigRequest proto.InternalMessageInfo

func (m *AutoConfigRequest) GetDatacenter() string {
	if m != nil {
		return m.Datacenter
	}
	return ""
}

func (m *AutoConfigRequest) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *AutoConfigRequest) GetSegment() string {
	if m != nil {
		return m.Segment
	}
	return ""
}

func (m *AutoConfigRequest) GetJWT() string {
	if m != nil {
		return m.JWT
	}
	return ""
}

func (m *AutoConfigRequest) GetConsulToken() string {
	if m != nil {
		return m.ConsulToken
	}
	return ""
}

func (m *AutoConfigRequest) GetCSR() string {
	if m != nil {
		return m.CSR
	}
	return ""
}

// AutoConfigResponse is the data structure sent in response to a AutoConfig.InitialConfiguration request
type AutoConfigResponse struct {
	// Config is the partial Consul configuration to inject into the agents own configuration
	Config *pbconfig.Config `protobuf:"bytes,1,opt,name=Config,proto3" json:"Config,omitempty"`
	// CARoots is the current list of Connect CA Roots
	CARoots *pbconnect.CARoots `protobuf:"bytes,2,opt,name=CARoots,proto3" json:"CARoots,omitempty"`
	// Certificate is the TLS certificate issued for the agent
	Certificate *pbconnect.IssuedCert `protobuf:"bytes,3,opt,name=Certificate,proto3" json:"Certificate,omitempty"`
	// ExtraCACertificates holds non-Connect certificates that may be necessary
	// to verify TLS connections with the Consul servers
	ExtraCACertificates  []string `protobuf:"bytes,4,rep,name=ExtraCACertificates,proto3" json:"ExtraCACertificates,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AutoConfigResponse) Reset()         { *m = AutoConfigResponse{} }
func (m *AutoConfigResponse) String() string { return proto.CompactTextString(m) }
func (*AutoConfigResponse) ProtoMessage()    {}
func (*AutoConfigResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccc5af992e5daf69, []int{1}
}
func (m *AutoConfigResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AutoConfigResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AutoConfigResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AutoConfigResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AutoConfigResponse.Merge(m, src)
}
func (m *AutoConfigResponse) XXX_Size() int {
	return m.Size()
}
func (m *AutoConfigResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AutoConfigResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AutoConfigResponse proto.InternalMessageInfo

func (m *AutoConfigResponse) GetConfig() *pbconfig.Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *AutoConfigResponse) GetCARoots() *pbconnect.CARoots {
	if m != nil {
		return m.CARoots
	}
	return nil
}

func (m *AutoConfigResponse) GetCertificate() *pbconnect.IssuedCert {
	if m != nil {
		return m.Certificate
	}
	return nil
}

func (m *AutoConfigResponse) GetExtraCACertificates() []string {
	if m != nil {
		return m.ExtraCACertificates
	}
	return nil
}

func init() {
	proto.RegisterType((*AutoConfigRequest)(nil), "autoconf.AutoConfigRequest")
	proto.RegisterType((*AutoConfigResponse)(nil), "autoconf.AutoConfigResponse")
}

func init() {
	proto.RegisterFile("proto/pbautoconf/auto_config.proto", fileDescriptor_ccc5af992e5daf69)
}

var fileDescriptor_ccc5af992e5daf69 = []byte{
	// 347 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xc1, 0x4a, 0xc3, 0x30,
	0x18, 0xc7, 0x8d, 0x9b, 0x9b, 0x4b, 0x41, 0x66, 0x76, 0x09, 0x13, 0xcb, 0xd8, 0x41, 0x86, 0x48,
	0x2b, 0x13, 0xaf, 0xc2, 0x8c, 0x1e, 0xf4, 0xe0, 0x21, 0x1b, 0x08, 0x5e, 0xa4, 0xeb, 0xb2, 0xad,
	0xb8, 0x25, 0xb5, 0xf9, 0x02, 0x3e, 0x8a, 0x6f, 0xe0, 0xab, 0x78, 0xd4, 0x37, 0x90, 0xf9, 0x22,
	0xd2, 0xb4, 0x95, 0x20, 0x9e, 0xfa, 0xef, 0xff, 0xf7, 0xfb, 0xe0, 0x6b, 0x13, 0xdc, 0x4f, 0x33,
	0x05, 0x2a, 0x4c, 0xa7, 0x91, 0x01, 0x15, 0x2b, 0x39, 0x0f, 0xf3, 0xf0, 0x98, 0xa7, 0x64, 0x11,
	0x58, 0x48, 0x76, 0x2b, 0xd6, 0x3d, 0xa8, 0xec, 0x82, 0x87, 0xae, 0xd6, 0x3d, 0x74, 0xa0, 0x14,
	0x31, 0x84, 0xe5, 0xb3, 0xc0, 0xfd, 0x37, 0x84, 0xf7, 0x47, 0x06, 0x14, 0xb3, 0x33, 0x5c, 0x3c,
	0x1b, 0xa1, 0x81, 0xf8, 0x18, 0x5f, 0x45, 0x10, 0xc5, 0x42, 0x82, 0xc8, 0x28, 0xea, 0xa1, 0x41,
	0x8b, 0x3b, 0x0d, 0x21, 0xb8, 0x7e, 0xa7, 0x66, 0x82, 0x6e, 0x5b, 0x62, 0x33, 0xa1, 0xb8, 0x39,
	0x16, 0x8b, 0xb5, 0x90, 0x40, 0xeb, 0xb6, 0xae, 0x5e, 0x49, 0x1b, 0xd7, 0x6e, 0xef, 0x27, 0x74,
	0xc7, 0xb6, 0x79, 0x24, 0x3d, 0xec, 0x31, 0x25, 0xb5, 0x59, 0x4d, 0xd4, 0x93, 0x90, 0xb4, 0x61,
	0x89, 0x5b, 0xe5, 0x33, 0x6c, 0xcc, 0x69, 0xb3, 0x98, 0x61, 0x63, 0xde, 0xff, 0x44, 0x98, 0xb8,
	0x9b, 0xea, 0x54, 0x49, 0x2d, 0xc8, 0x11, 0x6e, 0x14, 0x8d, 0x5d, 0xd3, 0x1b, 0xee, 0x05, 0xe5,
	0xe7, 0x97, 0x5e, 0x49, 0xc9, 0x31, 0x6e, 0xb2, 0x11, 0x57, 0x0a, 0xb4, 0xdd, 0xda, 0x1b, 0xb6,
	0x83, 0xea, 0x4f, 0x94, 0x3d, 0xaf, 0x04, 0x72, 0x8e, 0x3d, 0x26, 0x32, 0x48, 0xe6, 0x49, 0x1c,
	0x81, 0xa0, 0x35, 0xeb, 0x77, 0x7e, 0xfd, 0x1b, 0xad, 0x8d, 0x98, 0xe5, 0x06, 0x77, 0x3d, 0x72,
	0x8a, 0x3b, 0xd7, 0x2f, 0x90, 0x45, 0x6c, 0xe4, 0xb4, 0x9a, 0xd6, 0x7b, 0xb5, 0x41, 0x8b, 0xff,
	0x87, 0x2e, 0x2f, 0xde, 0x37, 0x3e, 0xfa, 0xd8, 0xf8, 0xe8, 0x6b, 0xe3, 0xa3, 0xd7, 0x6f, 0x7f,
	0xeb, 0xe1, 0x64, 0x91, 0xc0, 0xd2, 0x4c, 0x83, 0x58, 0xad, 0xc3, 0x65, 0xa4, 0x97, 0x49, 0xac,
	0xb2, 0x34, 0x3f, 0x33, 0x6d, 0x56, 0xe1, 0xdf, 0x5b, 0x31, 0x6d, 0xd8, 0xe6, 0xec, 0x27, 0x00,
	0x00, 0xff, 0xff, 0xe2, 0x1d, 0x6e, 0x48, 0x30, 0x02, 0x00, 0x00,
}

func (m *AutoConfigRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AutoConfigRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Datacenter) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.Datacenter)))
		i += copy(dAtA[i:], m.Datacenter)
	}
	if len(m.Node) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.Node)))
		i += copy(dAtA[i:], m.Node)
	}
	if len(m.Segment) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.Segment)))
		i += copy(dAtA[i:], m.Segment)
	}
	if len(m.JWT) > 0 {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.JWT)))
		i += copy(dAtA[i:], m.JWT)
	}
	if len(m.ConsulToken) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.ConsulToken)))
		i += copy(dAtA[i:], m.ConsulToken)
	}
	if len(m.CSR) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(len(m.CSR)))
		i += copy(dAtA[i:], m.CSR)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *AutoConfigResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AutoConfigResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Config != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(m.Config.Size()))
		n1, err := m.Config.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.CARoots != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(m.CARoots.Size()))
		n2, err := m.CARoots.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.Certificate != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintAutoConfig(dAtA, i, uint64(m.Certificate.Size()))
		n3, err := m.Certificate.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if len(m.ExtraCACertificates) > 0 {
		for _, s := range m.ExtraCACertificates {
			dAtA[i] = 0x22
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintAutoConfig(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AutoConfigRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Datacenter)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	l = len(m.Node)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	l = len(m.Segment)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	l = len(m.JWT)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	l = len(m.ConsulToken)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	l = len(m.CSR)
	if l > 0 {
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *AutoConfigResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Config != nil {
		l = m.Config.Size()
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	if m.CARoots != nil {
		l = m.CARoots.Size()
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	if m.Certificate != nil {
		l = m.Certificate.Size()
		n += 1 + l + sovAutoConfig(uint64(l))
	}
	if len(m.ExtraCACertificates) > 0 {
		for _, s := range m.ExtraCACertificates {
			l = len(s)
			n += 1 + l + sovAutoConfig(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAutoConfig(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAutoConfig(x uint64) (n int) {
	return sovAutoConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AutoConfigRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAutoConfig
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AutoConfigRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AutoConfigRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Datacenter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Datacenter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Segment", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Segment = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field JWT", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.JWT = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsulToken", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConsulToken = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CSR", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CSR = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAutoConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AutoConfigResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAutoConfig
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AutoConfigResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AutoConfigResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = &pbconfig.Config{}
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CARoots", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CARoots == nil {
				m.CARoots = &pbconnect.CARoots{}
			}
			if err := m.CARoots.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Certificate", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Certificate == nil {
				m.Certificate = &pbconnect.IssuedCert{}
			}
			if err := m.Certificate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraCACertificates", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAutoConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtraCACertificates = append(m.ExtraCACertificates, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAutoConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAutoConfig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAutoConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAutoConfig
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAutoConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAutoConfig
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthAutoConfig
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAutoConfig
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipAutoConfig(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthAutoConfig
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthAutoConfig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAutoConfig   = fmt.Errorf("proto: integer overflow")
)
