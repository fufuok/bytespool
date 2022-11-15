package main

import (
	"fmt"
	"time"

	"github.com/fufuok/bytespool/buffer"
)

func a(bb *buffer.Buffer) {
	// Buffers are not put back into the pool.
	defer bb.Put()
	fmt.Println("a:", bb)
}

func b(bb *buffer.Buffer) {
	// Buffers are not put back into the pool.
	defer bb.Release()

	time.Sleep(50 * time.Millisecond)

	// Buffer are safe to use.
	_ = bb.WriteByte('F')
	fmt.Println("b:", bb)
}

func main() {
	bb := buffer.Get()

	defer func() {
		// Here, the Buffer will be put back into the pool.
		ok := buffer.Release(bb)
		fmt.Println(ok)
	}()

	_, _ = bb.WriteString("ff")

	// Increment 2 reference counts
	bb.RefAdd(2)
	go a(bb)
	go b(bb)

	time.Sleep(200 * time.Millisecond)
	fmt.Println("main:", bb)

	// Output:
	// a: ff
	// b: ffF
	// main: ffF
	// true
}
