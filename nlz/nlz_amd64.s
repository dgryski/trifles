// func nlz(x uint32) uint32
TEXT ·nlz(SB),4,$0-12
        MOVL  x+0(FP), AX
        TESTL AX, AX
        JZ zero
        BSRL  AX, AX
        SUBL  $31, AX
        NEGL AX
        MOVL AX, ret+8(FP)
        RET
zero:
        MOVL $32, ret+8(FP)
        RET
