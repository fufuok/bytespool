package bytespool

import (
	"testing"
)

func BenchmarkCapacityPools(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := New(1024)
			Release(bs)
		}
	})
	b.Run("Make", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := Make(1024)
			Release(bs)
		}
	})
	b.Run("MakeMax", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bs := Make()
			Release(bs)
		}
	})
	b.Run("New.Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs := New(1024)
				Release(bs)
			}
		})
	})
	b.Run("Make.Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs := Make(1024)
				Release(bs)
			}
		})
	})
	b.Run("MakeMax.Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs := Make()
				Release(bs)
			}
		})
	})
}

// go test -run=^$ -benchmem -benchtime=1s -count=2 -bench=.
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/bytespool
// cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
// BenchmarkCapacityPools/New-4                     27641402                43.38 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New-4                     26251407                43.46 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4                    27594026                44.56 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4                    27791390                43.31 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4                 47187982                24.42 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4                 46407331                25.80 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           100000000                11.11 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           100000000                10.89 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          100000000                10.94 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          100000000                11.13 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       186346058                6.406 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       183609024                6.290 ns/op            0 B/op          0 allocs/op
