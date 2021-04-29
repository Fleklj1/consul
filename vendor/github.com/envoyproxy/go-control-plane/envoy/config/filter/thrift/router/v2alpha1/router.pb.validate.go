// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/filter/thrift/router/v2alpha1/router.proto

package envoy_config_filter_thrift_router_v2alpha1

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
var _router_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Router with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Router) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// RouterValidationError is the validation error returned by Router.Validate if
// the designated constraints aren't met.
type RouterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RouterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RouterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RouterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RouterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RouterValidationError) ErrorName() string { return "RouterValidationError" }

// Error satisfies the builtin error interface
func (e RouterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRouter.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RouterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RouterValidationError{}
