package bytespool

import "testing"

func TestBufPool(t *testing.T) {
	size := 1024
	bufPool := NewBufPool(size)
	buf := bufPool.Get()
	if len(buf) != size {
		t.Fatalf("expect buffer len == %d, but got %d", size, len(buf))
	}
	if cap(buf) != size {
		t.Fatalf("expect buffer cap == %d, but got %d", size, cap(buf))
	}

	buf[0] = 'f'
	bufPool.Put(buf)

	newBuf := bufPool.Get()
	if &buf[0] != &newBuf[0] {
		t.Fatal("expect buf and newBuf to be the same array")
	}
	if newBuf[0] != 'f' {
		t.Fatal("expect that newBuf may contain old data, but not")
	}
}
