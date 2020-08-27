// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/filter/http/rbac/v2/rbac.proto

package envoy_config_filter_http_rbac_v2

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
var _rbac_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on RBAC with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *RBAC) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetRules()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RBACValidationError{
				field:  "Rules",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetShadowRules()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RBACValidationError{
				field:  "ShadowRules",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// RBACValidationError is the validation error returned by RBAC.Validate if the
// designated constraints aren't met.
type RBACValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RBACValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RBACValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RBACValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RBACValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RBACValidationError) ErrorName() string { return "RBACValidationError" }

// Error satisfies the builtin error interface
func (e RBACValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRBAC.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RBACValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RBACValidationError{}

// Validate checks the field values on RBACPerRoute with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *RBACPerRoute) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetRbac()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RBACPerRouteValidationError{
				field:  "Rbac",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// RBACPerRouteValidationError is the validation error returned by
// RBACPerRoute.Validate if the designated constraints aren't met.
type RBACPerRouteValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RBACPerRouteValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RBACPerRouteValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RBACPerRouteValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RBACPerRouteValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RBACPerRouteValidationError) ErrorName() string { return "RBACPerRouteValidationError" }

// Error satisfies the builtin error interface
func (e RBACPerRouteValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRBACPerRoute.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RBACPerRouteValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RBACPerRouteValidationError{}
