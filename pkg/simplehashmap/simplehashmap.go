package simplehashmap

import (
	"fmt"
)

// Adapted from https://github.com/superhawk610/dumbhashmap/blob/master/dumbhashmap/dumbhashmap.go

const bucketCnt = 256

type bucket struct {
	entries [][]interface{}
}

type SimpleHashMap struct {
	buckets []bucket
}

func New(cap int) (h *SimpleHashMap) {
	s := &SimpleHashMap{make([]bucket, bucketCnt)}
	for i := range s.buckets {
		s.buckets[i].entries = make([][]interface{}, 0, cap/bucketCnt)
	}
	return s
}

func hash(key uint64) uint32 {
	return uint32(key) % bucketCnt
}

func (h SimpleHashMap) Get(key uint64) interface{} {
	b := h.buckets[hash(key)]
	if len(b.entries) == 0 {
		return nil
	}

	var value interface{}
	for _, e := range b.entries {
		if e[0] == key {
			value = e[1]
		}
	}
	return value
}

func (h SimpleHashMap) Set(key uint64, value interface{}) {
	b := &h.buckets[hash(key)]
	b.entries = append(b.entries, []interface{}{key, value})
}

func (h SimpleHashMap) Delete(key uint64) (ok bool) {
	b := &h.buckets[hash(key)]
	for i, e := range b.entries {
		if e[0] == key {
			ok = true
			b.entries = append(b.entries[:i], b.entries[i+1:]...)
		}
	}
	return ok
}

func (h SimpleHashMap) String() string {
	var strValue = ""
	for i, b := range h.buckets {
		strValue += fmt.Sprintf("%v{", i)
		for j, e := range b.entries {
			strValue += fmt.Sprintf("\n  %v[%v, %v]\n", j, e[0], e[1])
		}
		strValue += "}\n"
	}

	if strValue == "" {
		return "<empty SimpleHashMap>"
	}
	return strValue
}
