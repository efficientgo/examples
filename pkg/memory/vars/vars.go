package vars

import (
	"fmt"
	"unsafe"
)

// Example showing differences between values pointers and special types used as function arguments.
// Read more in "Efficient Go"; Example 5-4.

func myFunction(
	arg1 int, arg2 *int,
	arg3 biggie, arg4 *biggie,
	arg5 []byte, arg6 *[]byte,
	arg7 chan byte, arg8 map[string]int, arg9 func(),
) {
	// ...
	fmt.Println(unsafe.Sizeof(arg3))
}

type biggie struct {
	huge  [1e8]byte
	other *biggie
}
