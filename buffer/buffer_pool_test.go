package buffer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/fufuok/bytespool"
)

func TestBufferPool_GetPut(t *testing.T) {
	for i := 0; i < 10; i++ {
		want := fmt.Sprintf("i: %d", i)
		bb := Get()
		_, _ = bb.WriteString(want)
		if bb.String() != want {
			t.Fatalf("expect result: %q, got: %q", want, bb.String())
		}
		Put(bb)

		bb = Make(4)
		bb.SetString("i")
		_, _ = bb.WriteString(":")
		_ = bb.WriteByte(' ')
		bb.B = appendString(bb.B, strconv.Itoa(i))
		if bb.String() != want {
			t.Fatalf("expect result: %q, got: %q", want, bb.String())
		}
		ok := bb.Release()
		if !ok {
			t.Fatal("expect to release the buffer successfully, but not")
		}
	}
}

func TestBufferPool_Boundary(t *testing.T) {
	bb := Get()
	if bb.Len() != 0 || bb.Cap() != DefaultBufferSize {
		t.Fatal("buffer initial error")
	}
	bb = Get(0)
	if bb.Len() != 0 || bb.Cap() != MinSize() {
		t.Fatal("buffer initial error")
	}
	bb = Make(3)
	if bb.Len() != 0 || bb.Cap() != 4 {
		t.Fatal("buffer initial error")
	}
	bb = Make64(3)
	if bb.Len() != 0 || bb.Cap() != 4 {
		t.Fatal("buffer initial error")
	}
	bb = NewBytes([]byte("f"))
	if bb.Len() != 1 || bb.Cap() != 2 {
		t.Fatal("buffer initial error")
	}
	bb = NewString("f")
	if bb.Len() != 1 || bb.Cap() != 2 {
		t.Fatal("buffer initial error")
	}
	_, _ = bb.WriteString("ff")
	if bb.Len() != 3 || bb.Cap() != 4 {
		t.Fatal("buffer initial error")
	}

	SetCapacity(0, 7)
	defer SetCapacity(bytespool.DefaultCapacityPools.MinSize(), bytespool.DefaultCapacityPools.MaxSize())
	if MinSize() != 2 || MaxSize() != 7 {
		t.Fatalf("expect minSize is 2, maxSize is 7, but got: %d, %d", MinSize(), MaxSize())
	}
	bb = Get()
	if bb.Len() != 0 || bb.Cap() != DefaultBufferSize {
		t.Fatal("buffer initial error")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}

	bb = MakeMin()
	if bb.Len() != 0 || bb.Cap() != 2 {
		t.Fatal("buffer initial error")
	}
	bb = MakeMax()
	if bb.Len() != 0 || bb.Cap() != 8 {
		t.Fatal("buffer initial error")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}

	SetCapacity(0, 8)
	bb = MakeMax()
	if bb.Len() != 0 || bb.Cap() != 8 {
		t.Fatal("buffer initial error")
	}
	if !bb.Release() {
		t.Fatal("expect to release the buffer successfully, but not")
	}
}
