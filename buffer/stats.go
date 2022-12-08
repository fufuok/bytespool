package buffer

import (
	"github.com/fufuok/bytespool"
)

func RuntimeStats() map[string]uint64 {
	return bytespool.RuntimeStats(defaultPools.bs)
}
