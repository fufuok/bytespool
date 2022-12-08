package bytespool

import (
	"runtime/debug"
	"testing"
)

func TestRuntimeStats(t *testing.T) {
	var n, b, r uint64
	p := NewCapacityPools(2, 128)
	gc := debug.SetGCPercent(-1)

	n += 8
	buf := p.Get(6)
	p.Put(buf)

	r += 8
	_ = p.Make(8)

	n += 64
	buf = p.Get(63)
	p.Put(buf)

	b += 200
	_ = p.Make(200)

	debug.SetGCPercent(gc)
	stats := RuntimeStats(p)
	if stats["New"] != n {
		t.Fatalf("expect newCounter is %d, but got %d", n, stats["New"])
	}
	if stats["Big"] != b {
		t.Fatalf("expect bigCounter is %d, but got %d", n, stats["Big"])
	}
	if stats["Reuse"] != r {
		t.Fatalf("expect reuseCounter is %d, but got %d", n, stats["Reuse"])
	}
}
