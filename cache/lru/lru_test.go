package lru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type value string

func (v value) BytesNum() int64 {
	return int64(len(v))
}

func Example() {
	cache := NewCache()
	cache.Set("a", value("test"))
	got, _ := cache.Get("a")
	fmt.Println(got)
	cache.Del("a")
	_, ok := cache.Get("a")
	fmt.Println(ok)
	//Output:
	//test
	//false
}

func TestCache_Set(t *testing.T) {
	c := NewCache(WithMaxItem(3))
	c.Set("a", value("a"))
	c.Set("b", value("b"))
	c.Set("c", value("c"))
	c.Set("d", value("d"))

	_, ok := c.Get("a")
	assert.False(t, ok)
	_, ok = c.Get("b")
	assert.True(t, ok)
	c.Set("e", value("e"))
	_, ok1 := c.Get("b")
	_, ok2 := c.Get("c")
	assert.True(t, ok1)
	assert.False(t, ok2)
}
