// Package xerror custom error module
package xerror

import (
	"errors"
	"fmt"
)

// XError is a custom error type that carries an HTTP status code, an underlying error, and a user-facing message.
type XError struct {
	Code int
	Err  error
	Msg  string
}

// New creates a new XError with the given HTTP status code and underlying error.
// It panics if the error is nil.
func New(code int, err error) error {
	if err == nil {
		panic("xerror: err is nil")
	}
	return &XError{
		Code: code,
		Err:  err,
	}
}

// Error returns the string representation of the underlying error.
func (e *XError) Error() string {
	if e.Err == nil {
		return "unknown error"
	}
	return e.Err.Error()
}

// Wrap annotates the given error with a message while preserving its XError code.
// Returns nil if the input error is nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	var xe *XError
	if errors.As(err, &xe) {
		return &XError{
			Code: xe.Code,
			Err:  fmt.Errorf("%s: %w", msg, xe.Err),
		}
	}

	return &XError{
		Code: UnknownErr,
		Err:  fmt.Errorf("%s: %w", msg, err),
	}
}

// IfErr returns nil if err is nil, otherwise returns an XError constructed with the given code
func IfErr(code int, err error) error {
	if err == nil {
		return nil
	}

	return &XError{
		Code: code,
		Err:  err,
	}
}

// Code extracts the error code from an error. Returns OK if the error is nil, or UnknownErr if it is not an XError.
func Code(err error) int {
	if err == nil {
		return OK
	}

	var xe *XError
	if errors.As(err, &xe) {
		return xe.Code
	}

	return UnknownErr
}

// Msg extracts the user-facing message from an error. Returns "ok" if the error is nil,
// "internal error" if it is not an XError, or the XError Msg field otherwise.
func Msg(err error) string {
	if err == nil {
		return "ok"
	}

	var xe *XError
	if !errors.As(err, &xe) {
		return "internal error"
	}

	return xe.Msg
}

// IsCode reports whether the given error has the specified error code.
func IsCode(code int, err error) bool {
	if err == nil {
		return false
	}

	var xe *XError
	if errors.As(err, &xe) {
		return code == xe.Code
	}

	return false
}

// WithMsg attaches a user-facing error message for the frontend
func WithMsg(msg string, err error) error {
	if err == nil {
		return nil
	}

	var xe *XError
	if errors.As(err, &xe) {
		return &XError{
			Code: xe.Code,
			Err:  xe.Err,
			Msg:  msg,
		}
	}

	return &XError{
		Code: UnknownErr,
		Err:  err,
		Msg:  msg,
	}
}
