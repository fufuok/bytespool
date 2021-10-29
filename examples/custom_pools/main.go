package main

import (
	"github.com/fufuok/bytespool"
)

func main() {
	bspool := bytespool.NewCapacityPools(8, 1024)
	bs := bspool.MakeMax()
	bspool.Release(bs)
	bs = bspool.Make(64)
	bspool.Release(bs)
	bs = bspool.New(128)
	bspool.Release(bs)
}
