package leak

import (
	"net/http"
	"testing"

	"github.com/efficientgo/core/testutil"
	"go.uber.org/goleak"
)

func BenchmarkClient(b *testing.B) {
	defer goleak.VerifyNone(
		b,
		goleak.IgnoreTopFunction("testing.(*B).run1"),
		goleak.IgnoreTopFunction("testing.(*B).doBench"),
	)
	c := &http.Client{}
	defer c.CloseIdleConnections()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := c.Get("http://google.com")
		testutil.Ok(b, err)
		testutil.Ok(b, handleResp_Wrong(resp))
	}
}
