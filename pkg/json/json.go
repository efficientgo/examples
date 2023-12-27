package json

import (
	"encoding/json"
	"os"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
)

type Item struct {
	ID     int    `json:"id"`
	Name   string `json:"name"` // Max 32 chars.
	Size   [3]int `json:"size"` // Width, Height, Length.
	Weight int    `json:"weight"`
}

// Read more in "Efficient Go"; Example 3-3

type db struct {
	loaded []Item
}

func (d *db) load1(dbFile string) (err error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")

	return json.NewDecoder(f).Decode(&d.loaded)
}

func (d *db) load2(dbFile string) (err error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")

	st, err := f.Stat()
	if err != nil {
		return err
	}

	// Read at once, and unmarshal as a whole.
	b := make([]byte, int(st.Size()))
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	if n != int(st.Size()) {
		return errors.Wrapf(err, "read only %v/%v bytes", n, st.Size())
	}

	// Estimate number of items based on assumption that we each item has roughly
	// 88 bytes. This is best effort--in a worse case we over allocate a little,
	// or we will have to resize once.
	d.loaded = make([]Item, 0, int((st.Size())-24)/88)
	return json.Unmarshal(b, &d.loaded)
}
