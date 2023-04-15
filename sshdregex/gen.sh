re -r pcre -k pair -e '' -E main -pl go 'sshd\[\d{5}\]:\s*Failed' |sed -e 's/package mainfsm/package main/' >fsm.go
re -r pcre -k pair -e 'Unsafe' -E main -pl go_unsafe 'sshd\[\d{5}\]:\s*Failed' |sed -e 's/package mainfsm/package main/' >fsmunsafe.go
re -r pcre -k pair -e sshd -pl vmc 'sshd\[\d{5}\]:\s*Failed' > c/fsmc.c
re -r pcre -k pair -e '' -E asm -pl amd64_go 'sshd\[\d{5}\]:\s*Failed' > asm/fsmasm_amd64.s
ragel-go -G2 ragel.rl && gofmt -w ragel.go
