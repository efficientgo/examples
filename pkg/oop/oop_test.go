package block

import (
	"testing"
	"time"

"github.com/efficientgo/core/testutil"
)

func TestCompactor(t *testing.T) {
	now := time.Now()

	blocks := []Block{
		{id: "newer", start: now.Add(-2 * time.Hour), end: now},
		{id: "older", start: now.Add(-4 * time.Hour), end: now.Add(-2 * time.Hour)},
	}

	testutil.Equals(t,
		"older,newer: "+blocks[1].start.Format(time.RFC3339)+"-"+blocks[0].end.Format(time.RFC3339),
		Compact(blocks).String(),
	)
}
