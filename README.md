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
type BufPool struct{ ... }
    func NewBufPool(size int) *BufPool
type CapacityPools struct{ ... }
    func NewCapacityPools(minSize, maxSize int) *CapacityPools
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

### 🎨 Custom pools

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
go test -run=^$ -benchmem -benchtime=1s -count=2 -bench=.
goos: linux
goarch: amd64
pkg: github.com/fufuok/bytespool
cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
BenchmarkCapacityPools/New-4            56386340                21.24 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/New-4            56503125                21.21 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make-4           56200932                21.40 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make-4           56215285                21.43 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax-4        56522522                21.15 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax-4        56000730                21.45 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/New.Parallel-4           217137915                5.480 ns/op           0 B/op          0 allocs/op
BenchmarkCapacityPools/New.Parallel-4           212783748                5.912 ns/op           0 B/op          0 allocs/op
BenchmarkCapacityPools/Make.Parallel-4          212007224                5.541 ns/op           0 B/op          0 allocs/op
BenchmarkCapacityPools/Make.Parallel-4          211065468                5.583 ns/op           0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax.Parallel-4       217466509                5.525 ns/op           0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax.Parallel-4       218557538                5.524 ns/op           0 B/op          0 allocs/op
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