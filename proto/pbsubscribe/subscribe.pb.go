// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/pbsubscribe/subscribe.proto

package pbsubscribe

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	pbservice "github.com/hashicorp/consul/proto/pbservice"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// Topic enumerates the supported event topics.
type Topic int32

const (
	Topic_Unknown Topic = 0
	// ServiceHealth topic contains events for any changes to service health.
	Topic_ServiceHealth Topic = 1
	// ServiceHealthConnect topic contains events for any changes to service
	// health for connect-enabled services.
	Topic_ServiceHealthConnect Topic = 2
)

var Topic_name = map[int32]string{
	0: "Unknown",
	1: "ServiceHealth",
	2: "ServiceHealthConnect",
}

var Topic_value = map[string]int32{
	"Unknown":              0,
	"ServiceHealth":        1,
	"ServiceHealthConnect": 2,
}

func (x Topic) String() string {
	return proto.EnumName(Topic_name, int32(x))
}

func (Topic) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{0}
}

type CatalogOp int32

const (
	CatalogOp_Register   CatalogOp = 0
	CatalogOp_Deregister CatalogOp = 1
)

var CatalogOp_name = map[int32]string{
	0: "Register",
	1: "Deregister",
}

var CatalogOp_value = map[string]int32{
	"Register":   0,
	"Deregister": 1,
}

func (x CatalogOp) String() string {
	return proto.EnumName(CatalogOp_name, int32(x))
}

func (CatalogOp) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{1}
}

// SubscribeRequest used to subscribe to a topic.
type SubscribeRequest struct {
	// Topic identifies the set of events the subscriber is interested in.
	Topic Topic `protobuf:"varint,1,opt,name=Topic,proto3,enum=subscribe.Topic" json:"Topic,omitempty"`
	// Key is a topic-specific identifier that restricts the scope of the
	// subscription to only events pertaining to that identifier. For example,
	// to receive events for a single service, the service's name is specified
	// as the key.
	Key string `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	// Token is the ACL token to authenticate the request. The token must have
	// sufficient privileges to read the requested information otherwise events
	// will be filtered, possibly resulting in an empty snapshot and no further
	// updates sent.
	Token string `protobuf:"bytes,3,opt,name=Token,proto3" json:"Token,omitempty"`
	// Index is the raft index the subscriber has already observed up to. This
	// is zero on an initial streaming call, but then can be provided by a
	// client on subsequent re-connections such that the full snapshot doesn't
	// need to be resent if the client is up to date.
	Index uint64 `protobuf:"varint,4,opt,name=Index,proto3" json:"Index,omitempty"`
	// Datacenter specifies the Consul datacenter the request is targeted at.
	// If it's not the local DC the server will forward the request to
	// the remote DC and proxy the results back  to the subscriber. An empty
	// string defaults to the local datacenter.
	Datacenter string `protobuf:"bytes,5,opt,name=Datacenter,proto3" json:"Datacenter,omitempty"`
	// Namespace which contains the resources. If Namespace is not specified the
	// default namespace will be used.
	//
	// Namespace is an enterprise-only feature.
	Namespace string `protobuf:"bytes,6,opt,name=Namespace,proto3" json:"Namespace,omitempty"`
	// Partition which contains the resources. If Partition is not specified the
	// default partition will be used.
	//
	// Partition is an enterprise-only feature.
	Partition            string   `protobuf:"bytes,7,opt,name=Partition,proto3" json:"Partition,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubscribeRequest) Reset()         { *m = SubscribeRequest{} }
func (m *SubscribeRequest) String() string { return proto.CompactTextString(m) }
func (*SubscribeRequest) ProtoMessage()    {}
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{0}
}

func (m *SubscribeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubscribeRequest.Unmarshal(m, b)
}
func (m *SubscribeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubscribeRequest.Marshal(b, m, deterministic)
}
func (m *SubscribeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubscribeRequest.Merge(m, src)
}
func (m *SubscribeRequest) XXX_Size() int {
	return xxx_messageInfo_SubscribeRequest.Size(m)
}
func (m *SubscribeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubscribeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubscribeRequest proto.InternalMessageInfo

func (m *SubscribeRequest) GetTopic() Topic {
	if m != nil {
		return m.Topic
	}
	return Topic_Unknown
}

func (m *SubscribeRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SubscribeRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *SubscribeRequest) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *SubscribeRequest) GetDatacenter() string {
	if m != nil {
		return m.Datacenter
	}
	return ""
}

