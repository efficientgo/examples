package leak

import (
	"net/http"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
	"go.uber.org/goleak"
)

var c = &http.Client{}

func BenchmarkClient(b *testing.B) {
	defer goleak.VerifyNone(b)

	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://google.com")
		testutil.Ok(b, err)
		testutil.Equals(b, http.StatusOK, resp.StatusCode)

		// Not reading and not closing the response body...
	}
}