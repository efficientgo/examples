// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package leak

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Example case of leak in HTTP handlers.
// Read more in "Efficient Go"; Example 11-2.

func ComplexComputation() int {
	time.Sleep(1 * time.Second) // Computation.
	time.Sleep(1 * time.Second) // Cleanup.
	return 4
}

func Handle_VeryWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

// More examples of leaks in HTTP handlers.
// Read more in "Efficient Go"; Example 11-5.

func Handle_Wrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

func ComplexComputationWithCtx(ctx context.Context) (ret int) {
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second): // Computation.
		ret = 4
	}

	time.Sleep(1 * time.Second) // Cleanup.
	return ret
}

func Handle_AlsoWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context())
	}()

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

// Recommended code that does not leak.
// Read more in "Efficient Go"; Example 11-6.

func Handle_Better(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context())
	}()

	resp := <-respCh
	if r.Context().Err() != nil {
		return
	}

	_, _ = w.Write([]byte(strconv.Itoa(resp)))
}
