package lruk

import (
	"container/list"

	"github.com/driftingboy/localcache/cache"
	"github.com/driftingboy/localcache/cache/lru"
)

// lru-2(2q) = lru + newdata linked list(FIFO)
// Avoid a large number of thermal data failures caused by cold data loading, and do a simple isolation

type Cache struct {
	// old items in data, use lru strategy
	data *lru.Cache

	// max buffer size
	maxBufSize int
	// new items in bufferQueue, Avoid the failure of a large number of hot data caused by loading cold data
	bufQueue *list.List
	bufIndex map[string]*list.Element
	// How much time can a stay in the bufferqueue at most
	// maxStayTime
}

func NewCache(maxDataBytes int64, maxBufSize int) *Cache {
	return &Cache{
		data:       lru.NewCache(lru.WithMaxBytes(maxDataBytes)),
		maxBufSize: maxBufSize,
		bufQueue:   list.New(),
		bufIndex:   make(map[string]*list.Element),
	}
}

type Item struct {
	key   string
	value cache.Value
}

func (c *Cache) Set(key string, value cache.Value) {
	if ele, ok := c.bufIndex[key]; !ok {
		newE := c.bufQueue.PushFront(&Item{key: key, value: value})
		c.bufIndex[key] = newE

		if c.bufQueue.Len() > c.maxBufSize {
			c.DelBufOldest()
		}
	} else {
		c.data.Set(key, value)
		c.bufQueue.Remove(ele)
	}

}

func (c *Cache) Get(key string) (v cache.Value, ok bool) {
	if v, ok := c.data.Get(key); ok {
		return v, ok
	}

	if ele, ok := c.bufIndex[key]; ok {
		v := ele.Value.(*Item).value
		c.data.Set(key, v)
		c.bufQueue.Remove(ele)
		return v, ok
	}

	return nil, false
}

func (c *Cache) Del(key string) {
	if ele, ok := c.bufIndex[key]; ok {
		c.delBufElement(ele)
	}

	c.data.Del(key)
}

func (c *Cache) Clear() {}

func (c *Cache) DelBufOldest() {
	if c == nil {
		return
	}

	e := c.bufQueue.Back()
	if e == nil {
		return
	}

	c.delBufElement(e)
}

func (c *Cache) delBufElement(e *list.Element) {
	kv := e.Value.(*Item)
	c.bufQueue.Remove(e)
	delete(c.bufIndex, kv.key)
}

// today finsh it and cache interface
// tomorrow s algo 2 节， 1 leetcode
// x w cache 并发， http server 搭建 ... consisithash
