package lruk

import "fmt"

type key string

func (k key) Key() string {
	return string(k)
}

type value string

func (v value) BytesNum() int64 {
	return int64(len(v))
}

func Example() {
	cache := NewCache(1024, 10)
	cache.Set(key("a"), value("a"))
	v, ok := cache.Get(key("a"))
	fmt.Println(v, ok)
	//Output:
	//a true
}
