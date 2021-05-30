package block_test

import (
	"context"
	"fmt"

	block "github.com/efficientgo/examples/pkg/godoc"
	"github.com/oklog/ulid"
)

func ExampleDownload() {
	if err := block.Download(context.Background(), ulid.MustNew(0, nil), "here"); err != nil {
		fmt.Println(err)
	}
	// Output: downloaded
}
