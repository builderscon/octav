package tools

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: allocBuffer,
}

func allocBuffer() interface{} {
	return &bytes.Buffer{}
}

func GetBuffer() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

func ReleaseBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buf.Grow(0)
	bufPool.Put(buf)
}
