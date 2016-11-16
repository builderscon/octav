package service

import (
	"os"
	"sync"

	"github.com/builderscon/octav/octav/cache"
	pdebug "github.com/lestrrat/go-pdebug"
)

var cacheSvc *cache.Redis
var cacheOnce sync.Once

func Cache() *cache.Redis {
	cacheOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "127.0.0.1:6379"
		}

		if pdebug.Enabled {
			pdebug.Printf("Using redis server %s for cache backend", addr)
		}

		cacheSvc = cache.NewRedis(addr)
	})
	return cacheSvc
}
