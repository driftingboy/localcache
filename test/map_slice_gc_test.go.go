package test

// 这是一个内存数据库，kv存储，为了优化 lru、或者手动淘汰带来的gc压力，使用了 map+slice（tree）的方式。
// 对于大slice，如果达到一定限额，会定时回收没有使用的内存，节省内存
// 否则留着，减少内存分配次数

// 测试 gc
type SMap struct {
	index map[int]int
	items []interface{}
}
