package main

import (
	"encoding/json"
	"fmt"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	buffer.SetWithStats(true)
	for i := 0; i < 1000; i++ {
		bs := buffer.Make(10)
		buffer.Release(bs)
	}
	_ = buffer.Get(8)
	stats := buffer.RuntimeStatsSummary(10)
	js, _ := json.Marshal(stats)
	fmt.Println(string(js))

	// Output:
	// {"NewBytes":24,"OutBytes":0,"OutCount":0,"ReusedBytes":15984,"TopPools":[{"Rank":1,"Capacity":16,"ReuseHits":999}]}
}
