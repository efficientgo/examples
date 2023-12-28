// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package block

import (
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/google/uuid"
)

func TestCompactor(t *testing.T) {
	now := time.Now()

	block1 := Block{id: uuid.New(), start: now.Add(-2 * time.Hour), end: now}
	block2 := Block{id: uuid.New(), start: now.Add(-4 * time.Hour), end: now.Add(-2 * time.Hour)}

	compacted := Compact(block1, block2)
	testutil.Equals(t, 4*time.Hour, compacted.Duration())
	testutil.Equals(t,
		compacted.id.String()+": "+block2.start.Format(time.RFC3339)+"-"+block1.end.Format(time.RFC3339),
		compacted.String(),
	)
}
