// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package prealloc

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/efficientgo/core/testutil"
)

const size = 1e6

func testLinkedList(t *testing.T, l *SinglyLinkedList) {
	t.Helper()

	for i := 0; i < size; i++ {
		l.Insert(i)
	}

	expected := make([]int, size)
	for i := 0; i < size; i++ {
		expected[i] = size - i - 1
	}

	got := make([]int, 0, size)
	for curr := l.head; curr != nil; curr = curr.next {
		got = append(got, curr.value)
	}
	testutil.Equals(t, expected, got)

	// Remove all but last.
	for curr := l.head; curr.next != nil; curr = curr.next {
		l.Delete(curr)
	}

	var got2 []int
	for curr := l.head; curr != nil; curr = curr.next {
		got2 = append(got2, curr.value)
	}
	testutil.Equals(t, []int{0}, got2)
}

func TestSinglyLinkedList(t *testing.T) {
	testLinkedList(t, &SinglyLinkedList{})

	p := &SinglyLinkedList{}
	p.Grow(size)
	testLinkedList(t, p)
}

func BenchmarkSinglyLinkedList(b *testing.B) {
	const size = 1e6

	b.Run("normal", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l := &SinglyLinkedList{}
			for k := 0; k < size; k++ {
				l.Insert(k)
			}
		}
	})
	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			l := &SinglyLinkedList{}
			l.Grow(size)
			for k := 0; k < size; k++ {
				l.Insert(k)
			}
		}
	})
}

func BenchmarkSinglyLinkedList_Delete(b *testing.B) {
	const size = 1e6

	b.Run("normal", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			l := &SinglyLinkedList{}
			for k := 0; k < size; k++ {
				l.Insert(k)
			}
			b.StartTimer()

			// Remove all but last.
			for curr := l.head; curr.next != nil; curr = curr.next {
				l.Delete(curr)
			}
		}
	})
	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			l := &SinglyLinkedList{}
			l.Grow(size)
			for k := 0; k < size; k++ {
				l.Insert(k)
			}
			l.pool = nil // Dispose pool just in case.
			b.StartTimer()

			// Remove all but last.
			for curr := l.head; curr.next != nil; curr = curr.next {
				l.Delete(curr)
			}
		}
	})
}

func _printHeapUsage(prefix string) {
	m := runtime.MemStats{}

	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Println(prefix, float64(m.HeapAlloc)/1024.0, "KB")
}

func TestSinglyLinkedList_Delete1(t *testing.T) {
	l := &SinglyLinkedList{}
	for k := 0; k < size; k++ {
		l.Insert(k)
	}
	_printHeapUsage("Heap before deletions:        ")

	// Remove all but last.
	for curr := l.head; curr.next != nil; curr = curr.next {
		l.Delete(curr)
	}
	_printHeapUsage("Heap after deleting all - 1:  ")

	l.Delete(l.head)
	_printHeapUsage("Heap after last was deleted:   ")
}

func TestSinglyLinkedList_Delete2(t *testing.T) {
	l := &SinglyLinkedList{}
	l.Grow(size)
	for k := 0; k < size; k++ {
		l.Insert(k)
	}
	l.pool = nil // Dispose pool.
	_printHeapUsage("Heap before deletions:        ")

	// Remove all but last.
	for curr := l.head; curr.next != nil; curr = curr.next {
		l.Delete(curr)
	}
	_printHeapUsage("Heap after deleting all - 1:  ")

	l.Delete(l.head)
	_printHeapUsage("Heap after last was deleted:  ")
}

func TestSinglyLinkedList_Delete3(t *testing.T) {
	l := &SinglyLinkedList{}
	l.Grow(size)
	for k := 0; k < size; k++ {
		l.Insert(k)
	}
	l.pool = nil // Dispose pool.
	_printHeapUsage("Heap before deletions:        ")

	// Remove all but last.
	for curr := l.head; curr.next != nil; curr = curr.next {
		l.Delete(curr)
	}
	_printHeapUsage("Heap after deleting all - 1:  ")

	l.ClipMemory()

	_printHeapUsage("Heap after clipping:          ")

	l.Delete(l.head)
	_printHeapUsage("Heap after last was deleted:  ")
}
