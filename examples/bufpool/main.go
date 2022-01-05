package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	bufPool := bytespool.NewBufPool(32 * 1024)
	bs := bufPool.Get()

	data := []byte("test")
	n := copy(bs, data)
	// n: 4, bs: test
	fmt.Printf("n: %d, bs: %s\n", n, bs[:n])

	bufPool.Put(bs)
}
