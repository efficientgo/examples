package getter_test

import (
	"errors"
	"testing"

	"github.com/efficientgo/examples/pkg/getter"
	"github.com/efficientgo/tools/core/pkg/testutil"
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

func TestClosure(t *testing.T) {
	r := &testReporter{}

	ratio := getter.FailureRatio(r)
	testutil.Equals(t, 0., ratio)
	ratio = getter.FailureRatio2(r)
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
	ratio = getter.FailureRatio2(r)
	testutil.Equals(t, 3/4., ratio)

}
