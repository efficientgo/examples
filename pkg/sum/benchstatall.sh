#!/bin/bash

ARGS=""

for GOMAXPROCS in 1 2 8 12 24
do
	for WORKERS in 1 8 12 24 48
  do
    ARGS="${ARGS} vP${GOMAXPROCS}-W${WORKERS}.txt"
  done
done

benchstat ${ARGS}
