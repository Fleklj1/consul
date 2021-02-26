// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/type/metadata/v3/metadata.proto

package envoy_type_metadata_v3

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _metadata_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on MetadataKey with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *MetadataKey) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetKey()) < 1 {
		return MetadataKeyValidationError{
			field:  "Key",
			reason: "value length must be at least 1 bytes",
		}
	}

	if len(m.GetPath()) < 1 {
		return MetadataKeyValidationError{
			field:  "Path",
			reason: "value must contain at least 1 item(s)",
		}
	}

	for idx, item := range m.GetPath() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetadataKeyValidationError{
					field:  fmt.Sprintf("Path[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// MetadataKeyValidationError is the validation error returned by
// MetadataKey.Validate if the designated constraints aren't met.
type MetadataKeyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKeyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKeyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKeyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKeyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKeyValidationError) ErrorName() string { return "MetadataKeyValidationError" }

// Error satisfies the builtin error interface
func (e MetadataKeyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKey.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKeyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKeyValidationError{}

// Validate checks the field values on MetadataKind with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *MetadataKind) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Kind.(type) {

	case *MetadataKind_Request_:

		if v, ok := interface{}(m.GetRequest()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetadataKindValidationError{
					field:  "Request",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *MetadataKind_Route_:

		if v, ok := interface{}(m.GetRoute()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetadataKindValidationError{
					field:  "Route",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *MetadataKind_Cluster_:

		if v, ok := interface{}(m.GetCluster()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetadataKindValidationError{
					field:  "Cluster",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *MetadataKind_Host_:

		if v, ok := interface{}(m.GetHost()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetadataKindValidationError{
					field:  "Host",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return MetadataKindValidationError{
			field:  "Kind",
			reason: "value is required",
		}

	}

	return nil
}

// MetadataKindValidationError is the validation error returned by
// MetadataKind.Validate if the designated constraints aren't met.
type MetadataKindValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKindValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKindValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKindValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKindValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKindValidationError) ErrorName() string { return "MetadataKindValidationError" }

// Error satisfies the builtin error interface
func (e MetadataKindValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKind.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKindValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKindValidationError{}

// Validate checks the field values on MetadataKey_PathSegment with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MetadataKey_PathSegment) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Segment.(type) {

	case *MetadataKey_PathSegment_Key:

		if len(m.GetKey()) < 1 {
			return MetadataKey_PathSegmentValidationError{
				field:  "Key",
				reason: "value length must be at least 1 bytes",
			}
		}

	default:
		return MetadataKey_PathSegmentValidationError{
			field:  "Segment",
			reason: "value is required",
		}

	}

	return nil
}

// MetadataKey_PathSegmentValidationError is the validation error returned by
// MetadataKey_PathSegment.Validate if the designated constraints aren't met.
type MetadataKey_PathSegmentValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKey_PathSegmentValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKey_PathSegmentValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKey_PathSegmentValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKey_PathSegmentValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKey_PathSegmentValidationError) ErrorName() string {
	return "MetadataKey_PathSegmentValidationError"
}

// Error satisfies the builtin error interface
func (e MetadataKey_PathSegmentValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKey_PathSegment.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKey_PathSegmentValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKey_PathSegmentValidationError{}

// Validate checks the field values on MetadataKind_Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MetadataKind_Request) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// MetadataKind_RequestValidationError is the validation error returned by
// MetadataKind_Request.Validate if the designated constraints aren't met.
type MetadataKind_RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKind_RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKind_RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKind_RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKind_RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKind_RequestValidationError) ErrorName() string {
	return "MetadataKind_RequestValidationError"
}

// Error satisfies the builtin error interface
func (e MetadataKind_RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKind_Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKind_RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKind_RequestValidationError{}

// Validate checks the field values on MetadataKind_Route with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MetadataKind_Route) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// MetadataKind_RouteValidationError is the validation error returned by
// MetadataKind_Route.Validate if the designated constraints aren't met.
type MetadataKind_RouteValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKind_RouteValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKind_RouteValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKind_RouteValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKind_RouteValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKind_RouteValidationError) ErrorName() string {
	return "MetadataKind_RouteValidationError"
}

// Error satisfies the builtin error interface
func (e MetadataKind_RouteValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKind_Route.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKind_RouteValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKind_RouteValidationError{}

// Validate checks the field values on MetadataKind_Cluster with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *MetadataKind_Cluster) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// MetadataKind_ClusterValidationError is the validation error returned by
// MetadataKind_Cluster.Validate if the designated constraints aren't met.
type MetadataKind_ClusterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKind_ClusterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKind_ClusterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKind_ClusterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKind_ClusterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKind_ClusterValidationError) ErrorName() string {
	return "MetadataKind_ClusterValidationError"
}

// Error satisfies the builtin error interface
func (e MetadataKind_ClusterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKind_Cluster.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKind_ClusterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKind_ClusterValidationError{}

// Validate checks the field values on MetadataKind_Host with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *MetadataKind_Host) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// MetadataKind_HostValidationError is the validation error returned by
// MetadataKind_Host.Validate if the designated constraints aren't met.
type MetadataKind_HostValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetadataKind_HostValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetadataKind_HostValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetadataKind_HostValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetadataKind_HostValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetadataKind_HostValidationError) ErrorName() string {
	return "MetadataKind_HostValidationError"
}

// Error satisfies the builtin error interface
func (e MetadataKind_HostValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadataKind_Host.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetadataKind_HostValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetadataKind_HostValidationError{}
