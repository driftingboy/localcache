package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/driftingboy/localcache"
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

func main() {
	var port int
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.Parse()

	cc := localcache.NewCacheDB(
		"db1",
		loader,
		localcache.NewSyncCache(cache.NewCache(
			lru.Name, &cache.Config{
				MaxItems: 2 << 10,        // 1024
				MaxBytes: (2 << 10) * 16, // 16 kb
			},
		)),
	)
	localcache.RegisterDB(cc)

	endpoint := fmt.Sprintf("localhost:%d", port)
	hp := localcache.NewHTTPPoolOpts(endpoint, nil)

	hp.Set("localhost:8001", "localhost:8002", "localhost:8003")
	cc.RegisterPeerPicker(hp)

	hp.Log("v start %v", time.Now())

	log.Fatal(http.ListenAndServe(endpoint, hp))
}
