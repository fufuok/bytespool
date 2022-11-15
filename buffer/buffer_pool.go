package buffer

import (
	"sync"

	"github.com/fufuok/bytespool"
)

var defaultPools = &pools{
	bs: bytespool.DefaultCapacityPools,
}

var (
	// DefaultBufferSize is an initial allocation minimal capacity.
	DefaultBufferSize = 64
)

type pools struct {
	bs  *bytespool.CapacityPools
	buf sync.Pool
}

// SetCapacity initialize to the default byte slice pool.
// Divide into multiple pools according to the capacity scale.
// Maximum range of byte slice pool: [2,1<<31]
func SetCapacity(minSize, maxSize int) {
	defaultPools.bs = bytespool.NewCapacityPools(minSize, maxSize)
}

// Clone returns a copy of the Buffer.B.
// Atomically reset the reference count to 0.
func Clone(bb *Buffer) *Buffer {
	newBuf := NewBytes(bb.B)
	newBuf.RefReset()
	return newBuf
}

// Make return a Buffer with a byte slice of length 0.
// Capacity will not be 0, max(capacity, defaultPools.bs.MinSize())
func Make(capacity int) *Buffer {
	v := defaultPools.buf.Get()
	if v != nil {
		buf := v.(*Buffer)
		buf.B = defaultPools.bs.Make(capacity)
		buf.RefReset()
		return buf
	}
	return &Buffer{
		B: defaultPools.bs.Make(capacity),
		c: 0,
	}
}

func Make64(capacity uint64) *Buffer {
	return Make(int(capacity))
}

func MakeMax() *Buffer {
	return Make(defaultPools.bs.MaxSize())
}

func MakeMin() *Buffer {
	return Make(defaultPools.bs.MinSize())
}

func Get(capacity ...int) *Buffer {
	n := DefaultBufferSize
	if len(capacity) > 0 {
		n = capacity[0]
	}
	return Make(n)
}

// NewBytes returns a byte slice of the specified content.
func NewBytes(bs []byte) *Buffer {
	buf := Make(len(bs))
	buf.Set(bs)
	return buf
}

// NewString returns a byte slice of the specified content.
func NewString(s string) *Buffer {
	buf := Make(len(s))
	buf.SetString(s)
	return buf
}

// Similar to the built-in function to append elements to the end of a slice.
// If there is insufficient capacity,
// a new underlying array is allocated and the old array is reclaimed.
func appendBytes(buf []byte, elems ...byte) []byte {
	return defaultPools.bs.Append(buf, elems...)
}

func appendString(buf []byte, elems string) []byte {
	return defaultPools.bs.AppendString(buf, elems)
}

// Release put B back into the byte pool of the corresponding scale,
// and put the Buffer back into the buffer pool.
// Buffers smaller than the minimum capacity or larger than the maximum capacity are discarded.
func Release(bb *Buffer) (ok bool) {
	if bb.RefSwapDec() == 0 {
		ok = defaultPools.bs.Release(bb.B)
		bb.B = nil
		defaultPools.buf.Put(bb)
	}
	return
}

// Put is the same as b.Release.
func Put(bb *Buffer) {
	Release(bb)
}

func MinSize() int {
	return defaultPools.bs.MinSize()
}

func MaxSize() int {
	return defaultPools.bs.MaxSize()
}
