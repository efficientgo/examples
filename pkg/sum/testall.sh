#!/bin/bash


# // export var=v1 && go test -count 5 -benchtime 5s -run '^$' -bench . -memprofile=${var}.mem.pprof -cpuprofile=${var}.cpu.pprof > ${var}.txt
# func BenchmarkSum(b *testing.B) {
# 	w, err := strconv.Atoi(os.Getenv("WORKERS"))
# 	testutil.Ok(b, err)
#
# 	b.ReportAllocs()
# 	b.ResetTimer()
# 	for i := 0; i < b.N; i++ {
# 		Answer, _ = ConcurrentSum("input.txt", w)
# 	}
# }

for GOMAXPROCS in 1 2 8 12 24
do
	for WORKERS in 1 8 12 24 48
  do
  	export var="vP${GOMAXPROCS}-W${WORKERS}" && export GOMAXPROCS=${GOMAXPROCS} && export WORKERS=${WORKERS} && \
  	  go test -count 5 -benchtime 10s -run '^$' -bench . > ${var}.txt
  done
done

