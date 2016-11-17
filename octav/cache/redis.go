package cache

import (
	"sort"
	"strconv"
	"time"

	"github.com/builderscon/octav/octav/tools"

	cache "gopkg.in/go-redis/cache.v5"
	redis "gopkg.in/redis.v5"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type Redis struct {
	server *redis.Ring
	codec  *cache.Codec
	prefix string
	magic  string
}

type RedisOption interface {
	Configure(*Redis)
}

type RedisOptionFunc func(*Redis)

func (f RedisOptionFunc) Configure(r *Redis) {
	f(r)
}

func WithPrefix(s string) RedisOption {
	return RedisOptionFunc(func(r *Redis) {
		r.prefix = s
	})
}

func WithMagic(s string) RedisOption {
	return RedisOptionFunc(func(r *Redis) {
		r.magic = s
	})
}

func NewRedis(servers []string, options ...RedisOption) *Redis {
	sort.Strings(servers)

	addrs := make(map[string]string)
	for i := 1; i <= len(servers); i++ {
		addrs["server"+strconv.Itoa(i)] = servers[i-1]
	}

	r := redis.NewRing(&redis.RingOptions{
		Addrs: addrs,
	})
	c := &Redis{
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
	for _, o := range options {
		o.Configure(c)
	}
	return c
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
	it := cache.Item{
		Key:    key,
		Object: v,
		Func:   fn,
	}
	applyOptions(&it, options...)

	return c.codec.Do(&it)
}

func (c *Redis) Key(list ...string) string {
	// The magic number allows us to purge entire cache sets without
	// having to delete each one
	if c.magic != "" {
		list = append(list, c.magic)
	}

	if c.prefix != "" {
		list = append([]string{c.prefix}, list...)
	}

	buf := tools.GetBuffer()
	defer tools.ReleaseBuffer(buf)
	for i := 0; i < len(list); i++ {
		buf.WriteString(list[i])
		if i < len(list)-1 {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}