func (m *SubscribeRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *SubscribeRequest) GetPartition() string {
	if m != nil {
		return m.Partition
	}
	return ""
}

// Event describes a streaming update on a subscription. Events are used both to
// describe the current "snapshot" of the result as well as ongoing mutations to
// that snapshot.
type Event struct {
	// Index is the raft index at which the mutation took place. At the top
	// level of a subscription there will always be at most one Event per index.
	// If multiple events are published to the same topic in a single raft
	// transaction then the batch of events will be encoded inside a single
	// top-level event to ensure they are delivered atomically to clients.
	Index uint64 `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	// Payload is the actual event content.
	//
	// Types that are valid to be assigned to Payload:
	//	*Event_EndOfSnapshot
	//	*Event_NewSnapshotToFollow
	//	*Event_EventBatch
	//	*Event_ServiceHealth
	Payload              isEvent_Payload `protobuf_oneof:"Payload"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

type isEvent_Payload interface {
	isEvent_Payload()
}

type Event_EndOfSnapshot struct {
	EndOfSnapshot bool `protobuf:"varint,2,opt,name=EndOfSnapshot,proto3,oneof"`
}

type Event_NewSnapshotToFollow struct {
	NewSnapshotToFollow bool `protobuf:"varint,3,opt,name=NewSnapshotToFollow,proto3,oneof"`
}

type Event_EventBatch struct {
	EventBatch *EventBatch `protobuf:"bytes,4,opt,name=EventBatch,proto3,oneof"`
}

type Event_ServiceHealth struct {
	ServiceHealth *ServiceHealthUpdate `protobuf:"bytes,10,opt,name=ServiceHealth,proto3,oneof"`
}

func (*Event_EndOfSnapshot) isEvent_Payload() {}

func (*Event_NewSnapshotToFollow) isEvent_Payload() {}

func (*Event_EventBatch) isEvent_Payload() {}

func (*Event_ServiceHealth) isEvent_Payload() {}

func (m *Event) GetPayload() isEvent_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Event) GetEndOfSnapshot() bool {
	if x, ok := m.GetPayload().(*Event_EndOfSnapshot); ok {
		return x.EndOfSnapshot
	}
	return false
}

func (m *Event) GetNewSnapshotToFollow() bool {
	if x, ok := m.GetPayload().(*Event_NewSnapshotToFollow); ok {
		return x.NewSnapshotToFollow
	}
	return false
}

func (m *Event) GetEventBatch() *EventBatch {
	if x, ok := m.GetPayload().(*Event_EventBatch); ok {
		return x.EventBatch
	}
	return nil
}

func (m *Event) GetServiceHealth() *ServiceHealthUpdate {
	if x, ok := m.GetPayload().(*Event_ServiceHealth); ok {
		return x.ServiceHealth
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Event) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Event_EndOfSnapshot)(nil),
		(*Event_NewSnapshotToFollow)(nil),
		(*Event_EventBatch)(nil),
		(*Event_ServiceHealth)(nil),
	}
}

