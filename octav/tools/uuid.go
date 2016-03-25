package tools

import (
	"crypto/rand"
	"fmt"
)

func UUID() string {
	b := make([]byte, 16)
	rand.Reader.Read(b)
	b[6]=(b[6]&0x0F)|0x40
	b[8]=(b[8]&^0x40)|0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}