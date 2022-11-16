package buffer

import (
	"bytes"
	"errors"
	"testing"
)

var (
	testString    = "  Fufu 中　文\u2728->?\n*\U0001F63A   "
	testBytes     = []byte(testString)
	testStringLen = len(testString)
)

func TestBuffer_ZeroValue(t *testing.T) {
	var bb Buffer
	if bb.Len() != len(bb.B) || bb.Len() != 0 {
		t.Fatalf("except len=%d, got: %d", 0, bb.Len())
	}
	// bb.B is nil
	if bb.Cap() != cap(bb.B) || bb.Cap() != 0 {
		t.Fatalf("except len=%d, got: %d", 0, bb.Cap())
	}
	s := bb.String()
	if s != "" {
		t.Fatalf("unexpected result: %s", s)
	}
}

func TestBuffer_Base(t *testing.T) {
	bb := Get()
	if bb.Len() != len(bb.B) || bb.Len() != 0 {
		t.Fatalf("except len=%d, got: %d", 0, bb.Len())
	}
	if bb.Cap() != cap(bb.B) || bb.Cap() != DefaultBufferSize {
		t.Fatalf("except len=%d, got: %d", DefaultBufferSize, bb.Cap())
	}
	s := bb.String()
	if s != "" {
		t.Fatalf("unexpected result: %s", s)
	}
	bb.Set(testBytes)
	if bb.Len() != testStringLen || bb.String() != testString {
		t.Fatalf("unexpected result: %s", bb.String())
	}

	newBB := bb.Clone()
	n, err := bb.Write(testBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != testStringLen {
		t.Fatalf("except len=%d, got: %d", testStringLen, n)
	}
	if newBB.Len()*2 != bb.Len() {
		t.Fatalf("except newBB len=%d, got: %d", testStringLen, newBB.Len())
	}
	bb.Truncate(testStringLen)
	if !bytes.Equal(newBB.Bytes(), bb.B) {
		t.Fatalf("unexpected result: \n%s\n%s", bb.String(), newBB.String())
	}

	bb.Reset()
	if bb.Len() != 0 || bb.String() != "" {
		t.Fatalf("unexpected result: %s", bb.String())
	}
	if bb.Cap() < testStringLen*2 {
		t.Fatalf("except cap>=%d, got: %d", testStringLen*2, bb.Cap())
	}

	n, err = bb.WriteString("f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("except len=%d, got: %d", 1, n)
	}
	err = bb.WriteByte('f')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bb.String() != "ff" {
		t.Fatalf("unexpected result: %s", bb.String())
	}
	bb.Guarantee(4)
	if bb.Len() != 2 {
		t.Fatalf("except len=%d, got: %d", 2, bb.Len())
	}
	bb.Grow(4)
	if bb.Len() != 6 || bb.String() != "ffFufu" {
		t.Fatalf("unexpected result: %s", bb.String())
	}
	bb.Guarantee(400)
	if bb.Cap() < 402 {
		t.Fatalf("except cap=%d, got: %d", 402, bb.Cap())
	}
	if bb.Len() != 6 || bb.String() != "ffFufu" {
		t.Fatalf("unexpected result: %s", bb.String())
	}
}

func TestBuffer_RefAndRelease(t *testing.T) {
	bb := Get()
	if !bb.Release() {
		t.Fatal("expect to release the buffer successfully, but not")
	}
	bb = Get(2)
	bb.RefInc()
	bb.RefInc()
	bb.RefInc()
	c := bb.RefValue()
	if c != 3 {
		t.Fatalf("except reference counting=%d, got: %d", 3, c)
	}
	newBB := bb.Clone()
	newBS := bb.Copy()
	c = newBB.RefValue()
	if c != 0 {
		t.Fatalf("except reference counting=%d, got: %d", 3, c)
	}

	bb.Grow(2)
	if newBB.Len() == bb.Len() {
		t.Fatalf("except len=%d, got: %d", newBB.Len(), bb.Len())
	}
	if !bytes.Equal(newBB.Bytes(), newBS) {
		t.Fatalf("unexpected result: %s", newBS)
	}

	bb.RefStore(4)
	bb.RefDec()
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}
	// There will only be one success.
	if !bb.Release() {
		t.Fatal("expect to release the buffer successfully, but not")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}
	if bb.Release() {
		t.Fatal("expect to release the buffer failure, but not")
	}

	err := bb.Close()
	if !errors.Is(ErrClose, err) {
		t.Fatalf("unexpected result: %s", err)
	}
	err = newBB.Close()
	if err != nil {
		t.Fatalf("unexpected result: %s", err)
	}
}

func TestBuffer_Read(t *testing.T) {
	bb := New(10)
	copy(bb.B, testBytes)
	if !bytes.Equal(bb.B, testBytes[:10]) {
		t.Fatalf("unexpected result: %s", bb.String())
	}
	for i := 0; i <= 10; i++ {
		bs := make([]byte, i)
		n, err := bb.Read(bs)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if n != i {
			t.Fatalf("except len=%d, got: %d", i, n)
		}
		if !bytes.Equal(bs, testBytes[:i]) {
			t.Fatalf("unexpected result: %s", string(bs))
		}
	}
}

func TestBuffer_ReadFrom(t *testing.T) {
	prefix := "prefix"
	prefixLen := len(prefix)
	bb := NewString(prefix)
	for i := 0; i < 10; i++ {
		r := bytes.NewBufferString(testString)
		n, err := bb.ReadFrom(r)
		if int(n) != testStringLen {
			t.Fatalf("except n=%d, got: %d, i: %d", testStringLen, n, i)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		bufLen := bb.Len()
		expectedLen := prefixLen + (i+1)*testStringLen
		if bufLen != expectedLen {
			t.Fatalf("except length: %d, got: %d, i: %d", expectedLen, bufLen, i)
		}
		for j := 0; j < i; j++ {
			start := prefixLen + j*testStringLen
			b := bb.B[start : start+testStringLen]
			if string(b) != testString {
				t.Fatalf("except: %q, got: %q", testString, b)
			}
		}
	}
}

func TestBuffer_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	bb := NewBytes(testBytes)
	for i := 0; i < 10; i++ {
		n, err := bb.WriteTo(&buf)
		if int(n) != testStringLen {
			t.Fatalf("except n=%d, got: %d", testStringLen, n)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		s := buf.String()
		if s != testString {
			t.Fatalf("except: %q, got: %q", testString, s)
		}
		buf.Reset()
	}
}
