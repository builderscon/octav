package cache

import (
	"sort"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack"
	cache "gopkg.in/go-redis/cache.v5"
	redis "gopkg.in/redis.v5"
)

type Redis struct {
	server *redis.Ring
	codec  *cache.Codec
}

func NewRedis(servers ...string) *Redis {
	sort.Strings(servers)

	addrs := make(map[string]string)
	for i := 1; i <= len(servers); i++ {
		addrs["server"+strconv.Itoa(i)] = servers[i-1]
	}

	r := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
	})
	return &Redis{
		server: r,
		codec: &cache.Codec{
			Redis: r,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(key []byte, v interface{}) error {
				return msgpack.Unmarshal(key, v)
			},
		},
	}
}

func (c *Redis) Get(key string, v interface{}) error {
	return c.codec.Get(key, v)
}

func applyOptions(it *cache.Item, options ...Option) {
	for _, o := range options {
		switch o.Name() {
		case "expires":
			it.Expiration = o.Value().(time.Duration)
		}
	}
}

func (c *Redis) Set(key string, v interface{}, options ...Option) error {
	it := cache.Item{
		Key:    key,
		Object: v,
	}
	applyOptions(&it, options...)
	return c.codec.Set(&it)
}

func (c *Redis) Delete(key string) error {
	return c.codec.Delete(key)
}

func (c *Redis) GetOrSet(key string, v interface{}, fn func() (interface{}, error), options ...Option) (interface{}, error) {
	it := cache.Item {
		Key: key,
		Object: v,
		Func: fn,
	}
	applyOptions(&it, options...)

	return c.codec.Do(&it)
}
