package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func generateTestData(t testing.TB, count int) []Item {
	t.Helper()

	items := make([]Item, count)
	for i := 0; i < len(items); i += 5 {
		c := i
		items[c] = Item{
			ID:     c,
			Name:   fmt.Sprintf("T-shirt ABC (%d)", i), // 18.8
			Size:   [3]int{54, 49, 9},
			Weight: 500,
		}
		c++
		items[c] = Item{
			ID:     c,
			Name:   fmt.Sprintf("Hoodie (%d)", i),
			Size:   [3]int{64, 55, 15},
			Weight: 700,
		}
		c++
		items[c] = Item{
			ID:     c,
			Name:   fmt.Sprintf("Mug (%d)", i),
			Size:   [3]int{10, 8, 8},
			Weight: 300,
		}
		c++
		items[c] = Item{
			ID:     c,
			Name:   fmt.Sprintf("Water Bottle (%d)", i),
			Size:   [3]int{25, 7, 7},
			Weight: 200,
		}
		c++
		items[c] = Item{
			ID:     c,
			Name:   fmt.Sprintf("Phone Case (%d)", i),
			Size:   [3]int{15, 8, 3},
			Weight: 100,
		}
	}
	return items
}

const testFile = "./test_db.1e6.json"

func TestGenerateTestFile(t *testing.T) {
	t.Skip("only to generate a test file")
	items := generateTestData(t, 1e6)

	o, err := os.Create(testFile)
	testutil.Ok(t, err)
	defer func() { _ = o.Close() }()
	testutil.Ok(t, json.NewEncoder(o).Encode(items))
}

func TestLoad(t *testing.T) {
	d := &db{}
	for _, tc := range []func(string) error{
		d.load1,
		d.load2,
	} {
		t.Run("", func(t *testing.T) {
			inputFile := t.TempDir() + "/db_10k.json"

			items := generateTestData(t, 10e3)
			toCheck := items[39]

			o, err := os.Create(inputFile)
			testutil.Ok(t, err)
			defer func() { _ = o.Close() }()
			testutil.Ok(t, json.NewEncoder(o).Encode(items))

			d.loaded = nil
			testutil.Ok(t, tc(inputFile))
			testutil.Equals(t, toCheck, d.loaded[39])
		})
	}
}

func BenchmarkLoad(b *testing.B) {
	d := &db{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.loaded = nil
		testutil.Ok(b, d.load1(testFile))
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	// 10e6 = 9s -> Xeon 21s/35s
	// 10e3 = 8.5ms -> Xeon 21.5ms
	items := generateTestData(b, 10e3)
	buf := &bytes.Buffer{}
	testutil.Ok(b, json.NewEncoder(buf).Encode(items))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		items = nil
		testutil.Ok(b, json.Unmarshal(buf.Bytes(), &items))
	}
}
