# Example of different Thanos StoreAPI -> Parquet transformations.

Dataset, depending on matchers: 

A.`[type:RE__name:"__name__"__value:"continuous_app_metric9.{1}"]` 10%: `[]*export1.LabelMatcher{{Name: "__name__", Value: "continuous_app_metric9.{1}", Type: export1.LabelMatcher_RE}` - 1k series, 39 840 000 samples, every 15s. This exports into output parquet file 25MB size with 1 898 000 rows.
B.`[type:NEQ__name:"__name__"]` Everything: `[]*export1.LabelMatcher{{Name: "__name__", Value: "", Type: export1.LabelMatcher_NEQ}}`. 10k series, 398 400 000 samples, every 15s. This exports into output parquet file 241MB size with 18 980 000 rows.

Benchmarked on Hardware:

```
goos: linux
goarch: amd64
pkg: github.com/efficientgo/examples/pkg/parquet-export
cpu: Intel(R) Core(TM) i7-9850H CPU @ 2.60GHz
```

## Export 1: The most "naive" version.

```
BenchmarkParquetExportIntegration/[type:NEQ__name:"__name__"]
BenchmarkParquetExportIntegration/[type:NEQ__name:"__name__"]-12         	       3	48267774030 ns/op	30471885085 B/op	262049551 allocs/op
BenchmarkParquetExportIntegration/[type:RE__name:"__name__"__value:"continuous_app_metric9.{1}"]
BenchmarkParquetExportIntegration/[type:RE__name:"__name__"__value:"continuous_app_metric9.{1}"]-12         	      31	4863932773 ns/op	3005606573 B/op	26183210 allocs/op
```




