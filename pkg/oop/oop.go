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

var _ sort.Interface = sortable{}

type sortable []Block

func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i].start.Before(s[j].start) }
func (s sortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func Compact(blocks []Block) Block {
	sort.Sort(sortable(blocks))

	g := &Group{}
	for _, b := range blocks {
		g.Add(b)
	}
	return g.Block
}
