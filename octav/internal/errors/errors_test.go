package errors_test

import (
	"fmt"
	"testing"

	"github.com/builderscon/octav/octav/internal/errors"
	"github.com/stretchr/testify/assert"
)

type SimpleIgnorable struct {
	Ignore bool
}

func (e SimpleIgnorable) Error() string {
	return fmt.Sprintf("simple ignorable %t", e.Ignore)
}
func (e SimpleIgnorable) Ignorable() bool {
	return e.Ignore
}

type FinalizationRequired struct {
	Ignore   bool
	Callback func() error
}

func (e FinalizationRequired) Error() string {
	return fmt.Sprintf("finalize required (ignore=%t)", e.Ignore)
}

func (e FinalizationRequired) Ignorable() bool {
	return e.Ignore
}

func (e FinalizationRequired) FinalizeFunc() func() error {
	return e.Callback
}

func TestIgnorable(t *testing.T) {
	table := []struct {
		Error error
		Fn    func(assert.TestingT, bool, ...interface{}) bool
	}{
		{SimpleIgnorable{Ignore: true}, assert.True},
		{SimpleIgnorable{Ignore: false}, assert.False},
		{FinalizationRequired{Ignore: true}, assert.True},
		{FinalizationRequired{Ignore: false}, assert.False},
		{errors.New("regular error"), assert.False},
	}

	for _, data := range table {
		e := data.Error
		fn := data.Fn
		if !fn(t, errors.IsIgnorable(e), "should pass") {
			return
		}
	}
}

func dummyFinalize() error {
	return nil
}

func TestFinalizationRequired(t *testing.T) {
	table := []struct {
		Error error
		Fn    func(assert.TestingT, bool, ...interface{}) bool
	}{
		{SimpleIgnorable{Ignore: true}, assert.False},
		{SimpleIgnorable{Ignore: false}, assert.False},
		{FinalizationRequired{Ignore: true, Callback: dummyFinalize}, assert.True},
		{FinalizationRequired{Ignore: false, Callback: dummyFinalize}, assert.True},
		{errors.New("regular error"), assert.False},
	}

	for _, data := range table {
		e := data.Error
		fn := data.Fn
		cb, ok := errors.IsFinalizationRequired(e)
		if !fn(t, ok, "should pass") {
			return
		}
		if ok {
			if !assert.NotNil(t, cb, "callback should be non-nil") {
				return
			}
		}
	}
}
