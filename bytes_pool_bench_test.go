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
// BenchmarkCapacityPools/New-4            56386340                21.24 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New-4            56503125                21.21 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4           56200932                21.40 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make-4           56215285                21.43 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4        56522522                21.15 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax-4        56000730                21.45 ns/op            0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           217137915                5.480 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/New.Parallel-4           212783748                5.912 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          212007224                5.541 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/Make.Parallel-4          211065468                5.583 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       217466509                5.525 ns/op           0 B/op          0 allocs/op
// BenchmarkCapacityPools/MakeMax.Parallel-4       218557538                5.524 ns/op           0 B/op          0 allocs/op
