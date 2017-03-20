TEXT ·rdtsc(SB),0,$0-8
        RDTSC
	SHLQ    $32, DX
	ADDQ    DX, AX
	MOVQ    AX, ret+0(FP)
        RET

