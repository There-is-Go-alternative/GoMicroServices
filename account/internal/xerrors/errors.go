package xerrors

import (
	"fmt"
	"github.com/pkg/errors"
)

// Error is a simple type base on type string that help to define constant
// https://www.youtube.com/watch?v=pN_lm6QqHcw&list=PLMW8Xq7bXrG5B_gvikeSf3Du3NGBs4yVi&index=24
type Error string

// Error implement error interface
func (e Error) Error() string {
	return string(e)
}

// Is implement Go 1.13 new error interface
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok || e == nil || t == nil {
		return false
	}
	return string(*e) == string(*t)
}

// Unwrap implement Go 1.13 new error interface
func (e *Error) Unwrap() error {
	return nil
}

const (
	InternalError = iota
	ResourceNotFound
)

type ErrorWithCode struct {
	Code int
	Err  error
}

// Error implement error interface
func (e ErrorWithCode) Error() string {
	return fmt.Sprintf("[%d]: %v", e.Code, e.Err)
}

// Is implement Go 1.13 new error interface
func (e *ErrorWithCode) Is(target error) bool {
	t, ok := target.(*ErrorWithCode)
	if !ok || e == nil || t == nil {
		return false
	}
	return e.Code == t.Code && errors.Is(e.Err, t.Err)
}

// Unwrap implement Go 1.13 new error interface
func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}
