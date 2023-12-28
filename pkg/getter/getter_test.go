// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package getter_test

import (
	"errors"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/examples/pkg/getter"
)

type testReporter struct {
	r []getter.Report
}

func (r *testReporter) Get() []getter.Report {
	return r.r
}

type testReport struct {
	err error
}

func (r testReport) Error() error {
	return r.err
}

func TestFailureRatio(t *testing.T) {
	r := &testReporter{}

	ratio := getter.FailureRatio(r)
	testutil.Equals(t, 0., ratio)
	ratio = getter.FailureRatio_Better(r)
	testutil.Equals(t, 0., ratio)

	r.r = append(
		r.r,
		testReport{err: errors.New("a")},
		testReport{err: errors.New("b")},
		testReport{},
		testReport{err: errors.New("d")},
	)
	ratio = getter.FailureRatio(r)
	testutil.Equals(t, 3/4., ratio)
	ratio = getter.FailureRatio_Better(r)
	testutil.Equals(t, 3/4., ratio)
}
