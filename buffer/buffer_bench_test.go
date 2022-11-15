package buffer

import (
	"bytes"
	"testing"
)

func BenchmarkBuffer_Write(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		buf := Get()
		for pb.Next() {
			for i := 0; i < 10; i++ {
				if _, err := buf.Write(testBytes); err != nil {
					b.Fatal(err)
				}
			}
			buf.Reset()
		}
	})
}

func BenchmarkBuffer_Write_Std(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			for i := 0; i < 10; i++ {
				if _, err := buf.Write(testBytes); err != nil {
					b.Fatal(err)
				}
			}
			buf.Reset()
		}
	})
}

// go test -bench=. -benchmem
// goos: linux
// goarch: amd64
// pkg: github.com/fufuok/bytespool/buffer
// cpu: Intel(R) Xeon(R) Gold 6151 CPU @ 3.00GHz
// BenchmarkBuffer_Write-4         72282802                16.06 ns/op            0 B/op          0 allocs/op
// BenchmarkBuffer_Write_Std-4     65271292                18.50 ns/op            0 B/op          0 allocs/op
