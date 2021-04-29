// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/grpc_credential/v2alpha/file_based_metadata.proto

package envoy_config_grpc_credential_v2alpha

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
var _file_based_metadata_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on FileBasedMetadataConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *FileBasedMetadataConfig) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetSecretData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FileBasedMetadataConfigValidationError{
				field:  "SecretData",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for HeaderKey

	// no validation rules for HeaderPrefix

	return nil
}

// FileBasedMetadataConfigValidationError is the validation error returned by
// FileBasedMetadataConfig.Validate if the designated constraints aren't met.
type FileBasedMetadataConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FileBasedMetadataConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FileBasedMetadataConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FileBasedMetadataConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FileBasedMetadataConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FileBasedMetadataConfigValidationError) ErrorName() string {
	return "FileBasedMetadataConfigValidationError"
}

// Error satisfies the builtin error interface
func (e FileBasedMetadataConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFileBasedMetadataConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FileBasedMetadataConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FileBasedMetadataConfigValidationError{}
