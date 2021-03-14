package prealloc_test

import (
	"testing"

	"github.com/efficientgo/examples/pkg/prealloc"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestPrealloc(t *testing.T) {
	n := int(1e4)
	wp := prealloc.CreateSlice(n)
	testutil.Equals(t, n*7, len(wp))
	testutil.Equals(t, n*7+5776, cap(wp))

	p := prealloc.CreateSlice2(n)
	testutil.Equals(t, wp, p)
	testutil.Equals(t, n*7, len(p))
	testutil.Equals(t, len(p), cap(p))
}
