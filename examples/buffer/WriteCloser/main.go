package main

import (
	"fmt"
	"io"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.Get()

	_ = (io.WriteCloser)(bb)

	p := []byte("ff")
	n, err := bb.Write(p)
	fmt.Println(n, err)

	bs := make([]byte, 2)
	n = bb.CopyTo(bs)
	fmt.Println(n, string(bs))

	err = bb.Close()
	fmt.Println(err)

	// Output:
	// 2 <nil>
	// 2 ff
	// <nil>
}
