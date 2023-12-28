// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package generics

import (
	"sort"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func ExampleBasic() {
	toSort := []int{-20, 1, 10, 20}
	sort.Ints(toSort)

	toSort2 := []int{-20, 1, 10, 20}
	genericSortBasic[int](toSort2)

	// Output:
}

func TestSortSimple(t *testing.T) {
	expected := []int{-20, 1, 10, 20}
	unsorted := []int{10, 20, -20, 1}

	t.Run("sortable", func(t *testing.T) {
		toSort := make([]int, len(unsorted))
		copy(toSort, unsorted)

		sort.Ints(toSort)

		testutil.Equals(t, expected, toSort)
	})
	t.Run("generics", func(t *testing.T) {
		toSort := make([]int, len(unsorted))
		copy(toSort, unsorted)

		genericSortBasic[int](toSort)

		testutil.Equals(t, expected, toSort)
	})
}
