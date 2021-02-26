// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/type/v3/hash_policy.proto

package envoy_type_v3

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
var _hash_policy_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on HashPolicy with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *HashPolicy) Validate() error {
	if m == nil {
		return nil
	}

	switch m.PolicySpecifier.(type) {

	case *HashPolicy_SourceIp_:

		if v, ok := interface{}(m.GetSourceIp()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HashPolicyValidationError{
					field:  "SourceIp",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return HashPolicyValidationError{
			field:  "PolicySpecifier",
			reason: "value is required",
		}

	}

	return nil
}

// HashPolicyValidationError is the validation error returned by
// HashPolicy.Validate if the designated constraints aren't met.
type HashPolicyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HashPolicyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HashPolicyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HashPolicyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HashPolicyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HashPolicyValidationError) ErrorName() string { return "HashPolicyValidationError" }

// Error satisfies the builtin error interface
func (e HashPolicyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHashPolicy.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HashPolicyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HashPolicyValidationError{}

// Validate checks the field values on HashPolicy_SourceIp with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HashPolicy_SourceIp) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// HashPolicy_SourceIpValidationError is the validation error returned by
// HashPolicy_SourceIp.Validate if the designated constraints aren't met.
type HashPolicy_SourceIpValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HashPolicy_SourceIpValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HashPolicy_SourceIpValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HashPolicy_SourceIpValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HashPolicy_SourceIpValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HashPolicy_SourceIpValidationError) ErrorName() string {
	return "HashPolicy_SourceIpValidationError"
}

// Error satisfies the builtin error interface
func (e HashPolicy_SourceIpValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHashPolicy_SourceIp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HashPolicy_SourceIpValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HashPolicy_SourceIpValidationError{}
