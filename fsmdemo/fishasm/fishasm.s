#include "textflag.h"

// func Match(data string) int
TEXT    ·Match(SB), NOSPLIT, $0-24
	MOVQ data_base+0(FP), DI
	MOVQ data_len+8(FP), SI
	ADDQ   DI, SI

state_0:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x66
	JNE   state_0

state_1:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x66
	JE    state_1

	CMPL    DX, $0x69
	JNE   state_0

state_2:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x66
	JE    state_1

	CMPL    DX, $0x69
	JE    state_2

	CMPL    DX, $0x73
	JNE   state_0

// state 3
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x66
	JE    state_1

	CMPL    DX, $0x68
	JNE   state_0

// state 4
	MOVQ   $0x04, AX
	// elided jmp to finish

finish:
	MOVQ    AX, ret+16(FP)
	RET
