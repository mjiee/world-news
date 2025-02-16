package errorx

import (
	"fmt"
)

// BasicError is a basic error type
type BasicError struct {
	code    uint32
	message string
	err     error
}

// NewBasicError creates a new instance of BasicError with the provided error code and error message.
func NewBasicError(code uint32, message string) *BasicError {
	return &BasicError{
		code:    code,
		message: message,
	}
}

// Error implements the error interface.
func (e *BasicError) Error() string {
	if e.IsInternalError() {
		return fmt.Sprintf("%+v", e.err)
	}

	return e.message
}

// SetErr set the internal error.
func (e *BasicError) SetErr(err error) *BasicError {
	b := *e

	b.err = err

	return &b
}

// SetMessage set custom messages
func (e *BasicError) SetMessage(message string) *BasicError {
	b := *e

	b.message = message

	return &b
}

// IsInternalError checks whether the error is an internal error.
func (e *BasicError) IsInternalError() bool {
	return e.err != nil
}

// GetCode returns the error code.
func (e *BasicError) GetCode() uint32 {
	return e.code
}

// GetMessage return the error message.
func (e *BasicError) GetMessage() string {
	return e.message
}
