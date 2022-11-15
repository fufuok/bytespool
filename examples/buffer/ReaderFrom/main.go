package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.NewString("ff")

	_ = (io.ReaderFrom)(bb)

	fakeReader := bytes.NewBufferString("12345")
	n, err := bb.ReadFrom(fakeReader)
	fmt.Println(n, err)
	fmt.Println(bb)

	// Output:
	// 5 <nil>
	// ff12345
}
