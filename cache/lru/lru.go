package lru

import (
	"container/list"

	"github.com/driftingboy/localcache/cache"
)

// default config
var (
	defaultMaxItem        = 4096
	defaultMaxBytes int64 = 4096 * 1024
)

type Cache struct {
	// limit max items in cache
	maxItem int
	// limit max bytes in cache
	maxBytes int64

	// current data's bytes size
	bytesNum int64
	index    map[string]*list.Element
	data     *list.List

	afterDelKey func(key cache.Key, val cache.Value)
}

func NewCache(opts ...Option) *Cache {
	c := &Cache{
		maxItem:  defaultMaxItem,
		maxBytes: defaultMaxBytes,
		index:    make(map[string]*list.Element, defaultMaxItem),
		data:     list.New(),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

type Item struct {
	key   cache.Key
	value cache.Value
}

func (c *Cache) Set(key cache.Key, value cache.Value) {
	k := key.Key()
	if e, ok := c.index[k]; ok {
		c.data.MoveToFront(e)
		item := e.Value.(*Item)
		c.bytesNum += value.BytesNum() - item.value.BytesNum()
		item.value = value
	} else {
		newE := c.data.PushFront(&Item{key: key, value: value})
		c.index[k] = newE
		c.bytesNum += value.BytesNum()
	}

	// check max limit
	for c.isMaxLimit() {
		c.DelOldest()
	}
}

// is max limit reached
func (c *Cache) isMaxLimit() bool {
	if c.maxItem <= 0 && c.maxBytes <= 0 {
		return false
	} else if c.maxItem <= 0 && c.maxBytes > 0 {
		return c.bytesNum > c.maxBytes
	} else if c.maxItem > 0 && c.maxBytes <= 0 {
		return c.data.Len() > c.maxItem
	} else {
		return c.data.Len() > c.maxItem || c.bytesNum > c.maxBytes
	}
}

func (c Cache) Get(key cache.Key) (v cache.Value, ok bool) {
	e, ok := c.index[key.Key()]
	if !ok {
		return nil, false
	}

	c.data.MoveToFront(e)
	return e.Value.(*Item).value, true
}

func (c *Cache) Del(key cache.Key) {
	if c == nil {
		return
	}

	e, ok := c.index[key.Key()]
	if !ok {
		return
	}

	c.delElement(e)
}

func (c *Cache) Clear() {}

func (c *Cache) DelOldest() {
	if c == nil {
		return
	}

	e := c.data.Back()
	if e == nil {
		return
	}

	c.delElement(e)
}

func (c *Cache) delElement(e *list.Element) {
	kv := e.Value.(*Item)
	c.data.Remove(e)
	delete(c.index, kv.key.Key())
	c.bytesNum -= int64(len(kv.key.Key())) + kv.value.BytesNum()

	if c.afterDelKey != nil {
		c.afterDelKey(kv.key, kv.value)
	}
}
