package leak

import (
	"net/http"
	"testing"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
	"go.uber.org/goleak"
)

func handleResp_Wrong(resp *http.Response) error {
	// Not reading and not closing the response body...

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

func handleResp_StillWrong(resp *http.Response) error {
	defer func() {
		// https://pkg.go.dev/net/http#Client.Do
		// If the returned error is nil, the Response will contain a non-nil
		// Body which the user is expected to close. If the Body is not both
		// read to EOF and closed, the Client's underlying RoundTripper
		// (typically Transport) may not be able to re-use a persistent TCP
		// connection to the server for a subsequent "keep-alive" request.
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

func handleResp_Better(resp *http.Response) (err error) {
	defer errcapture.ExhaustClose(&err, resp.Body, "close")

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

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
		testutil.Ok(b, handleResp_StillWrong(resp))
	}
}
