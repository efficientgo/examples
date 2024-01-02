// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"bytes"
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/gobwas/pool/pbytes"
	"github.com/thanos-io/objstore"
)

func bench1(b *testing.B, labelFn func(ctx context.Context, objID string) (label, error)) {
	b.ReportAllocs()

	ctx := context.Background()
	var err error

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			_, err = labelFn(ctx, "10M.txt")
			testutil.Ok(b, err)
			continue
		}
		_, err = labelFn(ctx, "100M.txt")
		testutil.Ok(b, err)
	}
}

func bench2(b *testing.B, labelFn func(ctx context.Context, objID string) (label, error)) {
	b.ReportAllocs()

	ctx := context.Background()
	var err error

	wg := sync.WaitGroup{}
	wg.Add(4)

	b.ResetTimer()
	for g := 0; g < 4; g++ {
		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				_, err = labelFn(ctx, "10M.txt")
				testutil.Ok(b, err)
				_, err = labelFn(ctx, "100M.txt")
				testutil.Ok(b, err)
				runtime.GC()
			}
		}()
	}
	wg.Wait()

}

// BenchmarkLabeler recommended run options:
// $ export ver=v1 && go test -run '^$' -bench '^BenchmarkLabeler' -benchtime 100x -count 6 -cpu 1 -benchmem -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof | tee ${ver}.txt
func BenchmarkLabeler(b *testing.B) {
	ctx := context.Background()
	bkt := objstore.NewInMemBucket()

	buf := bytes.Buffer{}
	_, err := sumtestutil.CreateTestInputWithExpectedResult(&buf, 1e7)
	testutil.Ok(b, err)
	testutil.Ok(b, bkt.Upload(ctx, "10M.txt", &buf))

	buf.Reset()
	_, err = sumtestutil.CreateTestInputWithExpectedResult(&buf, 1e8)
	testutil.Ok(b, err)
	testutil.Ok(b, bkt.Upload(ctx, "100M.txt", &buf))

	b.Run("labelObject1", func(b *testing.B) {
		l := &labeler{bkt: bkt}

		bench1(b, l.labelObject1)
	})
	b.Run("labelObject2", func(b *testing.B) {
		l := &labeler{bkt: bkt}
		l.pool.New = func() any { return []byte(nil) }

		bench1(b, l.labelObject2)
	})
	b.Run("labelObject3", func(b *testing.B) {
		l := &labeler{bkt: bkt}

		l.bucketedPool = pbytes.New(1e3, 10e6)
		bench1(b, l.labelObject3)
	})
	b.Run("labelObject4", func(b *testing.B) {
		l := &labeler{bkt: bkt}

		bench1(b, l.labelObject4)
	})
}

func TestLabeler(t *testing.T) {
	ctx := context.Background()
	bkt := objstore.NewInMemBucket()

	buf := bytes.Buffer{}
	exp1, err := sumtestutil.CreateTestInputWithExpectedResult(&buf, 2e6)
	testutil.Ok(t, err)
	testutil.Ok(t, bkt.Upload(ctx, "2M.txt", &buf))

	buf.Reset()
	exp2, err := sumtestutil.CreateTestInputWithExpectedResult(&buf, 1e7)
	testutil.Ok(t, err)
	testutil.Ok(t, bkt.Upload(ctx, "100M.txt", &buf))

	t.Run("labelObjectNaive", func(t *testing.T) {
		l := &labeler{bkt: bkt}
		l.tmpDir = t.TempDir()

		ret, err := l.labelObject1(ctx, "2M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp1, ret.Sum)
		ret, err = l.labelObject1(ctx, "100M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp2, ret.Sum)
	})

	t.Run("labelObject1", func(t *testing.T) {
		l := &labeler{bkt: bkt}

		ret, err := l.labelObject1(ctx, "2M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp1, ret.Sum)
		ret, err = l.labelObject1(ctx, "100M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp2, ret.Sum)
	})
	t.Run("labelObject2", func(t *testing.T) {
		l := &labeler{bkt: bkt}
		l.pool.New = func() any { return []byte(nil) }

		ret, err := l.labelObject2(ctx, "2M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp1, ret.Sum)
		ret, err = l.labelObject2(ctx, "100M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp2, ret.Sum)
	})
	t.Run("labelObject3", func(t *testing.T) {
		l := &labeler{bkt: bkt}

		l.bucketedPool = pbytes.New(1e3, 10e6)

		ret, err := l.labelObject3(ctx, "2M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp1, ret.Sum)
		ret, err = l.labelObject3(ctx, "100M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp2, ret.Sum)
	})
	t.Run("labelObject4", func(t *testing.T) {
		l := &labeler{bkt: bkt}

		l.buf = make([]byte, 10e3)
		ret, err := l.labelObject4(ctx, "2M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp1, ret.Sum)
		ret, err = l.labelObject4(ctx, "100M.txt")
		testutil.Ok(t, err)
		testutil.Equals(t, exp2, ret.Sum)
	})
}
