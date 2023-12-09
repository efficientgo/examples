#!/bin/zsh

set -xe

export ver=v0 && \
 go test ./pkg/json -run '^$' -bench '^BenchmarkSell' -benchtime 5s -count 6 \
 -cpu 4 \
 -benchmem \
 -memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof \
 | tee ${ver}.txt
