package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type key string

func (k key) Key() string {
	return string(k)
}

type value string

func (v value) BytesNum() int64 {
	return int64(len(v))
}

func Example() {
	cache := NewCache()
	cache.Set(key("a"), value("test"))
	got, _ := cache.Get(key("a"))
	fmt.Println(got)
	cache.Del(key("a"))
	_, ok := cache.Get(key("a"))
	fmt.Println(ok)
	//Output:
	//test
	//false
}

func TestCache_Set(t *testing.T) {
	c := NewCache(WithMaxItem(3))
	c.Set(key("a"), value("a"))
	c.Set(key("b"), value("b"))
	c.Set(key("c"), value("c"))
	c.Set(key("d"), value("d"))

	_, ok := c.Get(key("a"))
	assert.False(t, ok)
	_, ok = c.Get(key("b"))
	assert.True(t, ok)
	c.Set(key("e"), value("e"))
	_, ok1 := c.Get(key("b"))
	_, ok2 := c.Get(key("c"))
	assert.True(t, ok1)
	assert.False(t, ok2)
}
