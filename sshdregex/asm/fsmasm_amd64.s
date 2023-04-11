//go:build amd64 && !purego

#include "textflag.h"

// func MatchASM(data []byte) int
TEXT    ·MatchASM(SB), NOSPLIT, $0-32
	MOVQ data_base+0(FP), DI
	MOVQ data_len+8(FP), SI
	ADDQ   DI, SI

state_0:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x73
	JNE   state_0

state_1:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x73
	JNE   state_0

state_2:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x68
	JE    state_3

	CMPL    DX, $0x73
	JE    state_2

	JMP   state_0

state_3:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x64
	JE    state_4

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_4:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x5b
	JE    state_5

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_5:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_0

	CMPL    DX, $0x39
	JBE   state_6

	CMPL    DX, $0x72
	JBE   state_0

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_6:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_0

	CMPL    DX, $0x39
	JBE   state_7

	CMPL    DX, $0x72
	JBE   state_0

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_7:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_0

	CMPL    DX, $0x39
	JBE   state_8

	CMPL    DX, $0x72
	JBE   state_0

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_8:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_0

	CMPL    DX, $0x39
	JBE   state_9

	CMPL    DX, $0x72
	JBE   state_0

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_9:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_0

	CMPL    DX, $0x39
	JBE   state_10

	CMPL    DX, $0x72
	JBE   state_0

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_10:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x5d
	JE    state_11

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_11:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x3a
	JE    state_12

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_12:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x09
	JE    state_12

	CMPL    DX, $0x20
	JE    state_12

	CMPL    DX, $0x46
	JE    state_13

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_13:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x61
	JE    state_14

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_14:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x69
	JE    state_15

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_15:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x6c
	JE    state_16

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_16:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x65
	JE    state_17

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_17:
	MOVQ    $-1, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x64
	JE    state_18

	CMPL    DX, $0x73
	JE    state_1

	JMP   state_0

state_18:
	MOVL   $0x12, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x73
	JNE   state_18

state_19:
	MOVL   $0x13, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x73
	JNE   state_18

state_20:
	MOVL   $0x14, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x68
	JE    state_21

	CMPL    DX, $0x73
	JE    state_20

	JMP   state_18

state_21:
	MOVL   $0x15, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x64
	JE    state_22

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_22:
	MOVL   $0x16, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x5b
	JE    state_23

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_23:
	MOVL   $0x17, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_18

	CMPL    DX, $0x39
	JBE   state_24

	CMPL    DX, $0x72
	JBE   state_18

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_24:
	MOVL   $0x18, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_18

	CMPL    DX, $0x39
	JBE   state_25

	CMPL    DX, $0x72
	JBE   state_18

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_25:
	MOVL   $0x19, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_18

	CMPL    DX, $0x39
	JBE   state_26

	CMPL    DX, $0x72
	JBE   state_18

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_26:
	MOVL   $0x1a, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_18

	CMPL    DX, $0x39
	JBE   state_27

	CMPL    DX, $0x72
	JBE   state_18

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_27:
	MOVL   $0x1b, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x2f
	JBE   state_18

	CMPL    DX, $0x39
	JBE   state_28

	CMPL    DX, $0x72
	JBE   state_18

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_28:
	MOVL   $0x1c, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x5d
	JE    state_29

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_29:
	MOVL   $0x1d, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x3a
	JE    state_30

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_30:
	MOVL   $0x1e, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x09
	JE    state_30

	CMPL    DX, $0x20
	JE    state_30

	CMPL    DX, $0x46
	JE    state_31

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_31:
	MOVL   $0x1f, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x61
	JE    state_32

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_32:
	MOVL   $0x20, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x69
	JE    state_33

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_33:
	MOVL   $0x21, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x6c
	JE    state_34

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

state_34:
	MOVL   $0x22, AX
	CMPQ    DI,SI
	JE      finish
	MOVBLZX (DI), DX
	ADDQ    $1, DI

	CMPL    DX, $0x73
	JE    state_19

	JMP   state_18

finish:
	MOVQ    AX, ret+24(FP)
	RET
