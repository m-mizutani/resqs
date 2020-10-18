package errors

import (
	"fmt"
)

type Error struct {
	msg    string
	parent error

	Values map[string]interface{}
}

func (x *Error) Error() string {
	if x.parent != nil {
		return fmt.Sprintf("%s: %v", x.msg, x.parent)
	}
	return x.msg
}

func (x *Error) With(key string, value interface{}) *Error {
	x.Values[key] = value
	return x
}

func New(msg string) *Error {
	return &Error{
		msg:    msg,
		Values: make(map[string]interface{}),
	}
}

func Wrap(parent error, msg string) *Error {
	err := New(msg)
	err.parent = parent
	return err
}
