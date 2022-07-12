package main

import (
	"testing"

	"github.com/efficientgo/e2e"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestCalculator_Sum(t *testing.T) {
	d, err := e2e.NewDockerEnvironment("calculator")
	testutil.Ok(t, err)
	t.Cleanup(d.Close)

	mon, err := e2emonitoring.Start(d)
	testutil.Ok(t, err)
}
