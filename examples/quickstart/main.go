package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	// len: 0, capacity: 8192 (Default maximum)
	bs := bytespool.Make()

	// Use...
	bs = append(bs, "abc"...)
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)

	// Put it back into the pool after use
	bytespool.Release(bs)

	// len: 0, capacity: 8 (Specified capacity)
	bs = bytespool.Make(8)
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	bytespool.Release(bs)

	// len: 8, capacity: 8 (Fixed length)
	bs = bytespool.New(8)
	copy(bs, "12345678")
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// Output:
	// len: 3, cap: 8192, value: abc
	// len: 0, cap: 8
	// len: 8, cap: 8, value: 12345678
}
