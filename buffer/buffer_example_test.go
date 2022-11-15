package buffer_test

import (
	"fmt"

	"github.com/fufuok/bytespool/buffer"
)

func ExampleBuffer() {
	bb := buffer.Get()

	bb.SetString("1")
	_, _ = bb.WriteString("22")
	_, _ = bb.Write([]byte("333"))
	_ = bb.WriteByte('x')
	bb.Truncate(6)

	fmt.Printf("result=%q", bb.String())

	// After use, put Buffer back in the pool.
	buffer.Put(bb)

	// Output:
	// result="122333"
}
