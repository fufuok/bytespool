package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	bytespool.InitDefaultPools(512, 4096)

	bs := bytespool.Make(10)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	bs = bytespool.MakeMax()
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	bs = bytespool.New(10240)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	ok := bytespool.Release(bs)
	fmt.Printf("Discard: %v\n", !ok)

	// Output:
	// len: 0, cap: 512
	// len: 0, cap: 4096
	// len: 10240, cap: 10240
	// Discard: true
}
