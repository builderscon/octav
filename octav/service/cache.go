package service

import (
	"os"
	"strconv"
	"sync"

	"github.com/builderscon/octav/octav/cache"
	pdebug "github.com/lestrrat/go-pdebug"
)

var cacheSvc *cache.Redis
var cacheOnce sync.Once
var DefaultCachePrefix = "octav"
var DefaultCacheMagic = ""

func Cache() *cache.Redis {
	cacheOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "127.0.0.1:6379"
		}

		if pdebug.Enabled {
			pdebug.Printf("Using redis server %s for cache backend", addr)
		}

		var options []cache.RedisOption
		prefix := os.Getenv("CACHE_PREFIX")
		if prefix == "" {
			prefix = DefaultCachePrefix
		}
		if prefix != "" {
			if pdebug.Enabled {
				pdebug.Printf("Using cache prefix %s", strconv.Quote(prefix))
			}
			options = append(options, cache.WithPrefix(prefix))
		}

		magic := os.Getenv("CACHE_MAGIC")
		if magic == "" {
			magic = DefaultCacheMagic
		}
		if magic != "" {
			if pdebug.Enabled {
				pdebug.Printf("Using cache magic %s", strconv.Quote(magic))
			}
			options = append(options, cache.WithMagic(magic))
		}

		cacheSvc = cache.NewRedis([]string{addr}, options...)
	})
	return cacheSvc
}
