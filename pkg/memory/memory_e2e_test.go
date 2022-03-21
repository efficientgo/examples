package memory

import (
	"context"
	"testing"

	"github.com/efficientgo/e2e"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
)

// export var=v1 && go test -count 1 -benchtime 1x -run '^$' -bench BenchmarkAllocButNotAccess -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt
func BenchmarkAllocButNotAccess(b *testing.B) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = AllocButNotAccess(cancelled)
	}
}

func TestMemE2e(t *testing.T) {
	e, err := e2e.NewDockerEnvironment("e2e", e2e.WithVerbose())
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	s, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)

	l, err := e2e.Containerize(e, "run", AllocButNotAccess)
	testutil.Ok(t, err)

	testutil.Ok(t, e2e.StartAndWaitReady(l))

	testutil.Ok(t, s.OpenUserInterfaceInBrowser(`/graph?g0.expr=container_memory_rss&g0.tab=0&g0.stacked=0&g0.range_input=15m&g1.expr=container_memory_working_set_bytes&g1.tab=0&g1.stacked=0&g1.range_input=15m&g2.expr=go_memstats_heap_alloc_bytes&g2.tab=0&g2.stacked=0&g2.range_input=15m&g3.expr=go_memstats_heap_idle_bytes&g3.tab=0&g3.stacked=0&g3.range_input=15m`))
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}