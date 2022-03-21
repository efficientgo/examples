package memory

import (
	"context"
	"fmt"
)

var bTest []byte

func AllocButNotAccess(ctx context.Context) error {
	// Allocate obsessively large amount of memory (10 GB), without accessing it (except first element), will use only few KBs of RSS.
	// Even though we never access is we can't allocate more than we have: fatal error: runtime: out of memory
	bTest := make([]byte, 10e9) // When used with global var, all is allocated.

	// Same if we generate the profile - somehow this triggers full allocation too.
	//if err := profiles.Heap("/shared/data/e2e-run"); err != nil {
	//	return err
	//}
	<-ctx.Done()
	fmt.Println(bTest[0])
	return nil
}
