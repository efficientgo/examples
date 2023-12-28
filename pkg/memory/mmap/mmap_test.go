// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package mmap

import (
	"fmt"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestMemoryMappedFile(t *testing.T) {
	f, err := OpenFileBacked("./test_file.txt", 20)
	testutil.Ok(t, err)

	t.Cleanup(func() {
		testutil.Ok(t, f.Close())
	})

	testutil.Equals(t, "is is a test stri", string(f.Bytes()[2:19]))
}

func TestMemoryMappedFileAppend(t *testing.T) {
	f, err := OpenFileBacked("./test_file.txt", 20)
	testutil.Ok(t, err)

	t.Cleanup(func() {
		testutil.Ok(t, f.Close())
	})

	b := f.Bytes()

	// Writing to b causes signal SIGSEGV: segmentation violation
	// b[2] = 'd'

	fmt.Println(len(b), cap(b))

	// Appending does not make much sense, since we replace memory mapped array by something on heap. But you can!
	b = append(b, '1')

	fmt.Println(len(b), cap(b))
}

func TestMemoryMappedAnnonymous(t *testing.T) {
	f, err := OpenAnonymous(20)
	testutil.Ok(t, err)

	t.Cleanup(func() {
		testutil.Ok(t, f.Close())
	})

	b := f.Bytes()
	b[2] = 'd'

	fmt.Println(len(b), cap(b))
}
