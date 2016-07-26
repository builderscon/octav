package errors

import (
	"github.com/lestrrat/go-pdebug"
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

	for err != nil {
		// If the error implements an ignorable error, return the value
		if ie, ok := err.(ignorableError); ok {
			return ie.Ignorable()
		}

		// chase the root cause
		if ce, ok := err.(causer); ok {
			err = ce.Cause()
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
	for err != nil {
		pdebug.Printf("%#v", err)
		if fre, ok := err.(finalizationRequiredError); ok {
			if cb := fre.FinalizeFunc(); cb != nil {
				return cb, true
			}
			return nil, false
		}

		if ce, ok := err.(causer); ok {
			err = ce.Cause()
			continue
		}

		break
	}
	return nil, false
}
