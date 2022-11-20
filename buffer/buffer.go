package buffer

import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/fufuok/bytespool/readerpool"
)

var (
	ErrTooLarge   = errors.New("buffer: too large")
	ErrTruncation = errors.New("buffer: truncation out of range")
	ErrGrow       = errors.New("buffer: negative count")
	ErrClose      = errors.New("buffer: failed to add it to the pool")

	_ io.WriteCloser = (*Buffer)(nil)
	_ io.ReaderFrom  = (*Buffer)(nil)
	_ io.WriterTo    = (*Buffer)(nil)
)

// Buffer similar to bytes.Buffer, but provides finer-grained multiplexing of underlying byte slices.
// The zero value for Buffer is an empty buffer ready to use, but capacity will be 0.
// It is recommended to use pool to initialize a Buffer:
// e.g.::
//     bb := buffer.Get()    // The initial capacity is 64 (DefaultBufferSize)
//     bb := buffer.Make(8)  // The initial capacity is 8
// After use, put it back in the pool:
//     bb.Put()
//     bb.Release()
type Buffer struct {
	B []byte

	// Reference counting ensures that only one successfully puts it back into the pool.
	// c: stores the number of additional references,
	// c == 1: there are 2 references in total.
	c int64
}

// Clone returns a copy of the Buffer.B.
// Atomically reset the reference count to 0.
func (bb *Buffer) Clone() *Buffer {
	return Clone(bb)
}

func (bb *Buffer) Bytes() []byte {
	return bb.B
}

// Copy return a copy of the Buffer data.
func (bb *Buffer) Copy() []byte {
	return defaultPools.bs.NewBytes(bb.B)
}

// CopyTo same as copy(p, bb.B).
func (bb *Buffer) CopyTo(p []byte) int {
	return copy(p, bb.B)
}

// String implements print.Stringer.
//
// if the Buffer is a nil pointer, it returns "" instead of "<nil>"
func (bb *Buffer) String() string {
	return string(bb.B)
}

// UnsafeString not immutable.
func (bb *Buffer) UnsafeString() string {
	return *(*string)(unsafe.Pointer(&bb.B))
}

func (bb *Buffer) Len() int {
	return len(bb.B)
}

func (bb *Buffer) Cap() int {
	return cap(bb.B)
}

// Grow grows the internal buffer such that 'n' bytes can be written without reallocating.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (bb *Buffer) Grow(n int) {
	bb.Guarantee(n)
	bb.B = bb.B[:len(bb.B)+n]
}

// Guarantee buffer will be guaranteed to have at least 'n' remaining capacity.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (bb *Buffer) Guarantee(n int) {
	if n < 0 {
		panic(ErrGrow)
	}
	bLen := bb.Len()
	bCap := bb.Cap()
	bSize := bLen + n
	if bCap >= bSize {
		return
	}
	if bSize > math.MaxInt32 {
		panic(ErrTooLarge)
	}
	buf := defaultPools.bs.Make(bSize)
	buf = append(buf, bb.B...)
	defaultPools.bs.Release(bb.B)
	bb.B = buf
}

// Truncate buffer data, keep data of specified length.
// It panics if n is negative or greater than the length of the buffer.
func (bb *Buffer) Truncate(n int) {
	if n < 0 || n > bb.Len() {
		panic(ErrTruncation)
	}
	bb.B = bb.B[:n]
}

// Reset is the same as Truncate(0).
func (bb *Buffer) Reset() {
	bb.B = bb.B[:0]
}

// Write implements io.Writer.
//
// The function appends all the data in p to Buffer.B.
// The returned error is always nil.
func (bb *Buffer) Write(p []byte) (int, error) {
	bb.B = appendBytes(bb.B, p...)
	return len(p), nil
}

// Append appends all data in p to Buffer.B.
func (bb *Buffer) Append(p []byte) {
	bb.B = append(bb.B, p...)
}

// WriteByte implements io.ByteWriter.
//
// The function appends the byte c to Buffer.B.
// The returned error is always nil.
func (bb *Buffer) WriteByte(c byte) error {
	bb.B = appendBytes(bb.B, c)
	return nil
}

