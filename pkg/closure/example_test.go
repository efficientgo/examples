package closure_test

import (
	"errors"
	"testing"

	"github.com/efficientgo/examples/pkg/closure"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

type testReport struct {
	err error
}

func (r testReport) Error() error {
	return r.err
}

func TestClosure(t *testing.T) {
	var reports []closure.Report
	reportsPtr := &reports
	r := closure.New(func() []closure.Report { return *reportsPtr })

	_, err := r.FailureRatio()
	testutil.NotOk(t, err)

	*reportsPtr = append(
		reports,
		testReport{err: errors.New("a")},
		testReport{err: errors.New("b")},
		testReport{},
		testReport{err: errors.New("d")},
	)
	ratio, err := r.FailureRatio()
	testutil.Ok(t, err)
	testutil.Equals(t, 3/4., ratio)

}
