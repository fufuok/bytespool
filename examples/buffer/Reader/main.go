package main

import (
	"fmt"
	"io"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.NewBytes([]byte("abc"))

	// *bytes.Reader
	rb := bb.GetReader()

	_ = (io.Reader)(rb)
	_ = (io.ByteReader)(rb)

	a, err := rb.ReadByte()
	fmt.Println(string(a), err)

	bs := make([]byte, 10)
	n, err := rb.Read(bs)
	fmt.Println(n, err, string(bs))

	// Output:
	// a
	// 2 <nil> bc
}
