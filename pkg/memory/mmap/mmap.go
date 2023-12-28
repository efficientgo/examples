// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package mmap

import (
	"os"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/merrors"
	"golang.org/x/sys/unix"
)

// Wrapper for using memory mapping.
// Read more in "Efficient Go"; Example 5-1.

type MemoryMap struct {
	f *os.File // nil if anonymous.
	b []byte
}

func OpenFileBacked(path string, size int) (mf *MemoryMap, _ error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b, err := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, merrors.New(f.Close(), err).Err()
	}

	return &MemoryMap{f: f, b: b}, nil
}

func OpenAnonymous(size int) (mf *MemoryMap, _ error) {
	b, err := unix.Mmap(0, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return nil, err
	}
	return &MemoryMap{f: nil, b: b}, nil
}

func (f *MemoryMap) Close() error {
	errs := merrors.New()
	errs.Add(unix.Munmap(f.b))
	if f.f != nil {
		errs.Add(f.f.Close())
	}
	return errs.Err()
}

func (f *MemoryMap) Bytes() []byte { return f.b }

func (f *MemoryMap) File() *os.File { return f.f }

func (f *MemoryMap) Advise(advise int) error {
	if f.f != nil {
		// TODO(bwplotka): Provide table what works in SHARED mode.
		return errors.New("Most of madvise calls works ony on MAP_ANON mappings.")
	}
	if err := unix.Madvise(f.b, advise); err != nil {
		return err
	}
	return nil
}