type EventBatch struct {
	Events               []*Event `protobuf:"bytes,1,rep,name=Events,proto3" json:"Events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventBatch) Reset()         { *m = EventBatch{} }
func (m *EventBatch) String() string { return proto.CompactTextString(m) }
func (*EventBatch) ProtoMessage()    {}
func (*EventBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{2}
}

func (m *EventBatch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventBatch.Unmarshal(m, b)
}
func (m *EventBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventBatch.Marshal(b, m, deterministic)
}
func (m *EventBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventBatch.Merge(m, src)
}
func (m *EventBatch) XXX_Size() int {
	return xxx_messageInfo_EventBatch.Size(m)
}
func (m *EventBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_EventBatch.DiscardUnknown(m)
}

var xxx_messageInfo_EventBatch proto.InternalMessageInfo

func (m *EventBatch) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type ServiceHealthUpdate struct {
	Op                   CatalogOp                   `protobuf:"varint,1,opt,name=Op,proto3,enum=subscribe.CatalogOp" json:"Op,omitempty"`
	CheckServiceNode     *pbservice.CheckServiceNode `protobuf:"bytes,2,opt,name=CheckServiceNode,proto3" json:"CheckServiceNode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ServiceHealthUpdate) Reset()         { *m = ServiceHealthUpdate{} }
func (m *ServiceHealthUpdate) String() string { return proto.CompactTextString(m) }
func (*ServiceHealthUpdate) ProtoMessage()    {}
func (*ServiceHealthUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab3eb8c810e315fb, []int{3}
}

func (m *ServiceHealthUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceHealthUpdate.Unmarshal(m, b)
}
func (m *ServiceHealthUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceHealthUpdate.Marshal(b, m, deterministic)
}
func (m *ServiceHealthUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceHealthUpdate.Merge(m, src)
}
func (m *ServiceHealthUpdate) XXX_Size() int {
	return xxx_messageInfo_ServiceHealthUpdate.Size(m)
}
func (m *ServiceHealthUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceHealthUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceHealthUpdate proto.InternalMessageInfo

func (m *ServiceHealthUpdate) GetOp() CatalogOp {
	if m != nil {
		return m.Op
	}
	return CatalogOp_Register
}

func (m *ServiceHealthUpdate) GetCheckServiceNode() *pbservice.CheckServiceNode {
	if m != nil {
		return m.CheckServiceNode
	}
	return nil
}

func init() {
	proto.RegisterEnum("subscribe.Topic", Topic_name, Topic_value)
	proto.RegisterEnum("subscribe.CatalogOp", CatalogOp_name, CatalogOp_value)
	proto.RegisterType((*SubscribeRequest)(nil), "subscribe.SubscribeRequest")
	proto.RegisterType((*Event)(nil), "subscribe.Event")
	proto.RegisterType((*EventBatch)(nil), "subscribe.EventBatch")
	proto.RegisterType((*ServiceHealthUpdate)(nil), "subscribe.ServiceHealthUpdate")
}

func init() {
	proto.RegisterFile("proto/pbsubscribe/subscribe.proto", fileDescriptor_ab3eb8c810e315fb)
}

var fileDescriptor_ab3eb8c810e315fb = []byte{
	// 527 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0x5d, 0x6f, 0xda, 0x3c,
	0x14, 0xc6, 0xb4, 0x85, 0xe6, 0xf0, 0xb6, 0xca, 0xeb, 0x32, 0x2d, 0xa2, 0x53, 0xc5, 0xd0, 0x54,
	0xb1, 0x4a, 0x23, 0x13, 0x93, 0xb6, 0xbb, 0x49, 0x83, 0xb6, 0x63, 0x9a, 0x04, 0x55, 0x68, 0x2f,
	0xb6, 0x3b, 0xe3, 0x9c, 0x91, 0x88, 0xd4, 0xf6, 0x12, 0x53, 0xd6, 0xfb, 0xed, 0x1f, 0xee, 0x07,
	0x4d, 0x31, 0x21, 0x04, 0xe8, 0x9d, 0xcf, 0xf3, 0xe1, 0x63, 0x9f, 0x0f, 0x78, 0xa9, 0x62, 0xa9,
	0xa5, 0xab, 0x26, 0xc9, 0x7c, 0x92, 0xf0, 0x38, 0x9c, 0xa0, 0x9b, 0x9f, 0x3a, 0x86, 0xa3, 0x56,
	0x0e, 0x34, 0x1a, 0xb9, 0x1a, 0xe3, 0x87, 0x90, 0xa3, 0x2b, 0xa4, 0x9f, 0xc9, 0x5a, 0x7f, 0x09,
	0xd8, 0xe3, 0x95, 0xd2, 0xc3, 0x9f, 0x73, 0x4c, 0x34, 0x3d, 0x87, 0x83, 0x5b, 0xa9, 0x42, 0xee,
	0x90, 0x26, 0x69, 0x1f, 0x77, 0xed, 0xce, 0xfa, 0x72, 0x83, 0x7b, 0x4b, 0x9a, 0xda, 0xb0, 0xf7,
	0x15, 0x1f, 0x9d, 0x72, 0x93, 0xb4, 0x2d, 0x2f, 0x3d, 0xd2, 0x7a, 0xea, 0x9c, 0xa1, 0x70, 0xf6,
	0x0c, 0xb6, 0x0c, 0x52, 0xf4, 0x8b, 0xf0, 0xf1, 0x97, 0xb3, 0xdf, 0x24, 0xed, 0x7d, 0x6f, 0x19,
	0xd0, 0x33, 0x80, 0x4b, 0xa6, 0x19, 0x47, 0xa1, 0x31, 0x76, 0x0e, 0x8c, 0xa1, 0x80, 0xd0, 0x17,
	0x60, 0x0d, 0xd9, 0x3d, 0x26, 0x8a, 0x71, 0x74, 0x2a, 0x86, 0x5e, 0x03, 0x29, 0x7b, 0xc3, 0x62,
	0x1d, 0xea, 0x50, 0x0a, 0xa7, 0xba, 0x64, 0x73, 0xa0, 0xf5, 0xa7, 0x0c, 0x07, 0x57, 0x0f, 0x28,
	0xf4, 0x3a, 0x37, 0x29, 0xe6, 0x3e, 0x87, 0xa3, 0x2b, 0xe1, 0x8f, 0x7e, 0x8c, 0x05, 0x53, 0x49,
	0x20, 0xb5, 0xf9, 0xc3, 0xe1, 0xa0, 0xe4, 0x6d, 0xc2, 0xb4, 0x0b, 0x27, 0x43, 0x5c, 0xac, 0xc2,
	0x5b, 0x79, 0x2d, 0xa3, 0x48, 0x2e, 0xcc, 0xef, 0x52, 0xf5, 0x53, 0x24, 0xfd, 0x00, 0x60, 0x52,
	0xf7, 0x98, 0xe6, 0x81, 0xf9, 0x72, 0xad, 0xfb, 0xac, 0x50, 0xc2, 0x35, 0x39, 0x28, 0x79, 0x05,
	0x29, 0xbd, 0x86, 0xa3, 0xf1, 0xb2, 0x43, 0x03, 0x64, 0x91, 0x0e, 0x1c, 0x30, 0xde, 0xb3, 0x82,
	0x77, 0x83, 0xbf, 0x53, 0x3e, 0xd3, 0x98, 0x3e, 0x7a, 0x03, 0xee, 0x59, 0x50, 0xbd, 0x61, 0x8f,
	0x91, 0x64, 0x7e, 0xeb, 0x7d, 0xf1, 0x2d, 0xb4, 0x0d, 0x15, 0x13, 0x25, 0x0e, 0x69, 0xee, 0xb5,
	0x6b, 0x1b, 0x8d, 0x35, 0x84, 0x97, 0xf1, 0xad, 0xdf, 0x04, 0x4e, 0x9e, 0xc8, 0x45, 0x5f, 0x41,
	0x79, 0xa4, 0xb2, 0xb1, 0xa8, 0x17, 0xdc, 0x7d, 0xa6, 0x59, 0x24, 0xa7, 0x23, 0xe5, 0x95, 0x47,
	0x8a, 0x7e, 0x06, 0xbb, 0x1f, 0x20, 0x9f, 0x65, 0x37, 0x0c, 0xa5, 0x8f, 0xa6, 0xc0, 0xb5, 0xee,
	0x69, 0x27, 0x9f, 0xc2, 0xce, 0xb6, 0xc4, 0xdb, 0x31, 0x5d, 0x7c, 0xca, 0x06, 0x91, 0xd6, 0xa0,
	0x7a, 0x27, 0x66, 0x42, 0x2e, 0x84, 0x5d, 0xa2, 0xff, 0x6f, 0xd5, 0xc9, 0x26, 0xd4, 0x81, 0xfa,
	0x06, 0xd4, 0x97, 0x42, 0x20, 0xd7, 0x76, 0xf9, 0xe2, 0x35, 0x58, 0xf9, 0xe3, 0xe8, 0x7f, 0x70,
	0xe8, 0xe1, 0x34, 0x4c, 0x34, 0xc6, 0x76, 0x89, 0x1e, 0x03, 0x5c, 0x62, 0xbc, 0x8a, 0x49, 0xf7,
	0x1b, 0x3c, 0x1f, 0x6b, 0xa6, 0xb1, 0x1f, 0x30, 0x31, 0xc5, 0x6c, 0x2b, 0x54, 0x3a, 0x4f, 0xf4,
	0x23, 0x58, 0xf9, 0x96, 0xd0, 0xd3, 0x62, 0x43, 0xb6, 0x76, 0xa7, 0xb1, 0x53, 0xd3, 0x56, 0xe9,
	0x2d, 0xe9, 0xb9, 0xdf, 0xdf, 0x4c, 0x43, 0x1d, 0xcc, 0x27, 0x1d, 0x2e, 0xef, 0xdd, 0x80, 0x25,
	0x41, 0xc8, 0x65, 0xac, 0x5c, 0x2e, 0x45, 0x32, 0x8f, 0xdc, 0x9d, 0x75, 0x9e, 0x54, 0x0c, 0xf4,
	0xee, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa0, 0xb3, 0x69, 0x51, 0xea, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StateChangeSubscriptionClient is the client API for StateChangeSubscription service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StateChangeSubscriptionClient interface {
	// Subscribe to a topic to receive events when there are changes to the topic.
	//
	// If SubscribeRequest.Index is 0 the event stream will start with one or
	// more snapshot events, followed by an EndOfSnapshot event. Subsequent
	// events will be a live stream of events as they happen.
	//
	// If SubscribeRequest.Index is > 0 it is assumed the client already has a
	// snapshot, and is trying to resume a stream that was disconnected. The
	// client will either receive a NewSnapshotToFollow event, indicating the
	// client view is stale and it must reset its view and prepare for a new
	// snapshot. Or, if no NewSnapshotToFollow event is received, the client
	// view is still fresh, and all events will be the live stream.
	//
	// Subscribe may return a gRPC status error with codes.ABORTED to indicate
	// the client view is now stale due to a change on the server. The client
	// must reset its view and issue a new Subscribe call to restart the stream.
	// This error is used when the server can no longer correctly maintain the
	// stream, for example because the ACL permissions for the token changed, or
	// because the server state was restored from a snapshot.
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (StateChangeSubscription_SubscribeClient, error)
}

type stateChangeSubscriptionClient struct {
	cc grpc.ClientConnInterface
}

func NewStateChangeSubscriptionClient(cc grpc.ClientConnInterface) StateChangeSubscriptionClient {
	return &stateChangeSubscriptionClient{cc}
}

func (c *stateChangeSubscriptionClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (StateChangeSubscription_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StateChangeSubscription_serviceDesc.Streams[0], "/subscribe.StateChangeSubscription/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &stateChangeSubscriptionSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StateChangeSubscription_SubscribeClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type stateChangeSubscriptionSubscribeClient struct {
	grpc.ClientStream
}

func (x *stateChangeSubscriptionSubscribeClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StateChangeSubscriptionServer is the server API for StateChangeSubscription service.
type StateChangeSubscriptionServer interface {
	// Subscribe to a topic to receive events when there are changes to the topic.
	//
	// If SubscribeRequest.Index is 0 the event stream will start with one or
	// more snapshot events, followed by an EndOfSnapshot event. Subsequent
	// events will be a live stream of events as they happen.
	//
	// If SubscribeRequest.Index is > 0 it is assumed the client already has a
	// snapshot, and is trying to resume a stream that was disconnected. The
	// client will either receive a NewSnapshotToFollow event, indicating the
	// client view is stale and it must reset its view and prepare for a new
	// snapshot. Or, if no NewSnapshotToFollow event is received, the client
	// view is still fresh, and all events will be the live stream.
	//
	// Subscribe may return a gRPC status error with codes.ABORTED to indicate
	// the client view is now stale due to a change on the server. The client
	// must reset its view and issue a new Subscribe call to restart the stream.
	// This error is used when the server can no longer correctly maintain the
	// stream, for example because the ACL permissions for the token changed, or
	// because the server state was restored from a snapshot.
	Subscribe(*SubscribeRequest, StateChangeSubscription_SubscribeServer) error
}

// UnimplementedStateChangeSubscriptionServer can be embedded to have forward compatible implementations.
type UnimplementedStateChangeSubscriptionServer struct {
}

func (*UnimplementedStateChangeSubscriptionServer) Subscribe(req *SubscribeRequest, srv StateChangeSubscription_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}

func RegisterStateChangeSubscriptionServer(s *grpc.Server, srv StateChangeSubscriptionServer) {
	s.RegisterService(&_StateChangeSubscription_serviceDesc, srv)
}

func _StateChangeSubscription_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StateChangeSubscriptionServer).Subscribe(m, &stateChangeSubscriptionSubscribeServer{stream})
}

type StateChangeSubscription_SubscribeServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type stateChangeSubscriptionSubscribeServer struct {
	grpc.ServerStream
}

func (x *stateChangeSubscriptionSubscribeServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

var _StateChangeSubscription_serviceDesc = grpc.ServiceDesc{
	ServiceName: "subscribe.StateChangeSubscription",
	HandlerType: (*StateChangeSubscriptionServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _StateChangeSubscription_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/pbsubscribe/subscribe.proto",
}
