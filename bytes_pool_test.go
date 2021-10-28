package bytespool

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"testing"
)

func TestCapacityPools(t *testing.T) {
	minSize := 64
	maxSize := 2048
	pools := NewCapacityPools(minSize, maxSize)
	tests := []struct {
		size        int
		scaleSize   int
		bytesLength int
		releaseOK   bool
	}{
		{0, 64, 0, true},
		{64, 64, 64, true},
		{128, 128, 128, true},
		{2000, 2048, 2000, true},
		{2047, 2048, 2047, true},
		{4096, 0, 4096, false},
		{5000, 0, 5000, false},
	}
	for _, v := range tests {
		t.Run(fmt.Sprintf("bytes.New(%d)", v.size), func(t *testing.T) {
			bp := pools.getPool(v.size)
			if bp == nil {
				if v.scaleSize > 0 {
					t.Fatalf("expect pool capacity is %d, but got nil", v.scaleSize)
				}
			} else if bp.capacity != v.scaleSize {
				t.Fatalf("expect pool capacity is %d, but got %d", v.scaleSize, bp.capacity)
			}

			buf := pools.New(v.size)
			if len(buf) != v.bytesLength {
				t.Fatalf("expect buffer len is %d, but got %d", v.bytesLength, len(buf))
			}
			if cap(buf) < v.scaleSize {
				t.Fatalf("expect buffer cap >= %d, but got %d", v.scaleSize, cap(buf))
			}

			buf = append(buf, '1')

			ok := pools.Release(buf)
			if ok != v.releaseOK {
				t.Fatalf("expect to release the buffer result is %v, but got %v", v.releaseOK, ok)
			}
		})
	}

	t.Run("bytes.Make()", func(t *testing.T) {
		bp := pools.pools[pools.maxIndex]
		if bp == nil {
			t.Fatalf("expect pool index is %d, but got nil", pools.maxIndex)
			return
		}
		buf := pools.Make()
		if len(buf) != 0 {
			t.Fatalf("expect buffer len is 0, but got %d", len(buf))
		}
		if cap(buf) < maxSize {
			t.Fatalf("expect buffer cap >= %d, but got %d", maxSize, cap(buf))
		}
	})
}

func TestCapacityPools_Boundary(t *testing.T) {
	pools := NewCapacityPools(0, 0)
	if pools.maxIndex != 0 {
		t.Fatal("expect have one pool, but not")
	}

	buf := pools.Make()
	if len(buf) != 0 {
		t.Fatalf("expect buffer len is 0, but got %d", len(buf))
	}
	if cap(buf) != minCapacity {
		t.Fatalf("expect buffer cap is %d, but got %d", minCapacity, cap(buf))
	}

	if !pools.Release(buf) {
		t.Fatal("expect to release the buffer successfully, but not")
	}

	buf = pools.New(3)
	if len(buf) != 3 {
		t.Fatalf("expect buffer len is 3, but got %d", len(buf))
	}
	if cap(buf) < 3 {
		t.Fatalf("expect buffer cap >= 3, but got %d", cap(buf))
	}

	if pools.Release(buf) {
		t.Fatal("expect to release the buffer failure, but not")
	}
}

func TestCapacityPools_Default(t *testing.T) {
	if defaultCapacityPools.maxIndex+1 != 13 {
		t.Fatalf("expect count default pools is 13, but got %d", defaultCapacityPools.maxIndex+1)
	}

	buf := Make()
	if len(buf) != 0 {
		t.Fatalf("expect buffer len is 0, but got %d", len(buf))
	}
	if cap(buf) != defaultMaxSize {
		t.Fatalf("expect buffer cap is %d, but got %d", defaultMaxSize, cap(buf))
	}

	abc := []byte("abc")
	buf = append(buf, abc...)

	Release(buf)

	newBuf := New(defaultMaxSize)
	if fmt.Sprintf("%p", newBuf) != fmt.Sprintf("%p", buf) {
		t.Fatal("expect the newBuf is the buf, but not")
	}
	if !bytes.Equal(abc, (newBuf)[:3]) {
		t.Fatal("expect that newBuf may contain old data, but not")
	}

	Release(newBuf)

	buf8 := New(8)
	copy(buf8, "12345678")
	if string(buf8) != "12345678" {
		t.Fatal("expect copy result is 123456789, but not")
	}

	buf8 = append(buf8, '9')

	// Disable GC to test re-acquire the same data
	gc := debug.SetGCPercent(-1)

	Release(buf8)

	buf16 := New(16)
	if &buf8[0] != &buf16[0] {
		t.Fatal("expect buf8 and buf16to be the same array")
	}
	if string(buf16[:9]) != "123456789" {
		t.Fatal("expect the buf8 is the buf16, but not")
	}

	// Re-enable GC
	debug.SetGCPercent(gc)

	minSize := 2
	maxSize := 8
	InitDefaultPools(minSize, maxSize)
	if defaultCapacityPools.maxIndex+1 != 3 {
		t.Fatal("expect count default pools is 3, but not")
	}

	buf = Make()
	if len(buf) != 0 {
		t.Fatalf("expect buffer len is 0, but got %d", len(buf))
	}
	if cap(buf) != maxSize {
		t.Fatalf("expect buffer cap is %d, but got %d", maxSize, cap(buf))
	}
	buf = Make(3)
	if len(buf) != 0 {
		t.Fatalf("expect buffer len is 0, but got %d", len(buf))
	}
	if cap(buf) != 4 {
		t.Fatalf("expect buffer cap is 4, but got %d", cap(buf))
	}
	buf = Make(33)
	if len(buf) != 0 {
		t.Fatalf("expect buffer len is 0, but got %d", len(buf))
	}
	if cap(buf) != 33 {
		t.Fatalf("expect buffer cap is 33, but got %d", cap(buf))
	}
	if Release(buf) {
		t.Fatal("expect to release the buffer failure, but not")
	}
	buf = append(buf, '1')
	if Release(buf) {
		t.Fatal("expect to release the buffer failure, but not")
	}

	InitDefaultPools(defaultMinSize, defaultMaxSize)
}
