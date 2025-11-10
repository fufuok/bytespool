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

func BenchmarkCapacityPoolsWithStatsEnabled(b *testing.B) {
	// Save the original state
	originalWithStats := GetWithStats()
	defer SetWithStats(originalWithStats) // Restore original state

	// Enable stats collection
	SetWithStats(true)

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
