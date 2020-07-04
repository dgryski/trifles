#!/bin/bash

set -x -e

iter=30
items=10

input=../matcher/domains.txt
regex_opt=regex-opt

for i in $(seq $iter); do
	$regex_opt $(shuf $input |head -n $items |  tr '\n' '|' | sed 's/|$//; s/\./\\./g; ') >regexp.in
	re -r pcre -l go -k pair -e fish "$(cat regexp.in)" >fishfsm/fishfsm.go
	re -r pcre -l amd64_go -k pair -e fish "$(cat regexp.in)" >fishasm/fishasm.s
	re -r pcre -l vmc -k pair -e fish "$(cat regexp.in)" >fishrx.c
	printf 'package main; import "regexp"; var fishrx = regexp.MustCompile(`%s`);' "$(cat regexp.in)" > fishrx.go
	go build
	./testlibfsm
done
