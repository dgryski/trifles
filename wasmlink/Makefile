export PATH := /usr/local/Cellar/llvm/16.0.1/bin:$(PATH)
TINYGOROOT=$(shell tinygo env TINYGOROOT)
SHELL := /bin/sh
MAIN=wasmlink
OBJS=main.o task.o runtime.o add.o

build:
	@ PATH=$(PATH) make build_all

build_all: main.o task.o runtime.o add.o

link_dog:
	EXTRA=rematch-dog/matcher.o $(MAKE) link

link_cat:
	EXTRA=rematch-cat/matcher.o $(MAKE) link

link: $(OBJS) $(EXTRA)
	wasm-ld --stack-first --no-demangle -o linked.wasm $(OBJS) $(EXTRA) /Users/dgryski/go/src/github.com/tinygo-org/tinygo/lib/wasi-libc/sysroot/lib/wasm32-wasi/libc.a 
	wasm-opt --asyncify linked.wasm -o linked-asyncify.wasm && rm -f linked.wasm
	mv linked-asyncify.wasm $(MAIN).wasm

task.o: $(TINYGOROOT)/src/internal/task/task_asyncify_wasm.S
	clang --target=wasm32-unknown-wasi -c -o $@ $<

runtime.o: $(TINYGOROOT)/src/runtime/asm_tinygowasm.S
	clang --target=wasm32-unknown-wasi -c -o $@ $<

add.o: add.wat
	wat2wasm --relocatable add.wat -o add.o

main.o: main.go
	tinygo build -target=wasi -o main.o

clean:
	rm -f main.o add.o runtime.o task.o linked.wasm linked-asyncify.wasm wasmlink.wasm

rematch-cat/matcher.o: rematch-cat/matcher.c
	clang --target=wasm32-unknown-wasi -c -o $@ $<
rematch-dog/matcher.o: rematch-dog/matcher.c
	clang --target=wasm32-unknown-wasi -c -o $@ $<

matchers:
	re -r pcre -k pair -e matcher -pl vmc meow >rematch-cat/matcher.c
	re -r pcre -k pair -e matcher -pl vmc woof >rematch-dog/matcher.c
