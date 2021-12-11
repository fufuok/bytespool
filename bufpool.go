package bytespool

// BufPool implements the httputil.BufferPool interface.
type BufPool struct {
	pool *CapacityPools
}

func NewBufPool(size int) *BufPool {
	return &BufPool{
		pool: NewCapacityPools(size, size),
	}
}

func (b *BufPool) Get() []byte {
	return b.pool.NewMin()
}

func (b *BufPool) Put(buf []byte) {
	b.pool.Put(buf)
}
