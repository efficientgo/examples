package memory

import (
	"context"
	"fmt"
	"unsafe"
)

func Run(ctx context.Context) error {
	// Allocate obsessively large amount of memory (100 GB), without accessing it (except first element), will use only few KBs of RSS.
	b := make([]byte, 100e9)
	fmt.Println(unsafe.Sizeof(b[0]))

	<-ctx.Done()
	fmt.Println(unsafe.Sizeof(b[0]))
	return nil
}
