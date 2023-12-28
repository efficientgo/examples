// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"context"
	"crypto/sha256"
	"io"
	"os"
	"sync"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/examples/pkg/profile/fd"
	"github.com/efficientgo/examples/pkg/sum"
	"github.com/gobwas/pool/pbytes"
	"github.com/thanos-io/objstore"
)

type labelFunc func(ctx context.Context, objID string) (label, error)

func bufferSize(fileSize int) int {
	s := fileSize / 64
	if s < 10e3 {
		return 10e3
	}
	return s
}

type labeler struct {
	bkt objstore.BucketReader

	tmpDir       string
	pool         sync.Pool
	bucketedPool *pbytes.Pool
	buf          []byte
}

func (l *labeler) labelObject1(ctx context.Context, objID string) (_ label, err error) {
	a, err := l.bkt.Attributes(ctx, objID)
	if err != nil {
		return label{}, err
	}

	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	defer errcapture.Do(&err, rc.Close, "close stream")

	buf := make([]byte, bufferSize(int(a.Size)))
	s, err := sum.Sum6Reader(rc, buf)
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID: objID,
		Sum:   s,
		// ...
	}, nil
}

func (l *labeler) labelObjectNaive(ctx context.Context, objID string) (_ label, err error) {
	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	// Download file first.
	// TODO(bwplotka): This is naive for book purposes.
	f, err := fd.CreateTemp(l.tmpDir, "cached-*")
	if err != nil {
		return label{}, err
	}
	defer func() {
		_ = f.Close()
		_ = os.RemoveAll(f.Name())
	}()

	h := sha256.New()

	// Write to both checksum hash and file.
	if _, err := io.Copy(f, io.TeeReader(rc, h)); err != nil {
		return label{}, err
	}
	if err := rc.Close(); err != nil {
		return label{}, err
	}

	s, err := sum.Sum(f.Name())
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID:    objID,
		Sum:      s,
		CheckSum: h.Sum(nil),
		// ...
	}, nil
}

func (l *labeler) labelObject2(ctx context.Context, objID string) (_ label, err error) {
	a, err := l.bkt.Attributes(ctx, objID)
	if err != nil {
		return label{}, err
	}

	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	defer errcapture.Do(&err, rc.Close, "close stream")

	bufSize := bufferSize(int(a.Size))
	buf := l.pool.Get().([]byte)
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	}
	defer func() { l.pool.Put(buf) }()

	s, err := sum.Sum6Reader(rc, buf[:bufSize])
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID: objID,
		Sum:   s,
		// ...
	}, nil
}

func (l *labeler) labelObject3(ctx context.Context, objID string) (_ label, err error) {
	a, err := l.bkt.Attributes(ctx, objID)
	if err != nil {
		return label{}, err
	}

	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	defer errcapture.Do(&err, rc.Close, "close stream")

	bufSize := bufferSize(int(a.Size))
	buf := l.bucketedPool.Get(bufSize, bufSize)
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	}
	defer func() { l.bucketedPool.Put(buf) }()

	s, err := sum.Sum6Reader(rc, buf[:bufSize])
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID: objID,
		Sum:   s,
		// ...
	}, nil
}

func (l *labeler) labelObject4(ctx context.Context, objID string) (_ label, err error) {
	a, err := l.bkt.Attributes(ctx, objID)
	if err != nil {
		return label{}, err
	}

	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	defer errcapture.Do(&err, rc.Close, "close stream")

	bufSize := bufferSize(int(a.Size))
	if cap(l.buf) < bufSize {
		l.buf = make([]byte, bufSize)
	}
	s, err := sum.Sum6Reader(rc, l.buf[:bufSize])
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID: objID,
		Sum:   s,
		// ...
	}, nil
}
