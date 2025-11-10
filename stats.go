package bytespool

// SetWithStats enables or disables statistics collection.
// When enabled, statistics will be collected, but this may affect performance.
// When disabled (default), all atomic operations for statistics are skipped for better performance.
// This function is not thread-safe and should be called before any pool operations.
func SetWithStats(t bool) {
	DefaultCapacityPools.SetWithStats(t)
}

// GetWithStats returns the current status of statistics collection.
// When true, statistics are being collected.
// When false (default), statistics are not being collected.
func GetWithStats() bool {
	return DefaultCapacityPools.GetWithStats()
}

// RuntimeStats returns runtime statistics for byte pools.
// The statistics include:
// - NewBytes: total bytes newly allocated for pools
// - OutBytes: total bytes allocated outside pools
// - OutCount: total number of bytes allocated outside pools
// - ReusedBytes: total bytes reused from pools
//
// The statistics collection can be enabled/disabled with SetWithStats().
// When disabled (default), all counters will be zero.
func RuntimeStats(ps ...*CapacityPools) map[string]uint64 {
	p := DefaultCapacityPools
	if len(ps) > 0 {
		p = ps[0]
	}

	if !p.GetWithStats() {
		return nil
	}

	nb := p.getTotalNewBytes()
	ob := p.getTotalOutBytes()
	oc := p.getOutCount()
	rb := p.getTotalReusedBytes()
	return map[string]uint64{
		"NewBytes":    nb,
		"OutBytes":    ob,
		"OutCount":    oc,
		"ReusedBytes": rb,
	}
}

// RuntimeSummary is a structured summary of runtime pool statistics.
// It contains global byte counters and the top pools by reuse hits.
type RuntimeSummary struct {
	NewBytes    uint64     // total bytes newly allocated for pools
	OutBytes    uint64     // total bytes allocated outside pools
	OutCount    uint64     // total number of bytes allocated outside pools
	ReusedBytes uint64     // total bytes reused from pools
	TopPools    []PoolStat // top pools by reuse hits (ranked)
}

// RuntimeStatsSummary returns a structured RuntimeSummary for the provided
// CapacityPools (or the default pools when none provided).
func RuntimeStatsSummary(topN int, ps ...*CapacityPools) RuntimeSummary {
	p := DefaultCapacityPools
	if len(ps) > 0 {
		p = ps[0]
	}

	if !p.GetWithStats() {
		return RuntimeSummary{}
	}

	summary := RuntimeSummary{
		NewBytes:    p.getTotalNewBytes(),
		OutBytes:    p.getTotalOutBytes(),
		OutCount:    p.getOutCount(),
		ReusedBytes: p.getTotalReusedBytes(),
	}
	if topN > 0 {
		summary.TopPools = p.getPoolReuseStats(topN)
	}
	return summary
}

// PoolStat represents a pool statistic entry
type PoolStat struct {
	Rank      int
	Capacity  int
	ReuseHits uint64
}

// PoolReuseStats returns the top N pool reuse statistics (by reuse hits).
// If n <= 0 it returns an empty slice.
func PoolReuseStats(topN int, ps ...*CapacityPools) []PoolStat {
	p := DefaultCapacityPools
	if len(ps) > 0 {
		p = ps[0]
	}

	if topN <= 0 || !p.GetWithStats() {
		return nil
	}

	return p.getPoolReuseStats(topN)
}
