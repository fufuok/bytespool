package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.NewBytes([]byte("ff"))

	_ = (io.WriterTo)(bb)

	buf := bytes.NewBuffer(nil)
	n, err := bb.WriteTo(buf)
	fmt.Println(n, err, buf.String())

	// Output:
	// 2 <nil> ff
}
