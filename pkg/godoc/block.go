// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

// Package block contains common functionality for interacting with TSDB blocks
// in the context of Thanos.
package block

import (
	"context"
	"fmt"

	"github.com/oklog/ulid"
)

const (
	// MetaFilename is the known JSON filename for meta information.
	MetaFilename = "meta.json"
)

// Download downloads directory that is mean to be block directory. If any of the files
// have a hash calculated in the meta file and it matches with what is in the destination path then
// we do not download it. We always re-download the meta file.
// BUG(bwplotka): No known bugs, but if there was one, it would be outlined here.
func Download(ctx context.Context, id ulid.ULID, dst string) error {
	fmt.Println("downloaded")
	return nil
}

// cleanUp cleans the partially uploaded files.
func cleanUp(ctx context.Context, id ulid.ULID) error {
	fmt.Println("cleaned")
	return nil
}
