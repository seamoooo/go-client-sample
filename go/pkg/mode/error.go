package mode

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// ErrInternal represents an internal error such as network outage. It
	// also happens when programming error.
	ErrInternal = "internal"

	// ErrUnknown represents an error is not known error which usually
	// means the error is not properly handled.
	ErrUnknown = "unknown"

	// ErrNotFound represents a resource is not found.
	ErrNotFound = "not-found"

	// ErrInvalidArgument represents a resource is not found.
	ErrInvalidArgument = "invalid-argument"
)

// Error is the standard error for motorist domain.
//
// NOTE: This error type is taken from
// https://middlemost.com/failure-is-your-domain/. Read the article first.
type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable message for end users (i.e., third party developers).
	// It must be gramatically correct American English.
	Message string

	// Op is the operation being performed, usually the name of the method
	// being invoked (Get, Put, etc.).
	Op string

	// The underlying error that triggered this one, if any.
	Err error
}

var _ error = (*Error)(nil)

// Error returns message for internal developers which should contain enough
// info to debug.
func (e *Error) Error() string {
	var b strings.Builder

	if e.Op != "" {
		fmt.Fprintf(&b, "%s: ", e.Op)
	}

	if e.Err != nil {
		b.WriteString(e.Err.Error())
	}

	if e.Code != "" {
		fmt.Fprintf(&b, "<%s> ", e.Code)
	}

	b.WriteString(e.Message)

	return b.String()
}

// ErrorCode returns error code of err when its type is Error. Otherwise it
// returns ErrUnknown.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		if e.Code == "" {
			if e.Err != nil {
				return ErrorCode(e.Err)
			}

			return ErrUnknown
		}

		return e.Code
	}

	return ErrUnknown
}

// ErrorMessage returns human-readable message of err for end users (i.e.,
// third party developers).
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		if e.Message != "" {
			return e.Message
		}

		return ErrorMessage(e.Err)
	}

	return "please contact technical support"
}
