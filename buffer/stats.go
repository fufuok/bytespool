package buffer

import (
	"github.com/fufuok/bytespool"
)

func SetWithStats(t bool) {
	defaultPools.bs.SetWithStats(t)
}

func GetWithStats() bool {
	return defaultPools.bs.GetWithStats()
}

func RuntimeStats() map[string]uint64 {
	return bytespool.RuntimeStats(defaultPools.bs)
}

func RuntimeStatsSummary(topN int) bytespool.RuntimeSummary {
	return bytespool.RuntimeStatsSummary(topN, defaultPools.bs)
}
