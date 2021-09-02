package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	execlib "os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/parquet-export/export1"
	"github.com/efficientgo/examples/pkg/parquet-export/export2"
	"github.com/efficientgo/examples/pkg/parquet-export/ref"
	"github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/efficientgo/tools/performance/pkg/profiles"
	"github.com/pkg/errors"
	"github.com/xitongsys/parquet-go-source/buffer"
	"github.com/xitongsys/parquet-go/reader"
	"google.golang.org/grpc"
)

var (
	generateDataPath = func() string { a, _ := filepath.Abs("generated"); return a }()
	maxTime          = `2021-07-20T00:00:00Z`
)

type exportFuncType func(ctx context.Context, address string, metricSelector []*ref.LabelMatcher, minTime, maxTime int64, w io.Writer) (seriesNum int, samplesNum int, err error)

type mockSeries struct {
	series []*export1.Series
}

func (m *mockSeries) Series(_ *export1.SeriesRequest, ss export1.Store_SeriesServer) error {
	for _, s := range m.series {
		if err := ss.Send(&export1.SeriesResponse{Result: &export1.SeriesResponse_Series{Series: s}}); err != nil {
			return err
		}
	}
	return nil
}

type sample struct {
	v float64
	t int64
}

func chunkFromSamples(t *testing.T, samples []sample) *export1.Chunk {
	chk := chunkenc.NewXORChunk()
	a, err := chk.Appender()
	testutil.Ok(t, err)

	for _, s := range samples {
		a.Append(s.t, s.v)
	}

	return &export1.Chunk{
		Type: export1.Chunk_Encoding(chk.Encoding() - 1),
		Data: chk.Bytes(),
	}
}

// TestParquetExport tests all export implementations to ensure functional correctness after and before optimizations.
func TestParquetExport(t *testing.T) {
	s := &mockSeries{
		series: []*export1.Series{
			{
				Labels: []*export1.Label{{Name: "__blockgen_target__", Value: "1"}, {Name: "__name__", Value: "continuous_app_metric0"}, {Name: "cluster", Value: "eu-1"}, {Name: "replica", Value: "0"}},
				Chunks: []*export1.AggrChunk{
					{Raw: chunkFromSamples(t, []sample{{v: -1, t: 0}, {v: 10, t: 2 * 60 * 1000}, {v: -10, t: 5 * 60 * 1000}, {v: 20, t: 7 * 60 * 1000}, {v: -20, t: 10 * 60 * 1000}, {v: 15, t: 15 * 60 * 1000}})},
					{Raw: chunkFromSamples(t, []sample{{v: -1, t: 17 * 60 * 1000}, {v: 10, t: 20 * 60 * 1000}})},
				},
			},
			{
				Labels: []*export1.Label{{Name: "__blockgen_target__", Value: "2"}, {Name: "__name__", Value: "continuous_app_metric0"}, {Name: "cluster", Value: "eu-1"}, {Name: "replica", Value: "0"}},
				Chunks: []*export1.AggrChunk{
					{Raw: chunkFromSamples(t, []sample{{v: -99, t: 0}, {v: 10, t: 2 * 60 * 1000}, {v: -10, t: 5 * 60 * 1000}, {v: 20, t: 7 * 60 * 1000}, {v: -20, t: 10 * 60 * 1000}, {v: 15, t: 15 * 60 * 1000}})},
					{Raw: chunkFromSamples(t, []sample{{v: -99, t: 17 * 60 * 1000}, {v: 10, t: 20 * 60 * 1000}})},
				},
			},
		},
	}

	srv := grpc.NewServer()
	export1.RegisterStoreServer(srv, s)

	l, err := net.Listen("tcp", "localhost:0")
	testutil.Ok(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = srv.Serve(l)
		wg.Done()
	}()

	for _, exportFn := range []exportFuncType{export1.Export5mAggregations, export2.Export5mAggregations} {
		t.Run("", func(t *testing.T) {
			b := bytes.Buffer{}
			seriesNum, samplesNum, err := exportFn(context.Background(), l.Addr().String(), []*ref.LabelMatcher{{Name: "__name__", Value: "", Type: ref.LabelMatcher_NEQ}}, 0, 1, &b)
			testutil.Ok(t, err)

			testutil.Equals(t, 2, seriesNum)
			testutil.Equals(t, 16, samplesNum)

			bb, err := buffer.NewBufferFile(b.Bytes())
			testutil.Ok(t, err)

			r, err := reader.NewParquetReader(bb, new(ref.Aggregation), 1)
			testutil.Ok(t, err)

			var aggr []ref.Aggregation
			for i := 0; i < int(r.GetNumRows()); i++ {
				a := make([]ref.Aggregation, 1)
				testutil.Ok(t, r.Read(&a))
				aggr = append(aggr, a[0])
			}
			r.ReadStop()

			testutil.Equals(t, []ref.Aggregation{
				{NameLabel: "continuous_app_metric0", TargetLabel: "1", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 300000, Count: 3, Sum: -1, Min: -10, Max: 10},
				{NameLabel: "continuous_app_metric0", TargetLabel: "1", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 720000, Count: 2, Sum: 0, Min: -20, Max: 20},
				{NameLabel: "continuous_app_metric0", TargetLabel: "1", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 1200000, Count: 3, Sum: 24, Min: -1, Max: 15},
				{NameLabel: "continuous_app_metric0", TargetLabel: "2", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 300000, Count: 3, Sum: -99, Min: -99, Max: 10},
				{NameLabel: "continuous_app_metric0", TargetLabel: "2", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 720000, Count: 2, Sum: 0, Min: -20, Max: 20},
				{NameLabel: "continuous_app_metric0", TargetLabel: "2", ClusterLabel: "eu-1", ReplicaLabel: "0", Timestamp: 1200000, Count: 3, Sum: -74, Min: -99, Max: 15},
			}, aggr)
		})
	}

	srv.Stop()
	wg.Wait()
}

