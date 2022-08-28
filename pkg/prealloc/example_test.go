package prealloc_test

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/examples/pkg/prealloc"
)

func TestPrealloc(t *testing.T) {
	n := int(1e4)
	wp := prealloc.CreateSlice(n)
	testutil.Equals(t, n*7, len(wp))
	testutil.Equals(t, n*7+5776, cap(wp))

	p := prealloc.CreateSlice2(n)
	testutil.Equals(t, wp, p)
	testutil.Equals(t, n*7, len(p))
	testutil.Equals(t, len(p), cap(p))
}

func optimizedPrealloc(b testutil.TB) {
	const size = 1e6

	slice := make([]string, 0, size)
	for i := 0; i < size; i++ {
		slice = append(slice, "something")
	}

	slice2 := make([]string, size)
	for i := 0; i < size; i++ {
		slice2[i] = "something"
	}

	m := make(map[int]string, size)
	for i := 0; i < size; i++ {
		m[i] = "something"
	}

	buf := bytes.Buffer{}
	buf.Grow(size)
	for i := 0; i < size; i++ {
		_ = buf.WriteByte('a')
	}

	builder := strings.Builder{}
	builder.Grow(size)
	for i := 0; i < size; i++ {
		builder.WriteByte('a')
	}

	buf2, _ := ReadAll1(bytes.NewReader(make([]byte, size)), size)
	buf3, _ := ReadAll2(bytes.NewReader(make([]byte, size)), size)

	// .Test
	if !b.IsBenchmark() {
		testutil.Equals(b, slice, slice2)
		testutil.Equals(b, buf.String(), builder.String())
		testutil.Equals(b, buf2, buf3)
		testutil.Equals(b, buf2, make([]byte, size))
	}
}

func unOptimizedPrealloc(b testutil.TB) {
	const size = 1e6

	var slice []string
	for i := 0; i < size; i++ {
		slice = append(slice, "something")
	}

	var slice2 []string
	for i := 0; i < size; i++ {
		slice2 = append(slice2, "something")
	}

	m := make(map[int]string)
	for i := 0; i < size; i++ {
		m[i] = "something"
	}

	buf := bytes.Buffer{}
	for i := 0; i < size; i++ {
		_ = buf.WriteByte('a')
	}

	builder := strings.Builder{}
	for i := 0; i < size; i++ {
		builder.WriteByte('a')
	}

	buf2, _ := io.ReadAll(bytes.NewReader(make([]byte, size)))
	buf3, _ := io.ReadAll(bytes.NewReader(make([]byte, size)))

	// .Test
	if !b.IsBenchmark() {
		testutil.Equals(b, slice, slice2)
		testutil.Equals(b, buf.String(), builder.String())
		testutil.Equals(b, buf2, buf3)
		testutil.Equals(b, buf2, make([]byte, size))
	}
}

func ReadAll1(r io.Reader, size int) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Grow(size)
	n, err := io.Copy(&buf, r)
	return buf.Bytes()[:n], err
}

func ReadAll2(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err == io.EOF {
		err = nil
	}
	return buf[:n], err
}

func TestPreallocs(t *testing.T) {
	unOptimizedPrealloc(testutil.NewTB(t))
	optimizedPrealloc(testutil.NewTB(t))
}

func BenchmarkReadAlls(b *testing.B) {
	const size = int(1e6)
	inner := make([]byte, size)
	b.Run("io.ReadAll", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := io.ReadAll(bytes.NewReader(inner))
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
	b.Run("ReadAll1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := ReadAll1(bytes.NewReader(inner), size)
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
	b.Run("ReadAll2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := ReadAll2(bytes.NewReader(inner), size)
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
}

type Node struct {
	next  *Node
	value int
}

type SinglyLinkedList struct {
	head *Node

	pool      []Node
	poolIndex int
}

func (l *SinglyLinkedList) Grow(len int) {
	l.pool = make([]Node, len)
	l.poolIndex = 0
}

func (l *SinglyLinkedList) Insert(value int) {
	var newNode *Node
	if len(l.pool) > l.poolIndex {
		newNode = &l.pool[l.poolIndex]
		l.poolIndex++
	} else {
		newNode = &Node{}
	}

	newNode.next = l.head
	newNode.value = value
	l.head = newNode
}

// NOTE: Showcase of kind-of leaking code.
func (l *SinglyLinkedList) Delete(n *Node) {
	if l.head == n {
		l.head = n.next
		return
	}

	for curr := l.head; curr != nil; curr = curr.next {
		if curr.next != n {
			continue
		}

		curr.next = n.next
		return
	}
}

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

func (l *SinglyLinkedList) ClipMemory() {
	var objs int
	for curr := l.head; curr != nil; curr = curr.next {
		objs++
	}

	l.pool = make([]Node, objs)
	l.poolIndex = 0
	for curr := l.head; curr != nil; curr = curr.next {
		oldCurr := curr
		curr = &l.pool[l.poolIndex]
		l.poolIndex++

		curr.next = oldCurr.next
		curr.value = oldCurr.value

		if oldCurr == l.head {
			l.head = curr
		}
	}
}
