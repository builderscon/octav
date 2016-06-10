package errors

import daverr "github.com/pkg/errors"

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
