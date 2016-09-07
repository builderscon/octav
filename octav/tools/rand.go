package tools

import (
	"crypto/rand"
	"math"
	"math/big"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var maxInt = big.NewInt(math.MaxInt64)

func RandomString(n int) string {
	b := make([]byte, n)
	// A randsrc.Int63() generates 63 random bits, enough for
	// letterIdxMax characters!
	cache, _ := rand.Int(rand.Reader, maxInt)
	for i, remain := n-1, letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, _ = rand.Int(rand.Reader, maxInt)
			remain = letterIdxMax
		}
		if idx := int(cache.Int64() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache.SetInt64(cache.Int64() >> letterIdxBits)
		remain--
	}

	return string(b)
}

func RandFloat64() float64 {
again:
	fbig, _ := rand.Int(rand.Reader, maxInt)
	f := float64(fbig.Int64()) / (1 << 63)
	if f == 1 {
		goto again
	}
	return f
}
