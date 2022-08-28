package generics

import (
	"sort"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
)

type Block struct {
	id         string
	start, end time.Time
	// ...
}

type sortable []Block

func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i].start.Before(s[j].start) }
func (s sortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type Comparable[T any] interface {
	Compare(T) int
}

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

func partitionGeneric[T Comparable[T]](arr []T, low, high int) int {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[i].Compare(pivot) > 0 {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return i
}

func quickSortGeneric[T Comparable[T]](arr []T, low, high int) {
	if low < high {
		p := partitionGeneric[T](arr, low, high)
		quickSortGeneric[T](arr, low, p-1)
		quickSortGeneric[T](arr, p+1, high)
	}
}

func TestQuickSortBlocks(t *testing.T) {
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

		quickSortInterface(sortable(b), 0, len(b)-1)

		testutil.Equals(t, expected, b)
	})
	t.Run("generics", func(t *testing.T) {
		b := make([]Block, len(unsorted))
		copy(b, unsorted)

		quickSortGeneric[Block](b, 0, len(b)-1)

		testutil.Equals(t, expected, b)
	})
}

func BenchmarkQuickSortBlock(b *testing.B) {
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

			quickSortInterface(sortable(toSort), 0, len(toSort)-1)
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

			quickSortGeneric[Block](toSort, 0, len(toSort)-1)
		}
	})
}
