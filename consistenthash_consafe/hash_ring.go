package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

type HashF func(key []byte) uint32

type HashRing struct {
	mutex sync.Mutex

	hashF    HashF
	replicas int

	// protect nodeHashes and hashNodeMap
	// TODO 1. 使用 指针看看代替 atomic，看是否data race 2. google 搜索 why
	value atomic.Value
}

type Value struct {
	// sorted slice contain node hash
	nodeHashes []int
	// map node hash and node string
	hashNodeMap map[int]string
}

func NewHashRing(replicas int, f HashF) *HashRing {
	h := &HashRing{
		hashF:    f,
		replicas: replicas,
	}
	h.value.Store(&Value{
		nodeHashes:  make([]int, 0),
		hashNodeMap: make(map[int]string),
	})
	if f == nil {
		h.hashF = crc32.ChecksumIEEE
	}

	return h
}

func (hr *HashRing) AddNode(keys ...string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	// copy on write
	v := hr.copyValue()

	for _, k := range keys {
		for i := 0; i < hr.replicas; i++ {
			rName := strconv.Itoa(i) + k
			rHash := int(hr.hashF([]byte(rName)))

			v.nodeHashes = append(v.nodeHashes, rHash)
			v.hashNodeMap[rHash] = k
		}
	}
	sort.Ints(v.nodeHashes)

	hr.value.Store(v)
}

func (hr *HashRing) GetNode(key string) (name string) {
	v := hr.loadValue()

	if len(v.nodeHashes) == 0 {
		return ""
	}

	h := int(hr.hashF([]byte(key)))

	idx := sort.SearchInts(v.nodeHashes, h)
	// 这里的 i 是最接近 key hash 值的 node hash 下标
	// 但是如果 key hash > max node hash， 则返回的是 len(nodeHashes)，所以需要映射
	if idx == len(v.nodeHashes) {
		idx = 0
	}

	return v.hashNodeMap[v.nodeHashes[idx]]
}

func (hr *HashRing) Remove(key string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	for i := 0; i < hr.replicas; i++ {
		rName := strconv.Itoa(i) + key
		rHash := int(hr.hashF([]byte(rName)))

		hr.removeByHash(rHash)
	}
}

func (hr *HashRing) removeByHash(hash int) {
	// copy on write
	v := hr.copyValue()

	idx := sort.SearchInts(v.nodeHashes, hash)
	if idx == len(v.nodeHashes) { // not found
		return
	} else if idx == len(v.nodeHashes)-1 { // last elem
		v.nodeHashes = v.nodeHashes[:idx]
	} else {
		v.nodeHashes = append(v.nodeHashes[:idx], v.nodeHashes[idx+1:]...)
	}
	delete(v.hashNodeMap, hash)

	hr.storeValue(v)
}

func (hr *HashRing) loadValue() *Value {
	return hr.value.Load().(*Value)
}

func (hr *HashRing) storeValue(v *Value) {
	hr.value.Store(v)
}

// deep copy
func (hr *HashRing) copyValue() *Value {
	v := hr.loadValue()

	newV := &Value{
		nodeHashes:  make([]int, len(v.nodeHashes)),
		hashNodeMap: make(map[int]string, len(v.hashNodeMap)),
	}
	for k, v := range v.hashNodeMap {
		newV.hashNodeMap[k] = v
	}
	copy(newV.nodeHashes, v.nodeHashes)

	return newV
}
