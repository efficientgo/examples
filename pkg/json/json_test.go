package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func generateTestData(t testing.TB) []Item {
	t.Helper()

	items := make([]Item, 10e3)
	for i := 0; i < len(items); i += 5 {
		c := i
		items[c] = Item{
			ID:           c,
			Name:         fmt.Sprintf("T-shirt ABC (%d)", i),
			Price:        20,
			Size:         Size{Height: 54, Width: 49, Depth: 10},
			PackagedSize: Size{Height: 60, Width: 50, Depth: 10},
			Weight:       500,
			Available:    10,
			Sold:         1231344,
		}
		c++
		items[c] = Item{
			ID:           c,
			Name:         fmt.Sprintf("Hoodie (%d)", i),
			Price:        82,
			Size:         Size{Height: 64, Width: 55, Depth: 15},
			PackagedSize: Size{Height: 70, Width: 60, Depth: 20},
			Weight:       700,
			Available:    244,
			Sold:         21311,
		}
		c++
		items[c] = Item{
			ID:           c,
			Name:         fmt.Sprintf("Mug (%d)", i),
			Price:        15,
			Size:         Size{Height: 10, Width: 8, Depth: 8},
			PackagedSize: Size{Height: 12, Width: 10, Depth: 10},
			Weight:       300,
			Available:    1234,
			Sold:         30,
		}
		c++
		items[c] = Item{
			ID:           c,
			Name:         fmt.Sprintf("Water Bottle (%d)", i),
			Price:        10,
			Size:         Size{Height: 25, Width: 7, Depth: 7},
			PackagedSize: Size{Height: 27, Width: 9, Depth: 9},
			Weight:       200,
			Available:    33214,
			Sold:         121110,
		}
		c++
		items[c] = Item{
			ID:           c,
			Name:         fmt.Sprintf("Phone Case (%d)", i),
			Price:        18,
			Size:         Size{Height: 15, Width: 8, Depth: 3},
			PackagedSize: Size{Height: 17, Width: 9, Depth: 5},
			Weight:       100,
			Available:    12,
			Sold:         23445,
		}
	}
	return items
}

func TestSell(t *testing.T) {
	for _, tc := range []func(string) error{
		sell0,
		sell,
		sell2,
		sell3,
		sell4,
	} {
		t.Run("", func(t *testing.T) {
			inputFile := t.TempDir() + "/db_10k.json"

			items := generateTestData(t)
			av, sold := items[39].Available, items[39].Sold

			o, err := os.Create(inputFile)
			testutil.Ok(t, err)
			defer func() { _ = o.Close() }()
			testutil.Ok(t, json.NewEncoder(o).Encode(items))

			testutil.Ok(t, tc(inputFile))

			f, err := os.Open(inputFile + ".updated")
			testutil.Ok(t, err)
			defer func() { _ = f.Close() }()

			items = items[:0]
			testutil.Ok(t, json.NewDecoder(f).Decode(&items))

			testutil.Equals(t, av-1, items[39].Available)
			testutil.Equals(t, sold+1, items[39].Sold)
		})
	}
}

func BenchmarkSell(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testutil.Ok(b, sell0("./db_10k.json"))
	}
}

func BenchmarkMarshal(b *testing.B) {
	items := generateTestData(b)
	buf := bytes.Buffer{}
	testutil.Ok(b, json.NewEncoder(&buf).Encode(items))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items = items[:0]
		testutil.Ok(b, json.Unmarshal(buf.Bytes(), &items))
	}
}

func BenchmarkWrite(b *testing.B) {
	out := make([]byte, 1024*1024*20)
	for i := range out {
		out[i] = byte(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o, err := os.Create("out.txt") // 4.5ms / 8.2ms fsync
		testutil.Ok(b, err)

		_, err = o.Write(out)
		testutil.Ok(b, err)

		//testutil.Ok(b, o.Sync()) // fsync.

		testutil.Ok(b, o.Close())
	}
	// 2560 frames 4.5ms -> ~1.7μs
}

func BenchmarkRead(b *testing.B) {
	out := make([]byte, 1024*1024*20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.Open("out.txt") // 1.7 ms read ~ 0.6 μs
		testutil.Ok(b, err)

		_, err = f.Read(out)
		testutil.Ok(b, err)

		f.Close()
	}
}
