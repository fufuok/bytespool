package bytespool

import (
	"sync/atomic"
)

func RuntimeStats(ps ...*CapacityPools) map[string]uint64 {
	p := DefaultCapacityPools
	if len(ps) > 0 {
		p = ps[0]
	}
	return map[string]uint64{
		"New":   atomic.LoadUint64(&p.newCounter),
		"Big":   atomic.LoadUint64(&p.bigCounter),
		"Reuse": atomic.LoadUint64(&p.reuseCounter),
	}
}
