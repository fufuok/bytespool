package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fufuok/bytespool"
)

func main() {
	// Custom pools
	bspool := bytespool.NewCapacityPools(8, 1024)
	bspool.SetWithStats(true)
	for i := 0; i < 1000; i++ {
		bs := bspool.Make(i)
		bspool.Release(bs)
	}
	_ = bspool.Get(1025)
	_ = bspool.Get(8)

	// Get global runtime stats
	stats := bytespool.RuntimeStats(bspool)
	fmt.Println("Runtime Stats:")
	// Sort keys for consistent output order
	keys := make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, stats[k])
	}

	// Get pool reuse stats
	poolStats := bytespool.PoolReuseStats(5, bspool)
	fmt.Println("Pool Reuse Stats:")
	if len(poolStats) == 0 {
		fmt.Println("  No pool reuse stats available")
	} else {
		for _, stat := range poolStats {
			fmt.Printf("  Rank %d: Capacity %d, ReuseHits %d times\n",
				stat.Rank, stat.Capacity, stat.ReuseHits)
		}
	}

	summary := bytespool.RuntimeStatsSummary(10, bspool)
	m, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(m))

	// Default pools
	bytespool.SetWithStats(true)
	stats = bytespool.RuntimeStats()
	fmt.Println("Default Pool Runtime Stats:")
	// Sort keys for consistent output order
	keys = make([]string, 0, len(stats))
	for k := range stats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, stats[k])
	}

	// Get default pool reuse stats
	poolStats = bytespool.PoolReuseStats(5)
	fmt.Println("Default Pool Reuse Stats:")
	if len(poolStats) == 0 {
		fmt.Println("  No pool reuse stats available")
	} else {
		for _, stat := range poolStats {
			fmt.Printf("  Rank %d: Capacity %d, ReuseHits %d times\n",
				stat.Rank, stat.Capacity, stat.ReuseHits)
		}
	}

	summary = bytespool.RuntimeStatsSummary(10)
	m, _ = json.MarshalIndent(summary, "", "  ")
	fmt.Println(string(m))

	// Output:
	// Runtime Stats:
	//  NewBytes: 2040
	//  OutBytes: 1025
	//  OutCount: 1
	//  ReusedBytes: 671448
	// Pool Reuse Stats:
	//  Rank 1: Capacity 1024, ReuseHits 486 times
	//  Rank 2: Capacity 512, ReuseHits 255 times
	//  Rank 3: Capacity 256, ReuseHits 127 times
	//  Rank 4: Capacity 128, ReuseHits 63 times
	//  Rank 5: Capacity 64, ReuseHits 31 times
	// {
	//  "NewBytes": 2040,
	//  "OutBytes": 1025,
	//  "OutCount": 1,
	//  "ReusedBytes": 671448,
	//  "TopPools": [
	//    {
	//      "Rank": 1,
	//      "Capacity": 1024,
	//      "ReuseHits": 486
	//    },
	//    {
	//      "Rank": 2,
	//      "Capacity": 512,
	//      "ReuseHits": 255
	//    },
	//    {
	//      "Rank": 3,
	//      "Capacity": 256,
	//      "ReuseHits": 127
	//    },
	//    {
	//      "Rank": 4,
	//      "Capacity": 128,
	//      "ReuseHits": 63
	//    },
	//    {
	//      "Rank": 5,
	//      "Capacity": 64,
	//      "ReuseHits": 31
	//    },
	//    {
	//      "Rank": 6,
	//      "Capacity": 32,
	//      "ReuseHits": 15
	//    },
	//    {
	//      "Rank": 7,
	//      "Capacity": 8,
	//      "ReuseHits": 9
	//    },
	//    {
	//      "Rank": 8,
	//      "Capacity": 16,
	//      "ReuseHits": 7
	//    }
	//  ]
	// }
	// Default Pool Runtime Stats:
	//  NewBytes: 0
	//  OutBytes: 0
	//  OutCount: 0
	//  ReusedBytes: 0
	// Default Pool Reuse Stats:
	//  No pool reuse stats available
	// {
	//  "NewBytes": 0,
	//  "OutBytes": 0,
	//  "OutCount": 0,
	//  "ReusedBytes": 0,
	//  "TopPools": null
	// }
}
