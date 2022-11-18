package generics

import (
	"sort"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
)

type sortable []Block

func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i].Compare(s[j]) > 0 }
func (s sortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Example2() {
	toSort := []Block{ /* ... */ }
	sort.Sort(sortable(toSort))

	toSort2 := []Block{ /* ... */ }
	genericSort[Block](toSort2)

	// Output:
}

func TestSortBlocks(t *testing.T) {
	n := time.Now()
	expected := []Block{
		{start: n.Add(-10 * time.Hour)},
		{start: n},
		{start: n.Add(5 * time.Hour)},
		{start: n.Add(8 * time.Hour)},
		{start: n.Add(20 * time.Hour)},
	}
	unsorted := []Block{
		{start: n.Add(20 * time.Hour)},
		{start: n.Add(5 * time.Hour)},
		{start: n},
		{start: n.Add(8 * time.Hour)},
		{start: n.Add(-10 * time.Hour)},
	}

	t.Run("sortable", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		sort.Sort(sortable(b))

		testutil.Equals(t, expected, b)
	})
	t.Run("sort.Slice", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		sort.Slice(b, func(i, j int) bool {
			return b[i].start.Before(b[j].start)
		})

		testutil.Equals(t, expected, b)
	})
	t.Run("generics", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		genericSort[Block](b)

		testutil.Equals(t, expected, b)
	})
}

func BenchmarkSortBlock(b *testing.B) {
	n := time.Now()
	unsorted := make([]Block, 1e6)
	for i := range unsorted {
		unsorted[i].start = n.Add(-1 * time.Duration(i))
	}

	b.Run("sortable", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			sort.Sort(sortable(toSort))
		}
	})
	b.Run("sort.Slice", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			sort.Slice(toSort, func(i, j int) bool {
				return toSort[i].start.Before(toSort[j].start)
			})
		}
	})
	b.Run("generics", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			genericSort[Block](toSort)
		}
	})
}

func insertionSortInterface(data sort.Interface) {
	for i := 1; i < data.Len(); i++ {
		for j := i; j > 0 && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func insertionSortGeneric[T Comparable[T]](data []T) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && data[j].Compare(data[j-1]) > 0; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func insertionSortBlocks(data []Block) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && data[j].Compare(data[j-1]) > 0; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

type compareFunc[T any] func(a, b T) int

func insertionSortGeneric2[T any](data []T, cmpFunc compareFunc[T]) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && cmpFunc(data[j], data[j-1]) > 0; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func CompareBlocks(a, b Block) int {
	if a.start.Before(b.start) {
		return 1
	}
	if a.start.Equal(b.start) {
		return 0
	}
	return -1
}

type lessFunc[T any] func(a, b T) bool

func insertionSortGenericLess[T any](data []T, lessFunc lessFunc[T]) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && lessFunc(data[j], data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func LessBlocks(a, b Block) bool {
	return a.start.Before(b.start)
}

func insertionSortBlocksLess(data []Block, lessFunc func(a, b Block) bool) {
	for i := 1; i < len(data); i++ {
		for j := i; j > 0 && lessFunc(data[j], data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func TestInsertionSortBlocks(t *testing.T) {
	n := time.Now()
	expected := []Block{
		{start: n.Add(-10 * time.Hour)},
		{start: n},
		{start: n.Add(5 * time.Hour)},
		{start: n.Add(8 * time.Hour)},
		{start: n.Add(20 * time.Hour)},
	}
	unsorted := []Block{
		{start: n.Add(20 * time.Hour)},
		{start: n.Add(5 * time.Hour)},
		{start: n},
		{start: n.Add(8 * time.Hour)},
		{start: n.Add(-10 * time.Hour)},
	}

	t.Run("sortable", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortInterface(sortable(b))

		testutil.Equals(t, expected, b)
	})
	t.Run("generics", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortGeneric[Block](b)

		testutil.Equals(t, expected, b)
	})
	t.Run("blocks", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortBlocks(b)

		testutil.Equals(t, expected, b)
	})
	t.Run("generics_func", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortGeneric2[Block](b, CompareBlocks)

		testutil.Equals(t, expected, b)
	})
	t.Run("generics_less", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortGenericLess[Block](b, LessBlocks)

		testutil.Equals(t, expected, b)
	})
	t.Run("blocks_less", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		insertionSortBlocksLess(b, LessBlocks)

		testutil.Equals(t, expected, b)
	})
}

func BenchmarkInsertionSortBlock(b *testing.B) {
	n := time.Now()
	unsorted := make([]Block, 1e4)
	for i := range unsorted {
		unsorted[i].start = n.Add(-1 * time.Duration(i))
	}

	b.Run("sortable", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortInterface(sortable(toSort))
		}
	})
	b.Run("generics", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortGeneric[Block](toSort)
		}
	})
	b.Run("blocks", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortBlocks(toSort)
		}
	})
	b.Run("generics_func", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortGeneric2[Block](toSort, CompareBlocks)
		}
	})
	b.Run("generics_less", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortGenericLess[Block](toSort, LessBlocks)
		}
	})
	b.Run("blocks_less", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			toSort := make([]Block, len(unsorted))
			copy(toSort, unsorted)
			b.StartTimer()

			insertionSortBlocksLess(toSort, LessBlocks)
		}
	})
}