// Testing export1 for now. Change it to other packages for better performance.
var exportFunction exportFuncType = export1.Export5mAggregations

// Recommended test args: -test.timeout 9999m for interactive mode experience.
func TestParquetExportIntegration(t *testing.T) {
	t.Parallel()

	testParquetExportIntegration(testutil.NewTB(t))
}

// Recommended test args: -test.benchmem -test.benchtime=1m
func BenchmarkParquetExportIntegration(b *testing.B) {
	testParquetExportIntegration(testutil.NewTB(b))
}

func testParquetExportIntegration(tb testutil.TB) {
	// Pinning to one CPU only for deterministic latency with 1 CPU.
	runtime.GOMAXPROCS(1)

	ctx := context.Background()

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
	testutil.Ok(tb, err)

	// Start isolated environment with given reference.
	e, err := e2e.NewDockerEnvironment("parquet_bench")
	testutil.Ok(tb, err)
	// Make sure resources (e.g docker containers, network, dir) are cleaned after test.
	tb.Cleanup(e.Close)

	var mon *e2emonitoring.Service
	var p e2e.Runnable
	if !tb.IsBenchmark() {
		// Start monitoring if you want to have interactive look on resources.
		mon, err = e2emonitoring.Start(e, e2emonitoring.WithCurrentProcessAsContainer())
		testutil.Ok(tb, err)

		// Schedule parquet tool, so we can check export produced parquet files.
		// See https://github.com/NathanHowell/parquet-tools for details.
		p = e.Runnable("parquet-tools").Init(
			e2e.StartOptions{
				Image:   "nathanhowell/parquet-tools",
				Command: e2e.NewCommandWithoutEntrypoint("tail", "-f", "/dev/null"),
			},
		)
		testutil.Ok(tb, e2e.StartAndWaitReady(p))
	}

	// Schedule StoreAPI gateway, pointing to local directory with generated dataset.
	testutil.Ok(tb, exec("cp", "-r", generateDataPath+"/.", filepath.Join(e.SharedDir(), "tsdb-data")))
	store := e2edb.NewThanosStore(e, "store", []byte(`type: FILESYSTEM
config:
  directory: "/shared/tsdb-data"
`))
	testutil.Ok(tb, e2e.StartAndWaitReady(store))

	parsedMaxTime, err := time.Parse(time.RFC3339, maxTime)
	testutil.Ok(tb, err)

	minTime := ref.TimestampFromTime(parsedMaxTime.Add(-7 * 24 * time.Hour))
	maxTime := ref.TimestampFromTime(parsedMaxTime)

	c, err := profiles.StartCPU(".", profiles.CPUTypeFGProf)
	testutil.Ok(tb, err)

	for _, tcase := range []struct {
		matchers []*ref.LabelMatcher
	}{
		{matchers: []*ref.LabelMatcher{{Name: "__name__", Value: "continuous_app_metric9.{1}", Type: ref.LabelMatcher_RE}}}, // 1k series.
		//{matchers: []*export1.LabelMatcher{{Name: "__name__", Value: "", Type: export1.LabelMatcher_NEQ}}}, // All, 10k series.
	} {
		tb.Run(fmt.Sprintf("%v", tcase.matchers), func(tb testutil.TB) {
			tb.ResetTimer()

			// Perform export.
			for i := 0; i < tb.N(); i++ {
				start := time.Now()

				f, err := os.OpenFile(filepath.Join(e.SharedDir(), "output.parquet"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
				testutil.Ok(tb, err)
				defer func() {
					if f != nil {
						testutil.Ok(tb, f.Close())
					}
				}()

				seriesNum, samplesNum, err := exportFunction(ctx, store.Endpoint("grpc"), tcase.matchers, minTime, maxTime, f)
				testutil.Ok(tb, err)
				testutil.Ok(tb, f.Close())
				f = nil

				if !tb.IsBenchmark() {
					fmt.Println("Export done in ", time.Since(start).String(), "exported", seriesNum, "series,", samplesNum, "samples")

					// TODO(bwplotka): Assert on it.
					// Validate if file is usable, by parquet tooling.
					stdout, stderr, err := p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "rowcount", "-d", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)

					stdout, stderr, err = p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "size", "-d", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)

					// Print 5 records.
					stdout, stderr, err = p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "head", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)
				}
			}
		})
	}

	testutil.Ok(tb, c())
	runtime.GC()
	testutil.Ok(tb, profiles.Heap("."))

	if !tb.IsBenchmark() {
		// Uncomment for extra interactive resources.
		testutil.Ok(tb, mon.OpenUserInterfaceInBrowser())
		testutil.Ok(tb, e2einteractive.RunUntilEndpointHit())
	}
}

func exec(cmd string, args ...string) error {
	if o, err := execlib.Command(cmd, args...).CombinedOutput(); err != nil {
		return errors.Wrap(err, string(o))
	}
	return nil
}
