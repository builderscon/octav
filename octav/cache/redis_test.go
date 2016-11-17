package cache_test

import (
	"testing"

	cache "github.com/builderscon/octav/octav/cache"
	"github.com/stretchr/testify/assert"
	redis "gopkg.in/redis.v5"
)

var redisAddr = "127.0.0.1:6379"

func redisAvailable() bool {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := client.Ping().Result(); err != nil {
		return false
	}
	return true
}

func TestRedis(t *testing.T) {
	var v, x  struct {
		Foo string
		Bar []byte
	}
	v.Foo = "Hello"
	v.Bar = []byte("World!")

	c := cache.NewRedis([]string{redisAddr})

	key := "foo"
	c.Delete(key)
	if !assert.Error(t, c.Get(key, &x), "Get should fail") {
		return
	}

	if !assert.NoError(t, c.Set(key, &v), "Set should succeed") {
		return
	}

	if !assert.NoError(t, c.Get(key, &x), "Get should succeed") {
		return
	}

	if !assert.Equal(t, x, v, "items should be equal") {
		return
	}

	if !assert.NoError(t, c.Delete(key), "Delete should succeed") {
		return
	}

	if !assert.Error(t, c.Get(key, &x), "Get should fail") {
		return
	}
}
