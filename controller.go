package localcache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/driftingboy/localcache/singleflight"
)

// A Loader loads local data for a key.
type Loader interface {
	// error need impl IsNotFound function
	Load(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type LoaderFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f LoaderFunc) Load(key string) ([]byte, error) {
	return f(key)
}

var (
	// protect dbs
	lock sync.RWMutex
	dbs  = make(map[string]*CacheDB)
)

type CacheDB struct {
	name string
	// load data in source. if nil，It will not go back to the source, and retrun directly
	loader     Loader
	peerPicker PeerPicker

	mainCache *SyncCache
	// TODO 3. hotCache  *SyncCache, 频繁访问的数据移到此做冷热隔离，hot cache 可以和maincache 选择不同的淘汰策略比如 不淘汰

	sf *singleflight.Group
}

// TODO 1. loader, peerPicker 都移入 opts
func NewCacheDB(name string, l Loader, mc *SyncCache) *CacheDB {
	if mc == nil {
		panic("mainCache can't set nil!")
	}
	return &CacheDB{
		name:      name,
		loader:    l,
		mainCache: mc,
		// hotCache:  hc,
		sf: &singleflight.Group{},
	}
}

func RegisterDB(cc *CacheDB) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := dbs[cc.name]; ok {
		panic("can't set duplicate cache conttroller!")
	}
	dbs[cc.name] = cc
}

func RegisterDBE(cc *CacheDB) error {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := dbs[cc.name]; ok {
		return errors.New("can't set duplicate cache conttroller!")
	}
	dbs[cc.name] = cc
	return nil
}

func GetDB(name string) (c *CacheDB, b bool) {
	lock.RLock()
	defer lock.RUnlock()

	c, b = dbs[name]
	return
}

func (cc *CacheDB) RegisterPeerPicker(p PeerPicker) {
	if cc.peerPicker != nil {
		panic("just register once!")
	}
	cc.peerPicker = p
}

var ErrNotFound = errors.New("notfound")

func (cc *CacheDB) Get(k string) (b ByteView, err error) {
	if k == "" {
		return ByteView{}, errors.New("key require!")
	}

	if b, ok := cc.mainCache.Get(k); ok {
		return b, nil
	}

	return cc.load(k)
}

func (cc *CacheDB) load(k string) (b ByteView, err error) {
	if cc.peerPicker == nil && cc.loader == nil {
		return ByteView{}, ErrNotFound
	}

	v, err, _ := cc.sf.Do(k, func() (interface{}, error) {
		if cc.peerPicker != nil {
			r, ok, isSelf := cc.peerPicker.Pick(k)
			if ok && !isSelf {
				b, err = cc.loadFormRemote(k, r)
				if err == nil {
					return b, nil
				}
				// TODO 如果err = cache miss，则返回；网络错误重试后说明节点宕机，去其他节点访问
				fmt.Printf("load in remote err: %v", err)
			}
		}

		if cc.loader != nil {
			b, err = cc.loadFormLocall(k)
			if err != nil {
				return ByteView{}, err
			}
			cc.Set(k, b)
		}
		return b, nil
	})

	if err != nil {
		return ByteView{}, err
	}

	return v.(ByteView), nil
}

func (cc *CacheDB) loadFormLocall(k string) (b ByteView, err error) {
	originBytes, err := cc.loader.Load(k)
	if err != nil {
		if isOriginDataNotFound(err) {
			return ByteView{}, ErrNotFound
		}
		return ByteView{}, err
	}

	b = ByteView{b: originBytes}

	return
}

func (cc *CacheDB) loadFormRemote(k string, rl RemoteLoader) (b ByteView, err error) {
	originBytes, err := rl.Load(context.TODO(), cc.name, k)
	if err != nil {
		if isOriginDataNotFound(err) {
			return ByteView{}, ErrNotFound
		}
		return ByteView{}, err
	}

	return ByteView{b: originBytes}, nil
}

func (cc *CacheDB) Set(k string, v ByteView) {
	cc.mainCache.Set(k, v)
}

func (cc *CacheDB) Del(k string) {
	cc.mainCache.Del(k)
}

func isOriginDataNotFound(err error) bool {
	if e, ok := err.(interface{ IsNotFound() bool }); ok {
		return e.IsNotFound()
	}
	return false
}
