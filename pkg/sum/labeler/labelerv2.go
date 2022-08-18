package main

import (
	"context"
	"sync"

	"github.com/efficientgo/examples/pkg/sum"
	"github.com/efficientgo/tools/core/pkg/errcapture"
	bktpool "github.com/gobwas/pool"
	"github.com/thanos-io/objstore"
)

const bufRatio = 64

type labeler struct {
	bkt objstore.BucketReader

	pool         sync.Pool
	bucketedPool *bktpool.Pool
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

	buf := make([]byte, int(a.Size/bufRatio))
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

	bufSize := int(a.Size / bufRatio)
	buf := l.pool.Get().([]byte)
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	}
	defer func() { l.pool.Put(buf) }()

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

	bufSize := int(a.Size / bufRatio)
	bufI, _ := l.bucketedPool.Get(bufSize)
	buf := bufI.([]byte)
	if cap(buf) < bufSize {
		buf = make([]byte, bufSize)
	}
	defer func() { l.bucketedPool.Put(buf, cap(buf)) }()

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

func (l *labeler) labelObject4(ctx context.Context, objID string) (_ label, err error) {
	if _, err := l.bkt.Attributes(ctx, objID); err != nil {
		return label{}, err
	}

	rc, err := l.bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	defer errcapture.Do(&err, rc.Close, "close stream")

	s, err := sum.Sum6Reader(rc, l.buf)
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
