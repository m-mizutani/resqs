package errors

import (
	"golang.org/x/xerrors"
)

// Error is ReSQS owned error structure. It can has Values by With function in order to keep context of error.
type Error struct {
	Err    error
	Values map[string]interface{}
}

// Error() is implemented according to basic Go error interface
func (x *Error) Error() string {
	return x.Err.Error()
}

// With sets a key and a value to Values. They should retrieves when showing error
func (x *Error) With(key string, value interface{}) *Error {
	x.Values[key] = value
	return x
}

// New returns Error structure with msg
func New(msg string) *Error {
	return &Error{
		Err:    xerrors.New(msg),
		Values: make(map[string]interface{}),
	}
}

// Wrap returns Error with cause error.
func Wrap(cause error, msg string) *Error {
	return &Error{
		Err:    xerrors.Errorf("%s: %w", msg, cause),
		Values: make(map[string]interface{}),
	}
}