// AppendByte appends the byte c to Buffer.B.
func (bb *Buffer) AppendByte(c byte) {
	bb.B = append(bb.B, c)
}

// WriteString implements io.StringWriter.
//
// The function appends the s to Buffer.B.
// The returned error is always nil.
func (bb *Buffer) WriteString(s string) (int, error) {
	bb.B = appendString(bb.B, s)
	return len(s), nil
}

// AppendString appends the s to Buffer.B.
func (bb *Buffer) AppendString(s string) {
	bb.B = append(bb.B, s...)
}

// Set sets Buffer.B to p.
func (bb *Buffer) Set(p []byte) {
	bb.B = appendBytes(bb.B[:0], p...)
}

// SetString sets Buffer.B to s.
func (bb *Buffer) SetString(s string) {
	bb.B = appendString(bb.B[:0], s)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to Buffer.B.
func (bb *Buffer) ReadFrom(r io.Reader) (int64, error) {
	bLen := bb.Len()
	bCap := bb.Cap()
	n := bLen
	p := bb.B[:bCap]
	for {
		if n == bCap {
			if n == math.MaxInt32 {
				return int64(n), ErrTooLarge
			}
			bCap *= 2
			if bCap > math.MaxInt32 {
				bCap = math.MaxInt32
			}
			pNew := defaultPools.bs.New(bCap)
			copy(pNew, p)
			defaultPools.bs.Release(p)
			p = pNew
		}
		nn, err := r.Read(p[n:])
		n += nn
		if err != nil {
			bb.B = p[:n]
			n -= bLen
			if err == io.EOF {
				return int64(n), nil
			}
			return int64(n), err
		}
	}
}

// WriteTo implements io.WriterTo.
func (bb *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(bb.B)
	return int64(n), err
}

// Close implements io.Closer.
func (bb *Buffer) Close() error {
	if Release(bb) {
		return nil
	}
	return ErrClose
}

// Release put B back into the byte pool of the corresponding scale,
// and put the Buffer back into the buffer pool.
// Buffers smaller than the minimum capacity or larger than the maximum capacity are discarded.
func (bb *Buffer) Release() bool {
	return Release(bb)
}

// Put is the same as b.Release.
func (bb *Buffer) Put() {
	Put(bb)
}

// RefInc atomically increment the reference count by 1.
func (bb *Buffer) RefInc() {
	bb.RefAdd(1)
}

// RefDec atomically decrement the reference count by 1.
func (bb *Buffer) RefDec() {
	bb.RefAdd(-1)
}

func (bb *Buffer) RefAdd(delta int64) {
	atomic.AddInt64(&bb.c, delta)
}

// RefStore atomically stores val into the reference count.
func (bb *Buffer) RefStore(val int64) {
	atomic.StoreInt64(&bb.c, val)
}

// RefValue atomically loads the reference count.
func (bb *Buffer) RefValue() int64 {
	return atomic.LoadInt64(&bb.c)
}

// RefSwapDec atomically decrement the reference count by 1 and return the old value.
func (bb *Buffer) RefSwapDec() (c int64) {
	for {
		c = atomic.LoadInt64(&bb.c)
		if atomic.CompareAndSwapInt64(&bb.c, c, c-1) {
			return
		}
	}
}

// RefReset atomically reset the reference count to 0.
func (bb *Buffer) RefReset() {
	atomic.StoreInt64(&bb.c, 0)
}

// GetReader returns an io.Reader with bb.B.
func (bb *Buffer) GetReader() *bytes.Reader {
	bb.RefInc()
	return readerpool.New(bb.B)
}

// PutReader put an io.Reader into the pool.
func (bb *Buffer) PutReader(r *bytes.Reader) bool {
	readerpool.Release(r)
	return bb.Release()
}

// PutAll put io.Reader and Buffer into the pool.
func (bb *Buffer) PutAll(r *bytes.Reader) {
	readerpool.Release(r)
	defaultPools.bs.Release(bb.B)
	bb.B = nil
	defaultPools.buf.Put(bb)
}
