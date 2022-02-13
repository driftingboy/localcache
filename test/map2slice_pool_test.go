package test

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

// 如何判断是否要使用 slice 替换 map 呢？
// 首先要知道两者的区别
// - slice 查询某个元素时间复杂度 O(n),但是方便复用
// - map 查询某个元素 O(1)，但是不方便复用内存
// 也就是要权衡 查询时间复杂度和内存分配（gc）的开销
// 在保证内存利用率的情况下保证最优的查询时间消耗
// slice 适合数据量小，访问高频的情况下使用

// 10 100 Benchmark_searchInMap-4   	  260016	      4812 ns/op	     420 B/op	       1 allocs/op
// 10 1000 Benchmark_searchInMap-4   	   29638	     39704 ns/op	     420 B/op	       1 allocs/op

// 100 100 Benchmark_searchInMap-4   	  142267	      8406 ns/op	    4188 B/op	       2 allocs/op
// 1000 100 Benchmark_searchInMap-4   	   17918	     63244 ns/op	   60256 B/op	     902 allocs/op
func Benchmark_searchInMap(b *testing.B) {
	size, useCount := 10, 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchInMap(size, useCount)
	}
}

// 10 100 Benchmark_searchInSlice-4   	  269385	      3890 ns/op	     320 B/op	      10 allocs/op
// 10 1000 Benchmark_searchInSlice-4   	   35610	     34716 ns/op	       320 B/op	      10 allocs/op

// 100 100 Benchmark_searchInSlice-4   	  226148	      5319 ns/op	       0 B/op	       0 allocs/op
// 1000 100 Benchmark_searchInSlice-4   	   45426	     25622 ns/op	       0 B/op	       0 allocs/op

var slicePool = &sync.Pool{
	New: func() interface{} {
		// return &slice{}
		return NewSlice(10)
	},
}

func Benchmark_searchInSlice(b *testing.B) {
	size, useCount := 10, 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchInSlicePool(size, useCount)
	}
}

// map delete 会引起gc，没法复用；
func searchInMap(mapSize int, useCount int) {
	// 1. init map
	m := make(map[int]string)
	for i := 0; i < mapSize; i++ {
		k := i
		v := strconv.Itoa(i)
		m[k] = v
	}

	// search
	var v string
	for i := 0; i < useCount; i++ {
		// use map
		v = m[rand.Intn(mapSize)]
		_ = v
	}
}

type KVPair struct {
	key   int
	value string
}

type slice struct {
	items []*KVPair
}

func NewSlice(size int) *slice {
	return &slice{
		items: make([]*KVPair, 0, size),
	}
}

func (s *slice) Append(kv *KVPair) {
	s.items = append(s.items, kv)
}

func (s slice) Len() int {
	return len(s.items)
}

func (s *slice) Reset() {
	if s.items == nil {
		return
	}
	s.items = s.items[:0]
}

func (s *slice) Find(key int) (value string, ok bool) {
	// 可优化为2分查找
	for _, item := range s.items {
		if item != nil && item.key == key {
			value = item.value
			ok = true
			return
		}
	}
	return
}

// slice = slice[:] 可以清理数据，但是不回收内存
func searchInSlicePool(size, useCount int) {
	s := slicePool.Get().(*slice)
	for i := 0; i < size; i++ {
		s.Append(&KVPair{key: i, value: strconv.Itoa(i)})
	}

	for i := 0; i < useCount; i++ {
		s.Find(rand.Intn(s.Len()))
	}

	s.Reset()
	slicePool.Put(s)
}
