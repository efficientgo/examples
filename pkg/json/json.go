package json

import (
	"encoding/json"
	"os"

	"github.com/efficientgo/core/errcapture"
)

type Item struct {
	ID           int    `json:"id"`
	Name         string `json:"name"` // Max 32 chars.
	Price        int    `json:"price"`
	Size         Size   `json:"size"`
	PackagedSize Size   `json:"packaged_size"`
	Weight       int    `json:"weight"`
	Available    int    `json:"available"`
	Sold         int    `json:"sold"`
}

type Size struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Depth  int `json:"depth"`
}

// Read more in "Efficient Go"; Example 3-3

func sell0(dbFile string) (err error) {
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
	if _, err := f.Read(b); err != nil {
		return err
	}

	items := make([]Item, 0, 10e3)
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	// For simplicity, let's assume we know it's ordered, we know it's the 39th
	// item, and we know it's still available.
	items[39].Available--
	items[39].Sold++

	// This will allocate a lot to grow.
	out, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile+".updated", out, os.ModePerm)
}

func sell(dbFile string) (err error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")

	var items []Item
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return err
	}

	// For simplicity, let's assume we know it's ordered, we know it's the 39th
	// item, and we know it's still available.
	items[39].Available--
	items[39].Sold++

	o, err := os.Create(dbFile + ".updated")
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, o.Close, "close")

	return json.NewEncoder(o).Encode(items)
}

func sell2(dbFile string) (err error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")

	items := make([]Item, 0, 10e3)
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return err
	}

	// For simplicity, let's assume we know it's ordered, we know it's the 39th
	// item, and we know it's still available.
	items[39].Available--
	items[39].Sold++

	o, err := os.Create(dbFile + ".updated")
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, o.Close, "close")

	return json.NewEncoder(o).Encode(items)
}

func sell3(dbFile string) (err error) {
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
	if _, err := f.Read(b); err != nil {
		return err
	}

	items := make([]Item, 0, 10e3)
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}

	// For simplicity, let's assume we know it's ordered, we know it's the 39th
	// item, and we know it's still available.
	items[39].Available--
	items[39].Sold++

	o, err := os.Create(dbFile + ".updated")
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, o.Close, "close")

	return json.NewEncoder(o).Encode(items)
}

func sell4(dbFile string) (err error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")

	items := make([]Item, 10e3)
	i := 0
	dec := json.NewDecoder(f)

	// TODO validation for wrong file format.
	if _, err := dec.Token(); err != nil {
		return err
	}
	for dec.More() {
		// This still allocates a lot, unfortunately we don't have access to internal
		// buffer of the decoder.
		if err := dec.Decode(&items[i]); err != nil {
			return err
		}
		i++
	}

	// For simplicity, let's assume we know it's ordered, we know it's the 39th
	// item, and we know it's still available.
	items[39].Available--
	items[39].Sold++

	o, err := os.Create(dbFile + ".updated")
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, o.Close, "close")

	// Can be further reduced by hand encoding in a streaming way.
	return json.NewEncoder(o).Encode(items)
}
