# 💫 BytesPool

Reuse used byte slices to achieve zero allocation.

The existing byte slices are stored in groups according to the capacity length range, and suitable byte slice objects are automatically allocated according to the capacity length when used.

## ✨ Features

- Customize the capacity range, or use the default pool.
- Get byte slices always succeed without panic.
- Optional length of 0 or fixed-length byte slices.
- Automatic garbage collection of big-byte slices.
- [BufPool](#-BufPool) implements the httputil.BufferPool interface.
- [Buffer](#-buffer) similar to bytes.Buffer, low-level byte slice multiplexing.
- High performance, See: [Benchmarks](#-benchmarks).

## ⚙️ Installation

```go
go get -u github.com/fufuok/bytespool
```

## 📚 Examples

Please see: [examples](examples)

Release warning: [examples/warning](examples/warning)

Simple reverse proxy: [examples/reverse_proxy](examples/reverse_proxy)

```go
package bytespool // import "github.com/fufuok/bytespool"

var DefaultCapacityPools = NewCapacityPools(defaultMinSize, defaultMaxSize)
func Append(buf []byte, elems ...byte) []byte
func AppendString(buf []byte, elems string) []byte
func Bytes(len, cap int) (b []byte)
func Clone(buf []byte) []byte
func Get(size int) []byte
func InitDefaultPools(minSize, maxSize int)
func Make(capacity int) []byte
func Make64(capacity uint64) []byte
func MakeMax() []byte
func MakeMin() []byte
func MaxSize() int
func MinSize() int
func New(size int) []byte
func New64(size uint64) []byte
func NewBytes(bs []byte) []byte
func NewMax() []byte
func NewMin() []byte
func NewString(s string) []byte
func Put(buf []byte)
func Release(buf []byte) bool
func RuntimeStats(ps ...*CapacityPools) map[string]uint64
func SetWithStats(t bool)
func GetWithStats() bool
type BufPool struct{ ... }
    func NewBufPool(size int) *BufPool
type CapacityPools struct{ ... }
    func NewCapacityPools(minSize, maxSize int) *CapacityPools
type PoolStat struct{ ... }
    func PoolReuseStatsN(topN int, ps ...*CapacityPools) []PoolStat
type RuntimeSummary struct{ ... }
    func RuntimeStatsSummary(topN int, ps ...*CapacityPools) RuntimeSummary
```

### ⚡️ Quickstart

```go
package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	// Get() is the same as New()
	bs := bytespool.Get(1024)
	// len: 1024, cap: 1024
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))

	// Put() is the same as Release(), Put it back into the pool after use
	bytespool.Put(bs)

	// len: 0, capacity: 8 (Specified capacity)
	bs = bytespool.Make(8)
	bs = append(bs, "abc"...)
	// len: 3, cap: 8
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	ok := bytespool.Release(bs)
	// true
	fmt.Println(ok)

	// len: 8, capacity: 8 (Fixed length)
	bs = bytespool.New(8)
	copy(bs, "12345678")
	// len: 8, cap: 8, value: 12345678
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// len: len("xyz"), capacity: 4
	bs = bytespool.NewString("xyz")
	// len: 3, cap: 4, value: xyz
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// Output:
	// len: 1024, cap: 1024
	// len: 3, cap: 8
	// true
	// len: 8, cap: 8, value: 12345678
	// len: 3, cap: 4, value: xyz
}
```

### ⏳ Automated reuse

```go
package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	// len: 0, cap: 4 (Specified capacity, automatically adapt to the capacity scale)
	bs3 := bytespool.Make(3)

	bs3 = append(bs3, "123"...)
	fmt.Printf("len: %d, cap: %d, %s\n", len(bs3), cap(bs3), bs3)

	bytespool.Release(bs3)

	// len: 4, cap: 4 (Fixed length)
	bs4 := bytespool.New(4)

	// Reuse of bs3
	fmt.Printf("same array: %v\n", &bs3[0] == &bs4[0])
	// Contain old data
	fmt.Printf("bs3: %s, bs4: %s\n", bs3, bs4[:3])

	copy(bs4, "xy")
	fmt.Printf("len: %d, cap: %d, %s\n", len(bs4), cap(bs4), bs4[:3])

	bytespool.Release(bs4)

	// Output:
	// len: 3, cap: 4, 123
	// same array: true
	// bs3: 123, bs4: 123
	// len: 4, cap: 4, xy3
}
```

### 🛠 Reset DefaultPools

```go
package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	bytespool.InitDefaultPools(512, 4096)

	bs := bytespool.Make(10)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	bs = bytespool.MakeMax()
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	bs = bytespool.New(10240)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	ok := bytespool.Release(bs)
	fmt.Printf("Discard: %v\n", !ok)

	// Output:
	// len: 0, cap: 512
	// len: 0, cap: 4096
	// len: 10240, cap: 10240
	// Discard: true
}
```

### 📊 Runtime Statistics

The library provides optional runtime statistics for monitoring byte slice usage:

```
// Check current statistics status
enabled := bytespool.GetWithStats() // Returns false by default

// Enable statistics (disabled by default for performance)
bytespool.SetWithStats(true)

// Get runtime statistics
stats := bytespool.RuntimeStats()
// Returns a map with keys:
// - "NewBytes": total bytes newly allocated for pools
// - "OutBytes": total bytes allocated outside pools
// - "OutCount": total number of bytes allocated outside pools
// - "ReusedBytes": total bytes reused from pools

// For custom pools
bspool := bytespool.NewCapacityPools(8, 1024)
stats = bytespool.RuntimeStats(bspool)
```

Note: Statistics are disabled by default to ensure maximum performance. Enable them only when needed for monitoring.

## 🎨 Custom pools

```go
package main

import (
	"github.com/fufuok/bytespool"
)

func main() {
	bspool := bytespool.NewCapacityPools(8, 1024)
	bs := bspool.MakeMax()
	bspool.Release(bs)
	bs = bspool.Make(64)
	bspool.Release(bs)
	bs = bspool.New(128)
	bspool.Release(bs)
}
```

### ♾ BufPool

Used to get fixed-length byte slices.

```go
package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	bufPool := bytespool.NewBufPool(32 * 1024)
	bs := bufPool.Get()

	data := []byte("test")
	n := copy(bs, data)
	// n: 4, bs: test
	fmt.Printf("n: %d, bs: %s\n", n, bs[:n])

	bufPool.Put(bs)
}
```

### 🔥 Buffer

Similar to bytes.Buffer, based on bytespool.

```go
package buffer // import "github.com/fufuok/bytespool/buffer"

var ErrTooLarge = errors.New("buffer: too large") ...
var DefaultBufferSize = 64
func GetReader(bs []byte) *bytes.Reader
func MaxSize() int
func MinSize() int
func Put(bb *Buffer)
func PutReader(r *bytes.Reader)
func Release(bb *Buffer) (ok bool)
func RuntimeStats() map[string]uint64
func SetCapacity(minSize, maxSize int)
type Buffer struct{ ... }
    func Clone(bb *Buffer) *Buffer
    func Get(capacity ...int) *Buffer
    func Make(capacity int) *Buffer
    func Make64(capacity uint64) *Buffer
    func MakeMax() *Buffer
    func MakeMin() *Buffer
    func New(size int) *Buffer
    func NewBuffer(buf []byte) *Buffer
    func NewBytes(bs []byte) *Buffer
    func NewString(s string) *Buffer
```

Please see:

- [DOC](buffer)
- [examples/buffer](examples/buffer)

```go
package main

import (
	"fmt"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.Get()

	bb.SetString("1")
	_, _ = bb.WriteString("22")
	_, _ = bb.Write([]byte("333"))
	_ = bb.WriteByte('x')
	bb.Truncate(6)

	fmt.Println("bb:", bb.String())

	bs := bb.Copy()
	bb.SetString("ff")
	fmt.Println("bs:", string(bs))
	fmt.Println("bb:", bb.String())

	// After use, put Buffer back in the pool.
	buffer.Put(bb)
	// or (safe)
	bb.Put()
	// or (safe)
	bb.Release()

	// Output:
	// bb: 122333
	// bs: 122333
	// bb: ff
}
```

## 🤖 Benchmarks

**byte slices**

```go
// # go test -run=^$ -benchmem -benchtime=1s -bench=.
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/bytespool
// cpu: AMD Ryzen 7 5700G with Radeon Graphics
// BenchmarkCapacityPools/New-16                   79853580                14.67 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-16                  82277400                15.29 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-16               87052480                13.77 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-16          450804770                2.481 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-16         523336443                2.406 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-16      597495465                1.909 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/New-16           72863752                15.69 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/Make-16          78273020                14.79 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/MakeMax-16       87188467                13.76 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/New.Parallel-16          34899675                35.38 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/Make.Parallel-16         33080406                35.70 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPoolsWithStatsEnabled/MakeMax.Parallel-16      34771252                33.53 ns/op            0 B/op          0 allocs/op
// BenchmarkParallelDirtBytes-16                                    7653496               140.6 ns/op          1024 B/op          1 allocs/op
// BenchmarkDirtBytes/size=1kb-16                                   9800404               121.7 ns/op          1024 B/op          1 allocs/op
// BenchmarkDirtBytes/size=3kb-16                                   5266185               232.2 ns/op          3072 B/op          1 allocs/op
// BenchmarkDirtBytes/size=5kb-16                                   2608358               426.5 ns/op          5376 B/op          1 allocs/op
// BenchmarkDirtBytes/size=7kb-16                                   1999676               585.7 ns/op          8192 B/op          1 allocs/op
// BenchmarkDirtBytes/size=9kb-16                                   5655252               211.8 ns/op          9472 B/op          1 allocs/op
// BenchmarkDirtBytes/size=11kb-16                                  2299012               515.4 ns/op         12288 B/op          1 allocs/op
// BenchmarkDirtBytes/size=13kb-16                                  3268006               368.2 ns/op         13568 B/op          1 allocs/op
// BenchmarkDirtBytes/size=15kb-16                                  1321743               923.3 ns/op         16384 B/op          1 allocs/op
// BenchmarkDirtBytes/size=17kb-16                                  3825888               306.5 ns/op         18432 B/op          1 allocs/op
// BenchmarkDirtBytes/size=19kb-16                                  2298004               525.0 ns/op         20480 B/op          1 allocs/op
// BenchmarkOriginBytes/size=1kb-16                                 6986295               176.7 ns/op          1024 B/op          1 allocs/op
// BenchmarkOriginBytes/size=3kb-16                                 3017343               394.2 ns/op          3072 B/op          1 allocs/op
// BenchmarkOriginBytes/size=5kb-16                                 1661990               717.0 ns/op          5376 B/op          1 allocs/op
// BenchmarkOriginBytes/size=7kb-16                                  899569              1243 ns/op            8192 B/op          1 allocs/op
// BenchmarkOriginBytes/size=9kb-16                                  955112              1121 ns/op            9472 B/op          1 allocs/op
// BenchmarkOriginBytes/size=11kb-16                                 740790              1502 ns/op           12288 B/op          1 allocs/op
// BenchmarkOriginBytes/size=13kb-16                                 727616              1635 ns/op           13568 B/op          1 allocs/op
// BenchmarkOriginBytes/size=15kb-16                                 490161              2079 ns/op           16384 B/op          1 allocs/op
// BenchmarkOriginBytes/size=17kb-16                                 521146              2052 ns/op           18432 B/op          1 allocs/op
// BenchmarkOriginBytes/size=19kb-16                                 469304              2387 ns/op           20480 B/op          1 allocs/op
// BenchmarkNormal4096Parallel-16                                     10000           6137806 ns/op        40960183 B/op      10001 allocs/op
// BenchmarkMCache4096Parallel-16                                     50416            110504 ns/op               0 B/op          0 allocs/op
// PASS
// ok      github.com/fufuok/bytespool     112.628s
```


**Buffer**

```go
go test -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/fufuok/bytespool/buffer
cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
BenchmarkBuffer_Write-4         72282802                16.06 ns/op            0 B/op          0 allocs/op
BenchmarkBuffer_Write_Std-4     65271292                18.50 ns/op            0 B/op          0 allocs/op
```







*ff*