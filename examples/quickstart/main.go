package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	// Get() is the same as New()
	bs := bytespool.Get(1024)
	// len: 1024, cap: 1024
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))

	// Put() is the same as Release(), Put it back into the pool after use
	bytespool.Put(bs)

	// len: 0, capacity: 8 (Specified capacity)
	bs = bytespool.Make(8)
	bs = append(bs, "abc"...)
	// len: 3, cap: 8
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))
	ok := bytespool.Release(bs)
	// true
	fmt.Println(ok)

	// len: 8, capacity: 8 (Fixed length)
	bs = bytespool.New(8)
	copy(bs, "12345678")
	// len: 8, cap: 8, value: 12345678
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// len: len("xyz"), capacity: 4
	bs = bytespool.NewString("xyz")
	// len: 3, cap: 4, value: xyz
	fmt.Printf("len: %d, cap: %d, value: %s\n", len(bs), cap(bs), bs)
	bytespool.Release(bs)

	// Output:
	// len: 1024, cap: 1024
	// len: 3, cap: 8
	// true
	// len: 8, cap: 8, value: 12345678
	// len: 3, cap: 4, value: xyz
}
