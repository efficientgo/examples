package generics

import (
	"sort"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type Block struct {
	id         uuid.UUID
	start, end time.Time
	// ...
}

func (b Block) Duration() time.Duration {
	return b.end.Sub(b.start)
}

type Group struct {
	Block

	children []uuid.UUID
	// ...
}

func (g *Group) Merge(b Block) {
	if g.end.IsZero() || g.end.Before(b.end) {
		g.end = b.end
	}
	if g.start.IsZero() || g.start.After(b.start) {
		g.start = b.start
	}

	g.children = append(g.children, b.id)

	// ...
}

func Compact(blocks ...Block) Block {
	sort.Sort(sortable(blocks))

	g := &Group{}
	g.id = uuid.New()
	for _, b := range blocks {
		g.Merge(b)
	}
	return g.Block
}

type Comparable[T any] interface {
	Compare(T) int
}

type sortable []Block

func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i].Compare(s[j]) > 0 }
func (s sortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Compare is anti-pattern: For type parameters, prefer functions to methods
// https://go.dev/blog/when-generics
func (b Block) Compare(other Block) int {
	if b.start.Before(other.start) {
		return 1
	}
	if b.start.Equal(other.start) {
		return 0
	}
	return -1
}

type genericSortable[T Comparable[T]] []T

func (s genericSortable[T]) Len() int           { return len(s) }
func (s genericSortable[T]) Less(i, j int) bool { return s[i].Compare(s[j]) > 0 }
func (s genericSortable[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genericSort[T Comparable[T]](slice []T) {
	sort.Sort(genericSortable[T](slice))
}

type genericSortableBasic[T constraints.Ordered] []T

func (s genericSortableBasic[T]) Len() int           { return len(s) }
func (s genericSortableBasic[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s genericSortableBasic[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genericSortBasic[T constraints.Ordered](slice []T) {
	sort.Sort(genericSortableBasic[T](slice))
}

func Example() {
	toSort := []int{-20, 1, 10, 20}
	sort.Ints(toSort)

	toSort2 := []int{-20, 1, 10, 20}
	genericSortBasic[int](toSort2)

	// Output:
}

func Example2() {
	toSort := []Block{ /* ... */ }
	sort.Sort(sortable(toSort))

	toSort2 := []Block{ /* ... */ }
	genericSort[Block](toSort2)

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
