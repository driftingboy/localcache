package lru

import "github.com/driftingboy/localcache/cache"

type Option func(c *Cache)

func WithMaxItem(mi int) Option {
	return func(c *Cache) {
		c.maxItem = mi
	}
}

func WithMaxBytes(mb int64) Option {
	return func(c *Cache) {
		c.maxBytes = mb
	}
}

func WithAfterDelFunc(f func(key string, val cache.Value)) Option {
	return func(c *Cache) {
		c.afterDelKey = f
	}
}
