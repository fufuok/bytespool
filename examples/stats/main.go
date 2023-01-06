package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	// Custom pools
	bspool := bytespool.NewCapacityPools(8, 1024)
	for i := 0; i < 1000; i++ {
		bs := bspool.Make(i)
		bspool.Release(bs)
	}
	_ = bspool.Get(1025)
	_ = bspool.Get(8)
	stats := bytespool.RuntimeStats(bspool)
	fmt.Println(stats)

	// Default pools
	stats = bytespool.RuntimeStats()
	fmt.Println(stats)

	// Output:
	// map[Big:1025 New:2040 Reuse:671448]
	// map[Big:0 New:0 Reuse:0]
}
