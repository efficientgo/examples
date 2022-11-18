package pools

import "sync"

// Example of simple buffering and sync.Pool as a solution.
// Read more in "Efficient Go"; Example 11-17, 11-18, 11-19.

func processUsingBuffer(buf []byte) {
	buf = buf[:0]

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}

	// Use buffer...
}

func processUsingPool_Wrong(p *sync.Pool) {
	buf := p.Get().([]byte)
	buf = buf[:0]

	defer p.Put(buf)

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}

	// Use buffer...
}

func processUsingPool(p *sync.Pool) {
	buf := p.Get().([]byte)
	buf = buf[:0]

	for i := 0; i < 1e6; i++ {
		buf = append(buf, 'a')
	}
	defer p.Put(buf)

	// Use buffer...
}
