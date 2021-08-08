package block

import (
	"fmt"
	"sort"
	"time"
)

type Block struct {
	id         string
	start, end time.Time
	// ...
}

func (b Block) String() string {
	return fmt.Sprint(b.id, ": ", b.start.Format(time.RFC3339), "-", b.end.Format(time.RFC3339))
}

type Group struct {
	Block
	// ...
}

func (g *Group) Add(b Block) {
	if g.end.IsZero() || g.end.Before(b.end) {
		g.end = b.end
	}
	if g.start.IsZero() || g.start.After(b.start) {
		g.start = b.start
	}

	if len(g.id) == 0 {
		g.id = b.id
	} else {
		g.id += "," + b.id
	}

	// ...
}

var _ sort.Interface = &Sortable{}

type Sortable []Block

func ToSortable(blocks []Block) sort.Interface {
	var s Sortable = blocks
	return &s
}

func (b *Sortable) Len() int           { return len(*b) }
func (b *Sortable) Less(i, j int) bool { return (*b)[i].start.Before((*b)[j].start) }
func (b *Sortable) Swap(i, j int)      { (*b)[i], (*b)[j] = (*b)[j], (*b)[i] }

func Compact(blocks []Block) Block {
	sort.Sort(ToSortable(blocks))

	g := &Group{}
	for _, b := range blocks {
		g.Add(b)
	}
	return g.Block
}
