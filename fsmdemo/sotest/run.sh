#!/bin/sh

go build -buildmode=c-shared -o fsmso.so github.com/dgryski/trifles/fsmdemo/sotest/fsmso
make runner
./runner
