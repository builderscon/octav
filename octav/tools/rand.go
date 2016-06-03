package tools

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randsrc = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	// A randsrc.Int63() generates 63 random bits, enough for
	// letterIdxMax characters!
	for i, cache, remain := n-1, randsrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randsrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}