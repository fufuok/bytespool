package bytespool

import (
	"math"
	"math/bits"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

const (
	minCapacity    = 2
	defaultMinSize = 2
	defaultMaxSize = 8 << 20 // 8 MiB
)

var DefaultCapacityPools = NewCapacityPools(defaultMinSize, defaultMaxSize)

type CapacityPools struct {
	minSize  int
	maxSize  int
	maxIndex int
	decIndex int
	pools    []*bytesPool
}

type bytesPool struct {
	capacity int
	pool     sync.Pool
}

// InitDefaultPools initialize to the default pool.
func InitDefaultPools(minSize, maxSize int) {
	DefaultCapacityPools = NewCapacityPools(minSize, maxSize)
}

// NewCapacityPools divide into multiple pools according to the capacity scale.
// Maximum range of byte slice pool: [minCapacity,1<<31]
func NewCapacityPools(minSize, maxSize int) *CapacityPools {
	var pools []*bytesPool
	if maxSize > math.MaxInt32 {
		maxSize = 1 << 31
	}
	if maxSize < minCapacity {
		maxSize = minCapacity
	}
	if minSize > maxSize {
		minSize = maxSize
	}
	if minSize < minCapacity {
		minSize = minCapacity
	}

	min := getIndex(minSize)
	max := getIndex(maxSize)
	for i := min; i <= max; i++ {
		pools = append(pools, newBytesPool(1<<i))
	}

	return &CapacityPools{
		minSize:  minSize,
		maxSize:  maxSize,
		maxIndex: len(pools) - 1,
		decIndex: min,
		pools:    pools,
	}
}

func newBytesPool(size int) *bytesPool {
	return &bytesPool{capacity: size}
}

// Clone return a copy of the byte slice
func (p *CapacityPools) Clone(buf []byte) []byte {
	return p.NewBytes(buf)
}

// Make return a byte slice of length 0.
func (p *CapacityPools) Make(capacity int) []byte {
	return p.New(capacity)[:0]
}

func (p *CapacityPools) Make64(capacity uint64) []byte {
	return p.New(int(capacity))[:0]
}

func (p *CapacityPools) MakeMax() []byte {
	return p.New(p.maxSize)[:0]
}

func (p *CapacityPools) MakeMin() []byte {
	return p.New(p.minSize)[:0]
}

// New return byte slice of the specified size.
// Warning: may contain old data.
// Warning: returned buf is never equal to nil
func (p *CapacityPools) New(size int) (buf []byte) {
	if size < 0 {
		size = 0
	}

	bp := p.getMakePool(size)
	if bp == nil {
		return make([]byte, size, size)
	}

	ptr, _ := bp.pool.Get().(unsafe.Pointer)
	if ptr == nil {
		return make([]byte, size, bp.capacity)
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(ptr)
	sh.Len = size
	sh.Cap = bp.capacity
	runtime.KeepAlive(ptr)
	return
}

func (p *CapacityPools) Get(size int) []byte {
	return p.New(size)
}

func (p *CapacityPools) New64(size uint64) []byte {
	return p.New(int(size))
}

func (p *CapacityPools) NewMax() []byte {
	return p.New(p.maxSize)
}

func (p *CapacityPools) NewMin() []byte {
	return p.New(p.minSize)
}

// NewBytes returns a byte slice of the specified content.
func (p *CapacityPools) NewBytes(bs []byte) []byte {
	buf := p.Make(len(bs))
	return append(buf, bs...)
}

// NewString returns a byte slice of the specified content.
func (p *CapacityPools) NewString(s string) []byte {
	buf := p.Make(len(s))
	return append(buf, s...)
}

// Append similar to the built-in function to append elements to the end of a slice.
// If there is insufficient capacity,
// a new underlying array is allocated and the old array is reclaimed.
func (p *CapacityPools) Append(buf []byte, elems ...byte) []byte {
	n := len(buf)
	c := cap(buf)
	m := n + len(elems)
	if c < m && c <= p.maxSize {
		bbuf := p.New(m)
		copy(bbuf, buf)
		copy(bbuf[n:], elems)
		p.Release(buf)
		return bbuf
	}
	return append(buf, elems...)
}

func (p *CapacityPools) AppendString(buf []byte, elems string) []byte {
	n := len(buf)
	c := cap(buf)
	m := n + len(elems)
	if c < m && c <= p.maxSize {
		bbuf := p.New(m)
		copy(bbuf, buf)
		copy(bbuf[n:], elems)
		p.Release(buf)
		return bbuf
	}
	return append(buf, elems...)
}

// Release put it back into the byte pool of the corresponding scale.
// Buffers smaller than the minimum capacity or larger than the maximum capacity are discarded.
func (p *CapacityPools) Release(buf []byte) bool {
	bp := p.getReleasePool(cap(buf))
	if bp == nil {
		return false
	}
	// array pointer
	bp.pool.Put(unsafe.Pointer(&buf[:1][0]))
	return true
}

func (p *CapacityPools) Put(buf []byte) {
	p.Release(buf)
}

func (p *CapacityPools) MinSize() int {
	return p.minSize
}

func (p *CapacityPools) MaxSize() int {
	return p.maxSize
}

func (p *CapacityPools) getMakePool(size int) *bytesPool {
	if size <= p.minSize {
		return p.pools[0]
	}
	if size == p.maxSize {
		return p.pools[p.maxIndex]
	}
	if size > p.maxSize {
		return nil
	}
	return p.pools[getIndex(size)-p.decIndex]
}

func (p *CapacityPools) getReleasePool(size int) *bytesPool {
	if size < p.minSize || size > p.maxSize {
		return nil
	}
	if size == p.minSize {
		return p.pools[0]
	}
	if size == p.maxSize {
		return p.pools[p.maxIndex]
	}
	idx := getIndex(size) - p.decIndex
	pool := p.pools[idx]
	if size < pool.capacity {
		pool = p.pools[idx-1]
	}
	return pool
}

func getIndex(n int) int {
	return bits.Len32(uint32(n) - 1)
}

func Clone(buf []byte) []byte {
	return DefaultCapacityPools.Clone(buf)
}

func Make(capacity int) []byte {
	return DefaultCapacityPools.Make(capacity)
}

func Make64(capacity uint64) []byte {
	return DefaultCapacityPools.Make64(capacity)
}

func MakeMax() []byte {
	return DefaultCapacityPools.MakeMax()
}

func MakeMin() []byte {
	return DefaultCapacityPools.MakeMin()
}

func New(size int) []byte {
	return DefaultCapacityPools.New(size)
}

func Get(size int) []byte {
	return DefaultCapacityPools.Get(size)
}

func New64(size uint64) []byte {
	return DefaultCapacityPools.New64(size)
}

func NewMax() []byte {
	return DefaultCapacityPools.NewMax()
}

func NewMin() []byte {
	return DefaultCapacityPools.NewMin()
}

func NewBytes(bs []byte) []byte {
	return DefaultCapacityPools.NewBytes(bs)
}

func NewString(s string) []byte {
	return DefaultCapacityPools.NewString(s)
}

func Append(buf []byte, elems ...byte) []byte {
	return DefaultCapacityPools.Append(buf, elems...)
}

func AppendString(buf []byte, elems string) []byte {
	return DefaultCapacityPools.AppendString(buf, elems)
}

func Release(buf []byte) bool {
	return DefaultCapacityPools.Release(buf)
}

func Put(buf []byte) {
	DefaultCapacityPools.Put(buf)
}

func MinSize() int {
	return DefaultCapacityPools.MinSize()
}

func MaxSize() int {
	return DefaultCapacityPools.MaxSize()
}
