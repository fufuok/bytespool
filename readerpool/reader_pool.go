package readerpool

import (
	"bytes"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewReader(nil)
	},
}

func New(b []byte) *bytes.Reader {
	r := pool.Get().(*bytes.Reader)
	r.Reset(b)
	return r
}

func Release(r *bytes.Reader) {
	r.Reset(nil)
	pool.Put(r)
}
