// +build !gcp
// +build debug

package log

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func logger() *debugLog {
	return DefaultLogger.(*debugLog)
}

func buffer() *bytes.Buffer {
	return logger().dst.(*bytes.Buffer)
}

func init() {
	logger().dst = &bytes.Buffer{}
}

func TestDefaultDebug(t *testing.T) {
	_, ok := DefaultLogger.(*debugLog)
	if !assert.True(t, ok, "DefaultLogger should be debugLog") {
		t.Logf("%#v", DefaultLogger)
		return
	}
}

func TestEmit(t *testing.T) {
	Info("Hello, World")
	if !assert.Equal(t, "[INFO] Hello, World\n", buffer().String()) {
		return
	}
	buffer().Reset()

	args := map[string]string{"foo": "1", "bar": "2"}
	Info(args)
	var m map[string]string
	b := buffer().Bytes()
	if !assert.NoError(t, json.Unmarshal(b[7:], &m)) {
		return
	}

	if !assert.Equal(t, args, m) {
		return
	}
	buffer().Reset()
}
