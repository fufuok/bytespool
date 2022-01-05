package main

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func main() {
	x := bytespool.New(8)
	copy(x, "........")
	y := x[4:]

	// Wrong release!
	bytespool.Release(y)

	// You should release (x) and stop using (x)!
	// bytespool.Release(x)

	z := bytespool.New(4)
	x[6] = 'x'

	fmt.Printf("%p, %p, %p\n", &x, &y, &z)
	fmt.Printf("%p, %p, %p\n", &x[0], &y[0], &z[0])
	fmt.Printf("%s\n%s\n%s\n", x, y, z)

	// The output is similar to:
	// 0xc0000044c0, 0xc000004500, 0xc000004520
	// 0xc00000a1a0, 0xc00000a1a4, 0xc00000a1a4
	// ......x.
	// ..x.
	// ..x.
}
