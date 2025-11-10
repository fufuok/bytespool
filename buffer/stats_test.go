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
	SetWithStats(true)
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
	if stats["NewBytes"] != n {
		t.Fatalf("expect newBytes is %d, but got %d", n, stats["NewBytes"])
	}
	if stats["OutBytes"] != b {
		t.Fatalf("expect outBytes is %d, but got %d", n, stats["OutBytes"])
	}
	if stats["ReusedBytes"] != r {
		t.Fatalf("expect reusedBytes is %d, but got %d", n, stats["ReusedBytes"])
	}
}
