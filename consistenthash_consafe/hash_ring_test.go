package consistenthash

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashRing_GetNode(t *testing.T) {
	r := NewHashRing(3, func(key []byte) uint32 {
		i, err := strconv.Atoi(string(key))
		if err != nil {
			panic(err)
		}

		return uint32(i)
	})
	// 1, 2, 3, 11, 12, 13, 21, 22, 23
	r.AddNode("1", "2", "3")

	assert.Equal(t, r.GetNode("1"), r.GetNode("11"))
	assert.Equal(t, r.GetNode("10"), r.GetNode("40"))
}

func TestHashRing_DataRace(t *testing.T) {
	r := NewHashRing(3, func(key []byte) uint32 {
		i, err := strconv.Atoi(string(key))
		if err != nil {
			panic(err)
		}

		return uint32(i)
	})
	// 1, 2, 3, 11, 12, 13, 21, 22, 23
	r.AddNode("1", "2", "3")

	var wg sync.WaitGroup

	wg.Add(3)
	// test data race
	go func() {
		r.AddNode("4")
		wg.Done()
	}()
	go func() {
		r.Remove("1")
		wg.Done()
	}()
	go func() {
		r.GetNode("1")
		wg.Done()
	}()
	wg.Wait()

	// 测试节点删除后，是否数据流向下一个节点
	assert.Equal(t, r.GetNode("1"), r.GetNode("11"))
	assert.Equal(t, r.GetNode("1"), r.GetNode("2"))
}
