package generics

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Example of generic sorting function for slice with any type that can be compared.
// Read more in "Efficient Go"; Example 2-14.

type genericSortableBasic[T constraints.Ordered] []T

func (s genericSortableBasic[T]) Len() int           { return len(s) }
func (s genericSortableBasic[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s genericSortableBasic[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genericSortBasic[T constraints.Ordered](slice []T) {
	sort.Sort(genericSortableBasic[T](slice))
}
