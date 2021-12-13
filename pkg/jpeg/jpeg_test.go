package jpeg

import (
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func BenchmarkDecEnc(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testutil.Ok(b, decodeEncode("./photo8k.jpg"))
	}
}
