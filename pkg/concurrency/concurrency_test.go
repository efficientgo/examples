package concurrency

import (
	"sync"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestFunction(t *testing.T) {
	function()
}

func TestConcurrency(t *testing.T) {
	var mu sync.Mutex
	var num int64
	// Do not do that at home. Globals are bad, doing it so example is simpler (:
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()

		num += 10
		return num
	}

	testutil.Equals(t, int64(10+20+30), sharingWithAtomic())
	testutil.Equals(t, int64(40+50+60), sharingWithMutex())
	testutil.Equals(t, int64(70+80+90), sharingWithChannel())
	testutil.Equals(t, int64(100+110+120), sharingWithShardedSpace())
}
