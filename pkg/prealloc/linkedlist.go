// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package prealloc

// Example of pre-allocation of linked list elements.
// Read more in "Efficient Go"; Example 11-14.

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

// Delete deletes node. However, this showcases kind-of leaking code.
// Read more in "Efficient Go"; Example 11-15.
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

// ClipMemory releases unused memory.
// Read more in "Efficient Go"; Example 11-16.
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
