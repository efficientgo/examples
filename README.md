# "Efficient Go" Book Code Examples

Hi 👋

![img.png](book.png)

My name is [Bartek Płotka](https://www.bwplotka.dev) and I wrote ["Efficient Go"](https://www.bwplotka.dev/book). This book teaches pragmatic approaches to software efficiency and optimizations. While the majority of the learnings works for any programming language, it's best to learn on specific examples. For that purpose I teach how to make my favorite language more efficient - [Go](https://go.dev).

In this open-source repository you can find all examples from the book with tests, additional comments and more! Play with the examples to learn more about CPU, memory, OS and Go runtime characteristics I explained in my book! 

> NOTE: Don't import this module to your production code--it is meant for learning purposes only. Instead, we maintain production grade utilities mentioned in the book in [core](https://github.com/efficientgo/core) and [e2e](https://github.com/efficientgo/e2e) modules.

## Index of Examples From Book

All examples from the book are buildable and tested in CI (as they should). See their location in the below table. For most of the code there exists corresponding `_test.go` file with tests and microbenchmarks when relevant.

For the book version before errata, see [v1.0 branch](https://github.com/efficientgo/examples/tree/v1).

> NOTE: Some function names in examples might be different in book vs in code, due to name clashes. Code in the book might be also simplified.
> I also use non-conventional naming convention with some functions e.g. `FailureRatio_Better`. In Go all names
> should have camelCase form, so e.g. `FailureRatioBetter`. However, I chose to keep underscore to separate different versions
> of the same functions for book purposes - so the "production name" should is still `FailureRatio`! (: 

| Example Ref    | Page    | Path to code/function in this repository.                                                                                               |
|----------------|---------|-----------------------------------------------------------------------------------------------------------------------------------------|
| Example 1-1    | 9       | [pkg/getter/getter.go:13 `FailureRatio`](pkg/getter/getter.go#L13)                                                                      |
| Example 1-2    | 10      | [pkg/getter/getter.go:29 `FailureRatio_Better`](pkg/getter/getter.go#L29)                                                               |
| Example 1-3    | 12      | [pkg/prealloc/slice.go:5 `createSlice`](pkg/prealloc/slice.go#L5)                                                                       |
| Example 1-4    | 12      | [pkg/prealloc/slice.go:14 `createSlice_Better`](pkg/prealloc/slice.go#L14)                                                              |
| Example 1-5    | 13      | [pkg/notationhungarian/hung.go:5 `structSystem`](pkg/notationhungarian/hung.go#L5)                                                      |
| Example 2-1    | 37      | [pkg/basic/basic.go:6](pkg/basic/basic.go#L6)                                                                                           |
| Example 2-2    | 42      | [pkg/export/export.go:6](pkg/export/export.go#L6)                                                                                       |
| Example 2-3    | 44      | [Prometheus main.go](https://github.com/prometheus/prometheus/blob/6b53aeb012080ab2a50acd229bfe9943125abfa6/cmd/prometheus/main.go#L15) |
| Example 2-4    | 48      | [pkg/errors/errors.go:8](pkg/errors/errors.go#L8)                                                                                       |
| Example 2-5    | 49      | [pkg/errors/errors.go:34](pkg/errors/errors.go#L34)                                                                                     |
| Example 2-6    | 51      | [pkg/basicserver/basicserver.go:10](pkg/basicserver/basicserver.go#L10)                                                                 |
| Example 2-7    | 52      | [pkg/unused/unused.go:8](pkg/unused/unused.go#L8)                                                                                       |
| Example 2-8    | 53-54   | [pkg/testing/max_test.go:11](pkg/testing/max_test.go#L11)                                                                               |
| Example 2-9    | 55      | [pkg/godoc](pkg/godoc)                                                                                                                  |  
| Example 2-10   | 56      | [pkg/godoc](pkg/godoc)                                                                                                                  |  
| Example 2-11   | 59      | [pkg/oop/oop.go:14](pkg/oop/oop.go#L14)                                                                                                 |  
| Example 2-12   | 62      | [sort.Interface from the standard library](https://go.dev/src/sort/sort.go)                                                             |  
| Example 2-13   | 62      | [pkg/oop/oop.go:58](pkg/oop/oop.go#L58)                                                                                                 |  
| Example 2-14   | 64      | [pkg/generics/sort.go:12](pkg/generics/sort.go#L12)                                                                                     |  
| Example 2-15   | 65      | [pkg/generics/blocks.go:19](pkg/generics/blocks.go#L19)                                                                                 |  
| Example 3-3    | 93      | [pkg/jpeg/jpeg.go:12](pkg/jpeg/jpeg.go#L12)                                                                                             |  
| Example 4-1    | 115     | [pkg/sum/sum.go:15 `Sum`](pkg/sum/sum.go#L15)                                                                                           |  
| Example 4-5    | 139     | [pkg/concurrency/concurrency.go:12](pkg/concurrency/concurrency.go#L12)                                                                 |  
| Example 4-6    | 140     | [pkg/concurrency/concurrency.go:36](pkg/concurrency/concurrency.go#L36)                                                                 |  
| Example 4-7    | 140     | [pkg/concurrency/concurrency.go:55](pkg/concurrency/concurrency.go#55)                                                                  |  
| Example 4-8    | 140     | [pkg/concurrency/concurrency.go:77](pkg/concurrency/concurrency.go#77)                                                                  |  
| Example 5-1    | 162     | [pkg/memory/mmap/mmap.go:14](pkg/memory/mmap/mmap.go#L14)                                                                               |  
| Example 5-2    | 164     | [pkg/memory/mmap/interactive/interactive_open.go:94](pkg/memory/mmap/interactive/interactive_open.go#L94)                               |  
| Example 5-3    | 165     | [pkg/memory/mmap/interactive/interactive_mmap.go:14](pkg/memory/mmap/interactive/interactive_mmap.go#L14)                               |  
| Example 5-4    | 179     | [pkg/memory/vars/vars.go:11](pkg/memory/vars/vars.go#L11)                                                                               |  
| Example 5-5    | 183     | [pkg/memory/slice/slice.go:30](pkg/memory/slice/slice.go#L30)                                                                           |  
| Example 5-6    | 187     | [pkg/memory/slice/slice.go:30](pkg/memory/slice/slice.go#L30)                                                                           |  
| Example 6-1    | 199     | [pkg/metrics/latency_test.go:31](pkg/metrics/latency_test.go#L31)                                                                       |  
| Example 6-2    | 200-201 | [pkg/metrics/latency_test.go:50](pkg/metrics/latency_test.go#L50)                                                                       |  
| Example 6-3    | 201     | [pkg/metrics/latency_test.go:74](pkg/metrics/latency_test.go#74)                                                                        |  
| Example 6-4    | 203     | [pkg/metrics/latency_test.go:85](pkg/metrics/latency_test.go#85)                                                                        |  
| Example 6-6    | 206     | [pkg/metrics/latency_test.go:107](pkg/metrics/latency_test.go#107)                                                                      |
| Example 6-7    | 211-212 | [pkg/metrics/latency_test.go:129](pkg/metrics/latency_test.go#129)                                                                      |
| Example 6-9    | 214     | [pkg/metrics/prom.yaml](pkg/metrics/prom.yaml)                                                                                          |
| Example 6-11   | 231     | [pkg/metrics/cpu_test.go:19](pkg/metrics/cpu_test.go#L19)                                                                               |
| Example 6-12   | 235     | [pkg/metrics/mem_test.go:69 `printMemRuntimeMetric`](pkg/metrics/mem_test.go#L69)                                                       |
| Example 7-1    | 241     | [pkg/sum/sum.go:15 `Sum`](pkg/sum/sum.go#L15)                                                                                           |
| Example 8-1    | 279     | [pkg/sum/sum_test.go:49 `BenchmarkSum`](pkg/sum/sum_test.go#L49)                                                                        |
| Example 8-2    | 281     | [pkg/sum/sum_test.go:49 `BenchmarkSum`](pkg/sum/sum_test.go#L49)                                                                        |
| Example 8-9    | 291     | [pkg/sum/sum_test.go:76 `TestSum`](pkg/sum/sum_test.go#L76)                                                                             |
| Example 8-10   | 292     | [pkg/sum/sum_test.go:49 `BenchmarkSum`](pkg/sum/sum_test.go#L49)                                                                        |
| Example 8-11   | 293     | [pkg/sum/sum_test.go:116 `TestBenchSum`](pkg/sum/sum_test.go#L116)                                                                      |
| Example 8-12   | 295     | [pkg/sum/sum_test.go:49 `BenchmarkSum`](pkg/sum/sum_test.go#L49)                                                                        |
| Example 8-13   | 296     | [pkg/sum/sum_test.go:25 `lazyCreateTestInput`](pkg/sum/sum_test.go#L25)                                                                 |
| Example 8-14   | 297     | [pkg/sum/sum_test.go:155 `BenchmarkSum_AcrossInputs`](pkg/sum/sum_test.go#L155)                                                         |
| Example 8-16   | 302     | [pkg/compileroptimizeaway/opt_away_test.go:11 `BenchmarkPopcnt_Wrong`](pkg/compileroptimizeaway/opt_away_test.go#11)                    |
| Example 8-18   | 304-305 | [pkg/compileroptimizeaway/opt_away_test.go:37 `BenchmarkPopcnt_Sink`](pkg/compileroptimizeaway/opt_away_test.go#L37)                    |
| Example 8-19   | 312-313 | [pkg/sum/labeler/labeler_e2e_test.go:39 `TestLabeler_LabelObject`](pkg/sum/labeler/labeler_e2e_test.go#L39)                             |
| Example 8-20   | 313-314 | [pkg/sum/labeler/labeler_e2e_test.go:39 `TestLabeler_LabelObject`](pkg/sum/labeler/labeler_e2e_test.go#L39)                             |
| Example 9-1    | 333     | [pkg/profile/fd/fd.go](pkg/profile/fd/fd.go)                                                                                            |
| Example 9-2    | 334-336 | [pkg/profile/fd/example/main.go](pkg/profile/fd/example/main.go)                                                                        |
| Example 9-4    | 344     | [pkg/profile/fd/example/main.go:52](pkg/profile/fd/example/main.go#L52)                                                                 |
| Example 9-5    | 358     | [pkg/profile/fd/http.go:12](pkg/profile/fd/http.go#L12)                                                                                 |
| Example 9-6    | 374-375 | [pkg/sum/labeler/labeler_e2e_test.go:39 `TestLabeler_LabelObject`](pkg/sum/labeler/labeler_e2e_test.go#L39)                             |
| Example 10-2   | 385     | [pkg/sum/sum_test.go:142 `BenchmarkSum_fgprof`](pkg/sum/sum_test.go#L142)                                                               |
| Example 10-3   | 388     | [pkg/sum/sum.go:43 `Sum2`](pkg/sum/sum.go#L43)                                                                                          |
| Example 10-4   | 390-391 | [pkg/sum/sum.go:94 `Sum3`](pkg/sum/sum.go#L94)                                                                                          |
| Example 10-5   | 393     | [pkg/sum/sum.go:143 `Sum4`](pkg/sum/sum.go#L143)                                                                                        |
| Example 10-7   | 396     | [pkg/sum/sum.go:191 `Sum5`](pkg/sum/sum.go#L191)                                                                                        |
| Example 10-8   | 398-399 | [pkg/sum/sum.go:252 `Sum6`](pkg/sum/sum.go#L252)                                                                                        | 
| Example 10-10  | 403     | [pkg/sum/sum_concurrent.go:18 `ConcurrentSum1`](pkg/sum/sum_concurrent.go#L18)                                                          | 
| Example 10-11  | 405-406 | [pkg/sum/sum_concurrent.go:49 `ConcurrentSum2`](pkg/sum/sum_concurrent.go#L49)                                                          | 
| Example 10-12  | 407-408 | [pkg/sum/sum_concurrent.go:109 `ConcurrentSum3`](pkg/sum/sum_concurrent.go#L109)                                                        | 
| Example 10-13  | 410     | [pkg/sum/sum_concurrent.go:181 `ConcurrentSum4`](pkg/sum/sum_concurrent.go#L181)                                                        | 
| Example 10-15  | 412     | [pkg/sum/sum.go:299 `Sum7`](pkg/sum/sum.go#L299)                                                                                        |
| Example 11-1   | 418     | [pkg/generics/dups.go:6](pkg/generics/dups.go#L6)                                                                                       |
| Example 11-2   | 428-429 | [pkg/leak/http_close.go:13](pkg/leak/http.go#L13)                                                                                       |
| Example 11-3   | 429-430 | [pkg/leak/http_close_test.go:17 `TestHandleCancel`](pkg/leak/http_test.go#L17)                                                          |
| Example 11-5   | 431-432 | [pkg/leak/http_close.go:39](pkg/leak/http.go#L39)                                                                                       |
| Example 11-6   | 433     | [pkg/leak/http_close.go:87 `Handle_Better`](pkg/leak/http.go#L87)                                                                       |
| Example 11-7   | 434     | [pkg/leak/http_close_test.go:81](pkg/leak/http_test.go#L81)                                                                             |
| Example 11-8   | 435-436 | [pkg/leak/file.go:15](pkg/leak/file.go#L15)                                                                                             |
| Example 11-9   | 437     | [pkg/leak/file.go:55](pkg/leak/file.go#L55)                                                                                             |
| Example 11-10  | 438-439 | [pkg/leak/http_exhaust.go:13](pkg/leak/http_exhaust.go#L13)                                                                             |
| Example 11-11  | 441     | [pkg/prealloc/prealloc_test.go:54](pkg/prealloc/prealloc_test.go#L54)                                                                   |
| Example 11-12  | 442-443 | [pkg/prealloc/prealloc_test.go:99](pkg/prealloc/prealloc_test.go#L99)                                                                   |
| Example 11-14  | 444-445 | [pkg/prealloc/linkedlist.go:6](pkg/prealloc/linkedlist.go#L6)                                                                           |
| Example 11-15  | 446     | [pkg/prealloc/linkedlist.go:39](pkg/prealloc/linkedlist.go#L39)                                                                         | 
| Example 11-16  | 448     | [pkg/prealloc/linkedlist.go:57](pkg/prealloc/linkedlist.go#L57)                                                                         | 
| Examples 11-17 | 450     | [pkg/pools/reuse.go:8](pkg/pools/reuse.go#L8)                                                                                           | 
| Examples 11-18 | 451-452 | [pkg/pools/reuse.go:8](pkg/pools/reuse.go#L8)                                                                                           | 
| Examples 11-19 | 452-453 | [pkg/pools/reuse.go:8](pkg/pools/reuse.go#L8)                                                                                           | 
| Examples 11-20 | 453     | [pkg/pools/reuse_test.go:39](pkg/pools/reuse_test.go#L39)                                                                               |
