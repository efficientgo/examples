// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/runutil"
	"github.com/efficientgo/core/testutil"
	"github.com/go-kit/log"
)

func TestGetProfile(t *testing.T) {
	tmpDir := t.TempDir()
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		<-errCh
	})

	go func() {
		errCh <- runMain(ctx, []string{`-objstore.config=type: FILESYSTEM
config:
  directory: "."`})
		close(errCh)
	}()

	rctx, rcancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer rcancel()
	testutil.Ok(t, runutil.RetryWithLog(log.NewLogfmtLogger(os.Stderr), 1*time.Second, rctx.Done(), func() error {
		res, err := http.Get("http://localhost:8080/debug/fgprof/profile?seconds=1")
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(tmpDir, "fgprof"))
		if err != nil {
			return err
		}

		if _, err := io.Copy(f, res.Body); err != nil {
			_ = f.Close()
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		fmt.Println(res.Status)

		if res.StatusCode != http.StatusOK {
			return errors.Newf("expected OK, got %v", res.StatusCode)
		}
		return nil
	}))

	select {
	case err := <-errCh:
		testutil.Ok(t, err)
		t.Fatal("expected to not fail")
	default:
	}
}
