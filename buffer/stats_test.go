package buffer

import (
	"runtime/debug"
	"testing"

	"github.com/fufuok/bytespool"
)

func TestRuntimeStats(t *testing.T) {
	defer func() {
		defaultPools.bs = bytespool.DefaultCapacityPools
	}()
	var n, b, r uint64
	SetCapacity(2, 128)
	gc := debug.SetGCPercent(-1)

	n += 8
	bb := Get(6)
	Put(bb)

	r += 8
	_ = Make(8)

	n += 64
	bb = Get(63)
	Put(bb)

	b += 200
	_ = Make(200)

	debug.SetGCPercent(gc)
	stats := RuntimeStats()
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
