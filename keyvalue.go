package localcache

// A ByteView holds an immutable view of bytes.
// Internally it wraps either a []byte or a string,
// but that detail is invisible to callers.
//
// A ByteView is meant to be used as a value type, not
// a pointer (like a time.Time).
type ByteView struct {
	b []byte
	// s string
}

// Len returns the view's length.
func (v ByteView) BytesNum() int64 {
	return int64(len(v.b))
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) Bytes() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
