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

## Export 1: Over-engineered Simple, but not optimised yet version.

```
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']-12         	      22	6457670042 ns/op	3213566347 B/op	28273104 allocs/op
```

```
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']-12         	      21	6949694906 ns/op	3213489541 B/op	28273365 allocs/op
```

GOMAXPROCS = 1

```
Export done in  7.931812877s exported 1000 series, 39840000 samples
output.parquet row count: 1898000
Total RowCount: 1898000
 
output.parquet: 25191410 bytes
Total Size: 25191410 bytes
```

Size of result: 25191410 bytes.
Max Working Set Memory: 1GB
Heap Max: 600MB
Allocs: 3.2MB

```
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']-12         	      21	6220319788 ns/op	3213490881 B/op	28273439 allocs/op
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']-12         	      21	6869591226 ns/op	3336041528 B/op	29544813 allocs/op
```

Main memory offender: 

```
github.com/xitongsys/parquet-go/marshal.(*ParquetStruct).Marshal
/home/bwplotka/go/pkg/mod/github.com/xitongsys/parquet-go@v1.6.0/marshal/marshal.go

  Total:      5.16GB     9.79GB (flat, cum)  8.83%
     66            .          .           	return nodes 
     67            .          .           } 
     68            .          .            
     69            .          .           type ParquetStruct struct{} 
     70            .          .            
     71            .          .           func (p *ParquetStruct) Marshal(node *Node, nodeBuf *NodeBufType) []*Node { 
     72            .          .           	var ok bool 
     73            .          .            
     74            .          .           	numField := node.Val.Type().NumField() 
     75       5.16GB     5.16GB           	nodes := make([]*Node, 0, numField) 
     76            .          .           	for j := 0; j < numField; j++ { 
     77            .     4.63GB           		tf := node.Val.Type().Field(j) 
     78            .          .           		name := tf.Name 
     79            .          .           		newNode := nodeBuf.GetNode() 
     80            .          .            
     81            .          .           		//some ignored item 
     82            .          .           		if newNode.PathMap, ok = node.PathMap.Children[name]; !ok { 
     83            .          .           			continue 
     84            .          .           		} 
     85            .          .            
     86            .          .           		newNode.Val = node.Val.Field(j) 
     87            .          .           		newNode.RL = node.RL 
     88            .          .           		newNode.DL = node.DL 
     89            .          .           		nodes = append(nodes, newNode) 
     90            .          .           	} 
     91            .          .           	return nodes 
     92            .          .           } 
```

Main CPU time user: Waiting for Store API response, Then chunk decoding

```
   Type: time
Showing nodes accounting for 26763.90s, 100% of 26763.90s total
----------------------------------------------------------+-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------+-------------
                                         11349.86s 42.81% |   runtime.selectgo
                                          9203.77s 34.71% |   runtime.chanrecv
                                          5843.87s 22.04% |   runtime.netpollblock
                                           117.16s  0.44% |   runtime.goparkunlock
 26514.66s 99.07% 99.07%  26514.66s 99.07%                | runtime.gopark
----------------------------------------------------------+-------------
                                            65.33s   100% |   github.com/efficientgo/examples/pkg/parquet-export/export1.everySample
    36.60s  0.14% 99.21%     65.33s  0.24%                | github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc.(*xorIterator).Next
                                            24.82s 37.99% |   github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc.(*xorIterator).readValue
                                             2.68s  4.10% |   encoding/binary.ReadVarint
                                             0.52s  0.79% |   encoding/binary.ReadUvarint
                                             0.32s  0.49% |   github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc.(*bstreamReader).readBit
                                             0.30s  0.46% |   github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc.(*bstreamReader).readBits
                                             0.10s  0.15% |   runtime.asyncPreemp
```

## Export 2: Simple, but not optimised yet version.

```
BenchmarkParquetExportIntegration/[type:NEQ__name:"__name__"]
BenchmarkParquetExportIntegration/[type:NEQ__name:"__name__"]-12         	       3	48267774030 ns/op	30471885085 B/op	262049551 allocs/op
BenchmarkParquetExportIntegration/[type:RE__name:"__name__"__value:"continuous_app_metric9.{1}"]
BenchmarkParquetExportIntegration/[type:RE__name:"__name__"__value:"continuous_app_metric9.{1}"]-12         	      31	4863932773 ns/op	3005606573 B/op	26183210 allocs/op
```

```
Export done in  5.298785971s exported 1000 series, 39840000 samples
output.parquet row count: 1898000
```

Size of result: 25191410 bytes.
Max Working Set Memory: 700MB
Heap Max: 470MB
Allocs: 3GB


GOMAXPROCS = 1


```
BenchmarkParquetExportIntegration/[__name__=~'continuous_app_metric9.{1}']-12         	      28	5086515883 ns/op	3005500617 B/op	26183167 allocs/op
```


