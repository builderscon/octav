package tools

import (
	"bytes"

	bufferpool "github.com/lestrrat/go-bufferpool"
)

var bufPool = bufferpool.New()

func GetBuffer() *bytes.Buffer {
	return bufPool.Get()
}

func ReleaseBuffer(buf *bytes.Buffer) {
	bufPool.Release(buf)
}
