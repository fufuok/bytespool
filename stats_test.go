package bytespool

import (
	"runtime/debug"
	"testing"
)

func TestRuntimeStats(t *testing.T) {
	var n, b, r uint64
	p := NewCapacityPools(2, 128)
	p.SetWithStats(true)
	gc := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(gc)

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

	stats := RuntimeStats(p)
	if stats["NewBytes"] != n {
		t.Fatalf("expect newBytes is %d, but got %d", n, stats["NewBytes"])
	}
	if stats["OutBytes"] != b {
		t.Fatalf("expect outBytes is %d, but got %d", b, stats["OutBytes"])
	}
	if stats["OutCount"] != 1 {
		t.Fatalf("expect outCount is %d, but got %d", 1, stats["OutCount"])
	}
	if stats["ReusedBytes"] != r {
		t.Fatalf("expect reusedBytes is %d, but got %d", r, stats["ReusedBytes"])
	}
	// Test MinSize and MaxSize in map
	if stats["MinSize"] != 2 {
		t.Fatalf("expect MinSize is 2, but got %d", stats["MinSize"])
	}
	if stats["MaxSize"] != 128 {
		t.Fatalf("expect MaxSize is 128, but got %d", stats["MaxSize"])
	}
}

func TestRuntimeStatsWithDefaultSettings(t *testing.T) {
	// Create a new pool to ensure clean state
	p := NewCapacityPools(2, 128)
	p.SetWithStats(true)
	gc := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(gc)

	// Perform some operations
	buf := p.Get(6)
	p.Put(buf)
	_ = p.Make(8)
	_ = p.Make(200)

	// Stats should have proper values when enabled
	stats := RuntimeStats(p)
	if stats["NewBytes"] == 0 {
		t.Fatalf("expect newBytes to be non-zero when stats enabled, but got %d", stats["NewBytes"])
	}
	if stats["OutBytes"] == 0 {
		t.Fatalf("expect outBytes to be non-zero when stats enabled, but got %d", stats["OutBytes"])
	}
	// Note: reusedBytes might be 0 if no reuse occurred in this test
}

func TestPoolReuseStatsN(t *testing.T) {
	pool := NewCapacityPools(8, 1024)
	pool.SetWithStats(true)
	for i := 0; i < 100; i++ {
		bs := pool.Make(100)
		pool.Release(bs)
	}

	stats := PoolReuseStats(5, pool)
	if len(stats) > 5 {
		t.Errorf("Expected at most 5 stats, got %d", len(stats))
	}
	for i := 1; i < len(stats); i++ {
		if stats[i-1].Rank != i || stats[i].Rank != i+1 {
			t.Error("Rankings are not sequential")
		}
	}
}

func TestRuntimeStatsSummary(t *testing.T) {
	pool := NewCapacityPools(8, 1024)
	pool.SetWithStats(true)
	stats := RuntimeStats(pool)
	requiredKeys := []string{"NewBytes", "OutBytes", "OutCount", "ReusedBytes", "MinSize", "MaxSize"}
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Missing required key: %s", key)
		}
	}
}

func TestRuntimeStatsSummaryStruct(t *testing.T) {
	pool := NewCapacityPools(8, 1024)
	pool.SetWithStats(true)

	summary := RuntimeStatsSummary(0, pool)
	// Test MinSize and MaxSize in struct
	if summary.MinSize != 8 {
		t.Fatalf("expect MinSize is 8, but got %d", summary.MinSize)
	}
	if summary.MaxSize != 1024 {
		t.Fatalf("expect MaxSize is 1024, but got %d", summary.MaxSize)
	}
	// Test that other fields are accessible
	_ = summary.NewBytes
	_ = summary.OutBytes
	_ = summary.OutCount
	_ = summary.ReusedBytes
	_ = summary.TopPools
}
