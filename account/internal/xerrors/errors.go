package xerrors

import (
	"fmt"
	"strings"
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
	return e.Code == t.Code && e.Err.Error() == t.Err.Error()
}

// Unwrap implement Go 1.13 new error interface
func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}

// errList is a type that represent a list of errors to be back-traced.
type errList []error

// Error implement error interface
// As errList is not exported and Concat is the only function that give access to it,
// no errors in errList should be nil
func (l errList) Error() string {
	filtered := make([]string, 0, len(l))
	for _, err := range l {
		filtered = append(filtered, err.Error())
	}
	return strings.Join(filtered, "\n")
}

// Concat takes a list of errors and return a single error out of it or nil if the list was empty.
// It ignores nil errors.
func Concat(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	filtered := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			filtered = append(filtered, err)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return errList(filtered)
}
