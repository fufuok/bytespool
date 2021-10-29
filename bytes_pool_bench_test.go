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
			bs := MakeMax()
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
				bs := MakeMax()
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
// BenchmarkCapacityPools/New-4                     27477343                44.06 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New-4                     26587846                43.59 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4                    26847370                44.09 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4                    27189642                43.44 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4                 52711242                22.77 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4                 52321828                24.49 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           100000000                10.55 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           100000000                10.52 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          100000000                10.48 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          100000000                10.57 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       204589562                5.929 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       199840848                5.931 ns/op            0 B/op          0 allocs/op
