package main

import (
	"context"
	"fmt"
	"os"
	execlib "os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/parquet-export/export1"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/pkg/errors"
)

var (
	generateDataPath = func() string { a, _ := filepath.Abs("generated"); return a }()
	maxTime          = `2021-07-20T00:00:00Z`
)

// Test args: -test.timeout 9999m
func TestParquetExport(t *testing.T) {
	t.Parallel()

	// Create 10k series for 1w of TSDB blocks. Cache them to 'generated' dir so we don't need to re-create on every run (it takes ~2m).
	_, err := os.Stat(generateDataPath)
	if os.IsNotExist(err) {
		err = exec(
			"sh", "-c",
			fmt.Sprintf("mkdir -p %s && "+
				"docker run -i quay.io/thanos/thanosbench:v0.2.0-rc.1 block plan -p continuous-1w-small --labels 'cluster=\"eu-1\"' --labels 'replica=\"0\"' --max-time=%s | "+
				"docker run -v %s/:/shared -i quay.io/thanos/thanosbench:v0.2.0-rc.1 block gen --output.dir /shared", generateDataPath, maxTime, generateDataPath),
		)
		if err != nil {
			_ = os.RemoveAll(generateDataPath)
		}
	}
	testutil.Ok(t, err)

	// Start isolated environment with given ref.
	e, err := e2e.NewDockerEnvironment("e2e_example")
	testutil.Ok(t, err)
	// Make sure resources (e.g docker containers, network, dir) are cleaned after test.
	t.Cleanup(e.Close)

	// Start monitoring if you want to have interactive look on resources.
	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)

	// Schedule parquet tool, so we can check export produced parquet files.
	// See https://github.com/NathanHowell/parquet-tools for details.
	p := e.Runnable("parquet-tools").Init(
		e2e.StartOptions{
			Image:   "nathanhowell/parquet-tools",
			Command: e2e.NewCommandWithoutEntrypoint("tail", "-f", "/dev/null"),
		},
	)

	// Schedule StoreAPI gateway, pointing to local directory with generated dataset.
	testutil.Ok(t, exec("cp", "-r", generateDataPath+"/.", filepath.Join(e.SharedDir(), "tsdb-data")))
	store := e2edb.NewThanosStore(e, "store", []byte(fmt.Sprintf(`type: FILESYSTEM
config:
  directory: %v
`, filepath.Join(e.SharedDir(), "tsdb-data"))))

	// Run both.
	testutil.Ok(t, e2e.StartAndWaitReady(p, store))

	start := time.Now()
	// Perform export.
	{
		f, err := os.Open(filepath.Join(e.SharedDir(), "output.parquet"))
		testutil.Ok(t, err)
		defer func() {
			if f != nil {
				testutil.Ok(t, f.Close())
			}
		}()

		parsedMaxTime, err := time.Parse(time.RFC3339, maxTime)
		testutil.Ok(t, err)

		testutil.Ok(t, export1.Export5mAggregations(
			context.Background(),
			store.Endpoint("http"),
			[]*export1.LabelMatcher{{Name: "__name__", Value: "continuous_app_metric99"}},
			int64(parsedMaxTime.Add(-7*24*time.Hour).Nanosecond()*int(time.Millisecond)),
			int64(parsedMaxTime.Nanosecond()*int(time.Millisecond)),
			f,
		))

		testutil.Ok(t, f.Close())
		f = nil
	}
	fmt.Println("Export done in ", time.Since(start).String())

	// Validate if file is usable.
	stdout, stderr, err := p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "rowcount", "/shared/output.parquet"))
	testutil.Ok(t, err)
	fmt.Println(stdout, stderr)

	// Uncomment for extra interactive resources.
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}

func exec(cmd string, args ...string) error {
	if o, err := execlib.Command(cmd, args...).CombinedOutput(); err != nil {
		return errors.Wrap(err, string(o))
	}
	return nil
}
