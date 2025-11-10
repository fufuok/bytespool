package bytespool

import (
	"math"
	"math/bits"
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	minCapacity    = 2
	defaultMinSize = 2
	defaultMaxSize = 8 << 20 // 8 MiB
)

var DefaultCapacityPools = NewCapacityPools(defaultMinSize, defaultMaxSize)

type CapacityPools struct {
	pools       []*bytesPool
	minSize     int
	maxSize     int
	maxIndex    int
	decIndex    int
	newBytes    uint64 // New bytes allocated for pools
	outBytes    uint64 // Bytes allocated outside pools
	outCount    uint64 // Number of bytes allocated outside pools
	reusedBytes uint64 // Bytes reused from pools
	withStats   bool   // Controls whether to collect statistics for this pool
}

// bytesPool represents a pool for a specific capacity
type bytesPool struct {
	pool      sync.Pool
	capacity  int
	reuseHits uint64 // Number of times byte slices were reused from this pool
}

// InitDefaultPools initialize to the default pool.
func InitDefaultPools(minSize, maxSize int) {
	DefaultCapacityPools = NewCapacityPools(minSize, maxSize)
}

// NewCapacityPools divide into multiple pools according to the capacity scale.
// Maximum range of byte slice pool: [minCapacity,math.MaxInt32]
func NewCapacityPools(minSize, maxSize int) *CapacityPools {
	var pools []*bytesPool
	if maxSize > math.MaxInt32 {
		maxSize = math.MaxInt32
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

	mn := getIndex(minSize)
	mx := getIndex(maxSize)
	for i := mn; i <= mx; i++ {
		pools = append(pools, newBytesPool(1<<i))
	}

	return &CapacityPools{
		pools:     pools,
		minSize:   minSize,
		maxSize:   maxSize,
		maxIndex:  len(pools) - 1,
		decIndex:  mn,
		withStats: false,
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
		if p.withStats {
			atomic.AddUint64(&p.outCount, 1)
			atomic.AddUint64(&p.outBytes, uint64(size))
		}
		return Bytes(size, size)
	}

	ptr, _ := bp.pool.Get().(*byte)
	if ptr == nil {
		if p.withStats {
			atomic.AddUint64(&p.newBytes, uint64(bp.capacity))
		}
		return Bytes(size, bp.capacity)
	}

	if p.withStats {
		// per-pool reuse counters
		atomic.AddUint64(&bp.reuseHits, 1)
		atomic.AddUint64(&p.reusedBytes, uint64(bp.capacity))
	}

	// go1.20
	// return unsafe.Slice(ptr, bp.capacity)[:size]

	sh := (*bytesHeader)(unsafe.Pointer(&buf))
	sh.Data = ptr
	sh.Len = size
	sh.Cap = bp.capacity
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

	// go1.20, store array pointer,
	// bp.pool.Put(unsafe.SliceData(buf))

	sh := (*bytesHeader)(unsafe.Pointer(&buf))
	bp.pool.Put(sh.Data)
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

// SetWithStats enables or disables statistics collection for this pool.
// When enabled, statistics will be collected, but this may affect performance.
// When disabled (default), all atomic operations for statistics are skipped for better performance.
// This function is not thread-safe and should be called before any pool operations.
func (p *CapacityPools) SetWithStats(t bool) {
	p.withStats = t
}

// GetWithStats returns the current status of statistics collection for this pool.
// When true, statistics are being collected.
// When false (default), statistics are not being collected.
func (p *CapacityPools) GetWithStats() bool {
	return p.withStats
}

// getTotalNewBytes returns the sum of new bytes allocated across all pools
func (p *CapacityPools) getTotalNewBytes() uint64 {
	return atomic.LoadUint64(&p.newBytes)
}

// getTotalOutBytes returns the sum of bytes allocated outside pools
func (p *CapacityPools) getTotalOutBytes() uint64 {
	return atomic.LoadUint64(&p.outBytes)
}

// getOutCount returns the number of times bytes were allocated outside pools
func (p *CapacityPools) getOutCount() uint64 {
	return atomic.LoadUint64(&p.outCount)
}

// getTotalReusedBytes returns the sum of bytes reused from pools
func (p *CapacityPools) getTotalReusedBytes() uint64 {
	return atomic.LoadUint64(&p.reusedBytes)
}

// getPoolReuseStats returns reuse statistics for each pool capacity
func (p *CapacityPools) getPoolReuseStats(n int) []PoolStat {
	if n <= 0 {
		return nil
	}

	// collect non-zero reuse hits
	type kv struct {
		cap  int
		hits uint64
	}
	arr := make([]kv, 0, len(p.pools))
	for _, bp := range p.pools {
		if bp == nil {
			continue
		}
		v := atomic.LoadUint64(&bp.reuseHits)
		if v == 0 {
			continue
		}
		arr = append(arr, kv{cap: bp.capacity, hits: v})
	}

	if len(arr) == 0 {
		return nil
	}

	// partial sort: if arr is larger than n, use sort.Slice and then trim.
	sort.Slice(arr, func(i, j int) bool { return arr[i].hits > arr[j].hits })
	if len(arr) > n {
		arr = arr[:n]
	}

	stats := make([]PoolStat, 0, len(arr))
	for i, kv := range arr {
		stats = append(stats, PoolStat{Rank: i + 1, Capacity: kv.cap, ReuseHits: kv.hits})
	}
	return stats
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
