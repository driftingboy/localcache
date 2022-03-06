package localcache

import (
	"sync"

	"github.com/driftingboy/localcache/cache"
)

// warpper lru with lock
// TODO 分片锁
// RW lock cannot be used because there are modification operations to move linked list elements in `get`
type SyncCache struct {
	mutex sync.Mutex

	c cache.Cache
}

func NewSyncCache(c cache.Cache) *SyncCache {
	return &SyncCache{
		c: c,
	}
}

func (sc *SyncCache) Set(k string, v ByteView) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.c.Set(k, v)
}

func (sc *SyncCache) Get(k string) (v ByteView, ok bool) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if v, ok := sc.c.Get(k); ok {
		return v.(ByteView), ok
	}
	return
}

func (sc *SyncCache) Del(k string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	sc.c.Del(k)
}
