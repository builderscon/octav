package errors

import (
	"database/sql"

	daverr "github.com/pkg/errors"
)

type httpCodeError struct {
	error
	httpCode int
}

func WithHTTPCode(err error, code ...int) error {
	if len(code) == 0 {
		code = []int{500}
	}
	return httpCodeError{
		error:    err,
		httpCode: code[0],
	}
}

func (e httpCodeError) HTTPCode() int {
	return e.httpCode
}

func Cause(err error) error {
	return daverr.Cause(err)
}

func Errorf(format string, args ...interface{}) error {
	return daverr.Errorf(format, args...)
}

func New(text string) error {
	return daverr.New(text)
}

func Wrap(cause error, message string) error {
	return daverr.Wrap(cause, message)
}

func Wrapf(cause error, format string, args ...interface{}) error {
	return daverr.Wrapf(cause, format, args...)
}

type causer interface {
	Cause() error
}

type ignorableError interface {
	Ignorable() bool
}

func IsIgnorable(err error) bool {
	if err == nil { // If the error is a nil error, then we can just ignore this
		return true
	}

	// Copy to e to make sure we don't tamper with err
	for e := err; e != nil; {
		// If the error implements an ignorable error, return the value
		if ie, ok := e.(ignorableError); ok {
			return ie.Ignorable()
		}

		// chase the root cause
		if ce, ok := e.(causer); ok {
			e = ce.Cause()
			continue
		}

		break
	}
	return false
}

type finalizationRequiredError interface {
	FinalizeFunc() func() error
}

func IsFinalizationRequired(err error) (func() error, bool) {
	if err == nil {
		return nil, false
	}

	for e := err; e != nil; {
		if fre, ok := e.(finalizationRequiredError); ok {
			if cb := fre.FinalizeFunc(); cb != nil {
				return cb, true
			}
			return nil, false
		}

		if ce, ok := e.(causer); ok {
			e = ce.Cause()
			continue
		}

		break
	}
	return nil, false
}

func findCause(err error, tester func(error) (error, bool)) (error, bool) {
	if err == nil {
		return nil, false
	}

	for e := err; e != nil; {
		if ie, ok := tester(e); ok {
			return ie, ok
		}

		if ce, ok := e.(causer); ok {
			e = ce.Cause()
			continue
		}

		break
	}
	return nil, false
}

func IsSQLNoRows(err error) bool {
	_, ok := findCause(err, func(v error) (error, bool) {
		if v == sql.ErrNoRows {
			return v, true
		}
		return nil, false
	})
	return ok
}
