// +build !gcp
// +build !debug

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDebug(t *testing.T) {
	_, ok := DefaultLogger.(nullLog)
	if !assert.True(t, ok, "DefaultLogger should be nullLog") {
		t.Logf("%#v", DefaultLogger)
		return
	}
}
