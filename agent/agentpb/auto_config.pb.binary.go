// Code generated by protoc-gen-go-binary. DO NOT EDIT.
// source: agent/agentpb/auto_config.proto

package agentpb

import (
	"github.com/golang/protobuf/proto"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (msg *AutoConfigRequest) MarshalBinary() ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (msg *AutoConfigRequest) UnmarshalBinary(b []byte) error {
	return proto.Unmarshal(b, msg)
}

// MarshalBinary implements encoding.BinaryMarshaler
func (msg *AutoConfigResponse) MarshalBinary() ([]byte, error) {
	return proto.Marshal(msg)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (msg *AutoConfigResponse) UnmarshalBinary(b []byte) error {
	return proto.Unmarshal(b, msg)
}
