package mmap

import (
	"os"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type MemoryMap struct {
	f *os.File // Empty if anonymous.
	b []byte
}

func OpenFile(path string, size int) (mf *MemoryMap, _ error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b, err := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return &MemoryMap{f: f, b: b}, nil
}

func Open(size int) (mf *MemoryMap, _ error) {
	b, err := unix.Mmap(0, 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON)
	if err != nil {
		return nil, err
	}

	return &MemoryMap{f: nil, b: b}, nil
}

func (f *MemoryMap) Close() error {
	if err := unix.Munmap(f.b); err != nil {
		if f.f != nil {
			_ = f.f.Close()
		}
		return err
	}

	if f.f != nil {
		return f.f.Close()
	}
	return nil
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
