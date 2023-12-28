#!/bin/zsh

set -xe

export ver=v1 && \
 go test ./pkg/json -run '^$' -bench '^BenchmarkLoad' -benchtime 10s -count 6 \
 -cpu 4 \
 -benchmem \
 -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof \
 | tee ${ver}.txt
