package main

import (
	"fmt"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	for i := 0; i < 1000; i++ {
		bs := buffer.Make(10)
		buffer.Release(bs)
	}
	_ = buffer.Get(8)
	stats := buffer.RuntimeStats()
	fmt.Println(stats)

	// Output:
	// map[Big:0 New:24 Reuse:15984]
}
