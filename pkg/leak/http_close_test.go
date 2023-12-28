// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package leak

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"go.uber.org/goleak"
)

// TestHandleCancel tests against leaks when cancelling the request.
// Read more in "Efficient Go"; Example 11-3.
func TestHandleCancel(t *testing.T) {
	defer goleak.VerifyNone(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "https://efficientgo.com", nil)

	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		Handle_Better(w, r.WithContext(ctx))
		wg.Done()
	}()
	// Immediately cancel.
	cancel()

	time.Sleep(3 * time.Second)
	wg.Wait()
}

func TestHandle(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Handle_VeryWrong", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_VeryWrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("Handle_Wrong", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_Wrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("Handle_AlsoWrong", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_AlsoWrong(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})

	t.Run("Handle_Better", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("", "https://efficientgo.com", nil)
		Handle_Better(w, r)

		testutil.Equals(t, http.StatusOK, w.Code)
		testutil.Equals(t, "4", w.Body.String())
	})
}

// Example of leaks in benchmarks (also wrong results).
// Read more in "Efficient Go"; Example 11-7.

func BenchmarkComplexComputation_Wrong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go func() { ComplexComputation() }()
		go func() { ComplexComputation() }()
	}
}

func BenchmarkComplexComputation_Better(b *testing.B) {
	defer goleak.VerifyNone(
		b,
		goleak.IgnoreTopFunction("testing.(*B).run1"),
		goleak.IgnoreTopFunction("testing.(*B).doBench"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			ComplexComputation()
		}()
		go func() {
			defer wg.Done()
			ComplexComputation()
		}()

		wg.Wait()
	}
}
