TEXT ·rdtsc(SB),0,$0-8
        RDTSC
        MOVL AX, 8(SP)
        MOVL DX, 12(SP)
        RET

