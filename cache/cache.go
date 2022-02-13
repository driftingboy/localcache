package cache

type Key interface {
	Key() string
}

type Value interface {
	BytesNum() int64
}

type Cache interface {
	Set(kKey, v Value)
	Get(k Key) (v Value, ok bool)
	Del(k Key)
	Clear()
}
