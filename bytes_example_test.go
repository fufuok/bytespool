package bytespool_test

import (
	"fmt"

	"github.com/fufuok/bytespool"
)

func ExampleNew() {
	// Get() is the same as New()
	bs := bytespool.Get(1024)
	// len: 1024, cap: 1024
	fmt.Printf("len: %d, cap: %d\n", len(bs), cap(bs))

	// Put() is the same as Release(), Put it back into the pool after use
	bytespool.Put(bs)

	// len: 0, cap: 4 (Specified capacity, automatically adapt to the capacity scale)
	bs3 := bytespool.Make(3)

	bs3 = append(bs3, "123"...)
	fmt.Printf("len: %d, cap: %d, %s\n", len(bs3), cap(bs3), bs3)

	bytespool.Release(bs3)

	// len: 4, cap: 4 (Fixed length)
	bs4 := bytespool.New(4)

	// Reuse of bs3
	fmt.Printf("same array: %v\n", &bs3[0] == &bs4[0])
	// Contain old data
	fmt.Printf("bs3: %s, bs4: %s\n", bs3, bs4[:3])

	copy(bs4, "xy")
	fmt.Printf("len: %d, cap: %d, %s\n", len(bs4), cap(bs4), bs4[:3])

	bytespool.Release(bs4)

	// Output:
	// len: 1024, cap: 1024
	// len: 3, cap: 4, 123
	// same array: true
	// bs3: 123, bs4: 123
	// len: 4, cap: 4, xy3
}
