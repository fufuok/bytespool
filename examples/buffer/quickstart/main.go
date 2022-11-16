package main

import (
	"fmt"

	"github.com/fufuok/bytespool/buffer"
)

func main() {
	bb := buffer.Get()

	bb.SetString("1")
	_, _ = bb.WriteString("22")
	_, _ = bb.Write([]byte("333"))
	_ = bb.WriteByte('x')
	bb.Truncate(6)

	fmt.Println("bb:", bb.String())

	bs := bb.Copy()
	bb.SetString("ff")
	fmt.Println("bs:", string(bs))
	fmt.Println("bb:", bb.String())

	// After use, put Buffer back in the pool.
	buffer.Put(bb)
	// or
	bb.Put()
	// or (safe)
	bb.Release()

	// Output:
	// bb: 122333
	// bs: 122333
	// bb: ff
}
