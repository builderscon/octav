package gettext_test

import (
	"testing"

	"github.com/builderscon/octav/octav/gettext"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	g := gettext.New("_test")
	g.AddDomain("messages")
	if !assert.Equal(t, "こんにちは、世界！", g.Get("ja", "messages", "Hello, World")) {
		return
	}
}
