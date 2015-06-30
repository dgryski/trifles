// func nlz(x uint32) uint32
TEXT ·nlz(SB),4,$0-12
        BSRL  x+0(FP), AX
        JZ zero
        SUBL  $31, AX
        NEGL AX
        MOVL AX, ret+8(FP)
        RET
zero:
        MOVL $32, ret+8(FP)
        RET
