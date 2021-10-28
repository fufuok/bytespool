# 💫 BytesPool

Reuse used byte slices to achieve zero allocation.

The existing byte slices are stored in groups according to the capacity length range, and suitable byte slice objects are automatically allocated according to the capacity length when used.

## ✨ Features

- Customize the capacity range, or use the default pool.
- Get byte slices always succeed without panic.
- Optional length of 0 or fixed-length byte slices.
- Automatic garbage collection of big-byte slices.
- High performance, See: [Benchmarks](#-benchmarks).

## ⚙️ Installation

```go
go get -u github.com/fufuok/bytespool
```

## 📚 Examples

please see: [examples](examples)

```go
package bytespool // import "github.com/fufuok/bytespool"

func InitDefaultPools(minSize, maxSize int)
func Make(capacity ...int) []byte
func New(size int) []byte
func Release(buf []byte) bool
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
	// len: 0, capacity: 8192 (Default maximum)
	bs := bytespool.Make()

	// Use...
	bs = append(bs, "abc"...)
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)

	// Put it back into the pool after use
	bytespool.Release(bs)

	// len: 0, capacity: 8 (Specified capacity)
	bs = bytespool.Make(8)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	// len: 8, capacity: 8 (Fixed length)
	bs = bytespool.New(8)
	copy(bs, "12345678")
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// Output:
	// len: 3, cap: 8192, value: abc
	// len: 0, cap: 8
	// len: 8, cap: 8, value: 12345678
}
```

### ⏳ Automated reuse

```go
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
```

### 🛠 Reset DefaultPools

```go
bytespool.InitDefaultPools(512, 4096)

bs := bytespool.Make(10)
fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
bytespool.Release(bs)

bs = bytespool.Make()
fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
bytespool.Release(bs)

bs = bytespool.New(10240)
fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
ok := bytespool.Release(bs)
fmt.Printf("Discard: %v", !ok)

// Output:
// len: 0, cap: 512
// len: 0, cap: 4096
// len: 10240, cap: 10240
// Discard: true
```

### 🎨 Custom pools

```go
bspool := bytespool.NewCapacityPools(8, 1024)
bs := bspool.Make()
bspool.Release(bs)
bs = bspool.Make(64)
bspool.Release(bs)
bs = bspool.New(128)
bspool.Release(bs)
```

## 🤖 Benchmarks

```go
go test -run=^$ -benchmem -benchtime=1s -count=2 -bench=.
goos: linux
goarch: amd64
pkg: github.com/fufuok/bytespool
cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
BenchmarkCapacityPools/New-4                     27641402                43.38 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/New-4                     26251407                43.46 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make-4                    27594026                44.56 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make-4                    27791390                43.31 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax-4                 47187982                24.42 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax-4                 46407331                25.80 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/New.Parallel-4           100000000                11.11 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/New.Parallel-4           100000000                10.89 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make.Parallel-4          100000000                10.94 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/Make.Parallel-4          100000000                11.13 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax.Parallel-4       186346058                6.406 ns/op            0 B/op          0 allocs/op
BenchmarkCapacityPools/MakeMax.Parallel-4       183609024                6.290 ns/op            0 B/op          0 allocs/op
```







*ff*