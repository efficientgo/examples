// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package leak

import (
	"net/http"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
)

// Examples of code which is leaking resources, because of not exhausted readers.
// Read more in "Efficient Go"; Example 11-10.

func handleResp_Wrong(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

func handleResp_StillWrong(resp *http.Response) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

func handleResp_Better(resp *http.Response) (err error) {
	defer errcapture.ExhaustClose(&err, resp.Body, "close")
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}
