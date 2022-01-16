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

var defaultCapacityPools = NewCapacityPools(defaultMinSize, defaultMaxSize)

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
	defaultCapacityPools = NewCapacityPools(minSize, maxSize)
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
	return &bytesPool{
		capacity: size,
		pool: sync.Pool{
			New: func() interface{} {
				buf := make([]byte, size, size)
				return &buf
			},
		},
	}
}

// Make return a byte slice of length 0.
func (p *CapacityPools) Make(capacity int) (buf []byte) {
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
// Length is size, may contain old data.
func (p *CapacityPools) New(size int) (buf []byte) {
	if size < 0 {
		size = 0
	}
	if size > p.maxSize {
		return make([]byte, size, size)
	}

	bp := p.getPool(size)
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

func (p *CapacityPools) NewBytes(bs []byte) []byte {
	buf := p.Make(len(bs))
	return append(buf, bs...)
}

func (p *CapacityPools) NewString(s string) []byte {
	buf := p.Make(len(s))
	return append(buf, s...)
}

func (p *CapacityPools) NewMax() []byte {
	return p.New(p.maxSize)
}

func (p *CapacityPools) NewMin() []byte {
	return p.New(p.minSize)
}

// Release put it back into the pool of the corresponding scale.
// Discard buffer larger than the maximum capacity.
func (p *CapacityPools) Release(buf []byte) bool {
	n := cap(buf)
	if n == 0 || n > p.maxSize {
		return false
	}
	bp := p.getPool(n)
	if n != bp.capacity {
		return false
	}
	// array pointer
	bp.pool.Put(unsafe.Pointer(&buf[:1][0]))
	return true
}

func (p *CapacityPools) Put(buf []byte) {
	p.Release(buf)
}

func (p *CapacityPools) getPool(size int) *bytesPool {
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

func getIndex(n int) int {
	return bits.Len32(uint32(n) - 1)
}

func Make(size int) []byte {
	return defaultCapacityPools.Make(size)
}

func Make64(size uint64) []byte {
	return defaultCapacityPools.Make64(size)
}

func MakeMax() []byte {
	return defaultCapacityPools.MakeMax()
}

func MakeMin() []byte {
	return defaultCapacityPools.MakeMin()
}

func New(size int) []byte {
	return defaultCapacityPools.New(size)
}

func Get(size int) []byte {
	return defaultCapacityPools.Get(size)
}

func New64(size uint64) []byte {
	return defaultCapacityPools.New64(size)
}

func NewBytes(bs []byte) []byte {
	return defaultCapacityPools.NewBytes(bs)
}

func NewString(s string) []byte {
	return defaultCapacityPools.NewString(s)
}

func NewMax() []byte {
	return defaultCapacityPools.NewMax()
}

func NewMin() []byte {
	return defaultCapacityPools.NewMin()
}

func Release(buf []byte) bool {
	return defaultCapacityPools.Release(buf)
}

func Put(buf []byte) {
	defaultCapacityPools.Put(buf)
}
