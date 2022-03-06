package localcache

import (
	"fmt"
	"testing"

	"github.com/driftingboy/localcache/cache"
	"github.com/driftingboy/localcache/cache/lru"
)

type LoadErr string

func (l LoadErr) Error() string {
	return string(l)
}

func (l LoadErr) IsNotFound() bool {
	return string(l) == "10400"
}

type MyLoader map[string]string

func (m MyLoader) Load(key string) ([]byte, error) {
	if v, ok := m[key]; ok {
		return []byte(v), nil
	}
	return nil, LoadErr("10400")
}

var loader = MyLoader{
	"1": "a",
	"2": "b",
	"3": "c",
	"4": "d",
}

func Example() {
	sc := NewSyncCache(cache.NewCache(lru.Name, nil))
	cc := NewCacheDB("c-1", loader, sc)
	data, _ := cc.Get("1")
	fmt.Println("key 1, data:", data)
	_, err := cc.Get("5")
	fmt.Println("key 5, db load", err.Error())

	// cc1 no load, if cache miss, return ErrNotFound
	cc1 := NewCacheDB("c-2", nil, NewSyncCache(cache.NewCache(lru.Name, nil)))
	_, err = cc1.Get("1")
	fmt.Println("key 1, cache load", err.Error())

	cc1.Set("1", ByteView{b: []byte("a")})
	data, _ = cc1.Get("1")
	fmt.Println("key 1, data:", data)
	//Output:
	//key 1, data: a
	//key 5, db load notfound
	//key 1, cache load notfound
	//key 1, data: a
}

func TestCacheController_Get(t *testing.T) {

}
