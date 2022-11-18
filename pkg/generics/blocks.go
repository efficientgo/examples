package generics

import (
	"sort"
	"time"

	"github.com/google/uuid"
)

type Block struct {
	id         uuid.UUID
	start, end time.Time
	// ...
}

// Example of generic sorting function for any type that supports `Compare` method.
// Read more in "Efficient Go"; Example 2-15.

type Comparable[T any] interface {
	Compare(T) int
}

type genericSortable[T Comparable[T]] []T

func (s genericSortable[T]) Len() int           { return len(s) }
func (s genericSortable[T]) Less(i, j int) bool { return s[i].Compare(s[j]) > 0 }
func (s genericSortable[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func genericSort[T Comparable[T]](slice []T) {
	sort.Sort(genericSortable[T](slice))
}

// Compare is anti-pattern: For type parameters, prefer functions to methods.
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
