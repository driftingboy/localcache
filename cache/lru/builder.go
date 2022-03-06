package lru

import "github.com/driftingboy/localcache/cache"

const Name = "lru"

var _ cache.CacheBuilder = (*CacheBuilder)(nil)

func init() {
	cache.RegisterCache(&CacheBuilder{})
}

type CacheBuilder struct{}

func (cb CacheBuilder) Name() string {
	return Name
}

func (cb CacheBuilder) BuildFunc() func(c *cache.Config) cache.Cache {
	return func(c *cache.Config) cache.Cache {
		if c == nil {
			return NewCache()
		}

		opts := make([]Option, 0)
		if c.MaxBytes != 0 {
			opts = append(opts, WithMaxBytes(c.MaxBytes))
		}
		if c.MaxItems != 0 {
			opts = append(opts, WithMaxItem(c.MaxItems))
		}
		return NewCache(opts...)
	}
}
