package cache

import (
	"fmt"
	"strings"
)

type Cache interface {
	Set(k string, v Value)
	Get(k string) (v Value, ok bool)
	Del(k string)
	Clear()
}

type Value interface {
	BytesNum() int64
}

type CacheBuilder interface {
	Name() string
	BuildFunc() func(c *Config) Cache
}

type Config struct {
	MaxItems int
	MaxBytes int64
}

var cacheMap = make(map[string]func(c *Config) Cache)

func RegisterCache(cb CacheBuilder) {
	if cb == nil {
		panic("CacheBuilder can not be nil!")
	}
	if cb.Name() == "" || cb.BuildFunc() == nil {
		panic("CacheBuilder.Name and BuildFunc can not be empty!")
	}

	cacheMap[strings.ToLower(cb.Name())] = cb.BuildFunc()
}

// if Config is nil, There will be different default values according to different implementations
// For example, the default value of the `non obsolete` may be larger than `lru`
func NewCache(name string, c *Config) Cache {
	if f, ok := cacheMap[strings.ToLower(name)]; ok {
		return f(c)
	}
	panic(fmt.Errorf("typ: %v, need register init func!", name))
}
