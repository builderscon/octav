package service

import (
	"os"
	"sync"

	"github.com/builderscon/octav/octav/cache"
)

var cacheSvc *cache.Redis
var cacheOnce sync.Once

func Cache() *cache.Redis {
	cacheOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr != "" {
			addr = "redis:6379"
		}
		cacheSvc = cache.NewRedis(addr)
	})
	return cacheSvc
}
