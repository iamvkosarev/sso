// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: sso/v1/sso.proto

package sso

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
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
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on RegisterUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RegisterUserRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterUserRequestMultiError, or nil if none found.
func (m *RegisterUserRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterUserRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetEmail()); l < 5 || l > 254 {
		err := RegisterUserRequestValidationError{
			field:  "Email",
			reason: "value length must be between 5 and 254 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = RegisterUserRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 6 || l > 128 {
		err := RegisterUserRequestValidationError{
			field:  "Password",
			reason: "value length must be between 6 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return RegisterUserRequestMultiError(errors)
	}

	return nil
}

func (m *RegisterUserRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *RegisterUserRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// RegisterUserRequestMultiError is an error wrapping multiple validation
// errors returned by RegisterUserRequest.ValidateAll() if the designated
// constraints aren't met.
type RegisterUserRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterUserRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterUserRequestMultiError) AllErrors() []error { return m }

// RegisterUserRequestValidationError is the validation error returned by
// RegisterUserRequest.Validate if the designated constraints aren't met.
type RegisterUserRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterUserRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterUserRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterUserRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterUserRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterUserRequestValidationError) ErrorName() string {
	return "RegisterUserRequestValidationError"
}

// Error satisfies the builtin error interface
func (e RegisterUserRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterUserRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterUserRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterUserRequestValidationError{}

// Validate checks the field values on RegisterUserResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RegisterUserResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterUserResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterUserResponseMultiError, or nil if none found.
func (m *RegisterUserResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterUserResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	if len(errors) > 0 {
		return RegisterUserResponseMultiError(errors)
	}

	return nil
}

// RegisterUserResponseMultiError is an error wrapping multiple validation
// errors returned by RegisterUserResponse.ValidateAll() if the designated
// constraints aren't met.
type RegisterUserResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterUserResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterUserResponseMultiError) AllErrors() []error { return m }

// RegisterUserResponseValidationError is the validation error returned by
// RegisterUserResponse.Validate if the designated constraints aren't met.
type RegisterUserResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterUserResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterUserResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterUserResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterUserResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterUserResponseValidationError) ErrorName() string {
	return "RegisterUserResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RegisterUserResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterUserResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterUserResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterUserResponseValidationError{}

// Validate checks the field values on LoginUserRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *LoginUserRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LoginUserRequestMultiError, or nil if none found.
func (m *LoginUserRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginUserRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetEmail()); l < 5 || l > 254 {
		err := LoginUserRequestValidationError{
			field:  "Email",
			reason: "value length must be between 5 and 254 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = LoginUserRequestValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 6 || l > 128 {
		err := LoginUserRequestValidationError{
			field:  "Password",
			reason: "value length must be between 6 and 128 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return LoginUserRequestMultiError(errors)
	}

	return nil
}

func (m *LoginUserRequest) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *LoginUserRequest) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

// LoginUserRequestMultiError is an error wrapping multiple validation errors
// returned by LoginUserRequest.ValidateAll() if the designated constraints
// aren't met.
type LoginUserRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginUserRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginUserRequestMultiError) AllErrors() []error { return m }

// LoginUserRequestValidationError is the validation error returned by
// LoginUserRequest.Validate if the designated constraints aren't met.
type LoginUserRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginUserRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginUserRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginUserRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginUserRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginUserRequestValidationError) ErrorName() string { return "LoginUserRequestValidationError" }

// Error satisfies the builtin error interface
func (e LoginUserRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginUserRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginUserRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginUserRequestValidationError{}

// Validate checks the field values on LoginUserResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *LoginUserResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginUserResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LoginUserResponseMultiError, or nil if none found.
func (m *LoginUserResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginUserResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	// no validation rules for UserId

	if len(errors) > 0 {
		return LoginUserResponseMultiError(errors)
	}

	return nil
}

// LoginUserResponseMultiError is an error wrapping multiple validation errors
// returned by LoginUserResponse.ValidateAll() if the designated constraints
// aren't met.
type LoginUserResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginUserResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginUserResponseMultiError) AllErrors() []error { return m }

// LoginUserResponseValidationError is the validation error returned by
// LoginUserResponse.Validate if the designated constraints aren't met.
type LoginUserResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginUserResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginUserResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginUserResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginUserResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginUserResponseValidationError) ErrorName() string {
	return "LoginUserResponseValidationError"
}

// Error satisfies the builtin error interface
func (e LoginUserResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginUserResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginUserResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginUserResponseValidationError{}

// Validate checks the field values on VerifyTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *VerifyTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on VerifyTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// VerifyTokenRequestMultiError, or nil if none found.
func (m *VerifyTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *VerifyTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return VerifyTokenRequestMultiError(errors)
	}

	return nil
}

// VerifyTokenRequestMultiError is an error wrapping multiple validation errors
// returned by VerifyTokenRequest.ValidateAll() if the designated constraints
// aren't met.
type VerifyTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m VerifyTokenRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m VerifyTokenRequestMultiError) AllErrors() []error { return m }

// VerifyTokenRequestValidationError is the validation error returned by
// VerifyTokenRequest.Validate if the designated constraints aren't met.
type VerifyTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VerifyTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VerifyTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VerifyTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VerifyTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VerifyTokenRequestValidationError) ErrorName() string {
	return "VerifyTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e VerifyTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVerifyTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VerifyTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VerifyTokenRequestValidationError{}

// Validate checks the field values on VerifyTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *VerifyTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on VerifyTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// VerifyTokenResponseMultiError, or nil if none found.
func (m *VerifyTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *VerifyTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	if len(errors) > 0 {
		return VerifyTokenResponseMultiError(errors)
	}

	return nil
}

// VerifyTokenResponseMultiError is an error wrapping multiple validation
// errors returned by VerifyTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type VerifyTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m VerifyTokenResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m VerifyTokenResponseMultiError) AllErrors() []error { return m }

// VerifyTokenResponseValidationError is the validation error returned by
// VerifyTokenResponse.Validate if the designated constraints aren't met.
type VerifyTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e VerifyTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e VerifyTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e VerifyTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e VerifyTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e VerifyTokenResponseValidationError) ErrorName() string {
	return "VerifyTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e VerifyTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sVerifyTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = VerifyTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = VerifyTokenResponseValidationError{}
