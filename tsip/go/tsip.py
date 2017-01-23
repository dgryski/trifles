import peachpy.x86_64

k0 = Argument(uint64_t)
k1 = Argument(uint64_t)
p_base = Argument(ptr())
p_len = Argument(int64_t)
p_cap = Argument(int64_t)

def sipround(v0, v1):
    ADD(v0, v1)

    ROL(v1, 13)
    XOR(v1, v0)

    ROL(v0, 35)
    ADD(v0, v1)

    ROL(v1, 17)
    XOR(v1, v0)

    ROL(v0, 21)

with Function("HashASM", (k0, k1, p_base, p_len, p_cap), uint64_t, target=uarch.default) as function:

    reg_v0 = GeneralPurposeRegister64()
    reg_v1 = GeneralPurposeRegister64()

    LOAD.ARGUMENT(reg_v0, k0)
    LOAD.ARGUMENT(reg_v1, k1)

    reg_magic = GeneralPurposeRegister64()
    MOV(reg_magic, 0x736f6d6570736575)
    XOR(reg_v0, reg_magic)
    MOV(reg_magic, 0x646f72616e646f6d)
    XOR(reg_v1, reg_magic)

    reg_p = GeneralPurposeRegister64()
    reg_p_len = GeneralPurposeRegister64()
    LOAD.ARGUMENT(reg_p, p_base)
    LOAD.ARGUMENT(reg_p_len, p_len)

    reg_b = GeneralPurposeRegister64()
    MOV(reg_b, reg_p_len)
    SHL(reg_b, 56)

    reg_m = GeneralPurposeRegister64()

    loop = Loop()
    CMP(reg_p_len, 8)
    JL(loop.end)
    with loop:
        MOV(reg_m, [reg_p])

        XOR(reg_v1, reg_m)
        sipround(reg_v0, reg_v1)
        XOR(reg_v0, reg_m)

        ADD(reg_p, 8)
        SUB(reg_p_len, 8)
        CMP(reg_p_len, 8)
        JGE(loop.begin)

    # no support for jump tables
    labels = [Label("sw%d" % i) for i in range(0, 8)]

    for i in range(0,7):
        CMP(reg_p_len, i)
        JE(labels[i])

    char = GeneralPurposeRegister64()
    for i in range(7,0,-1):
        LABEL(labels[i])
        MOVZX(char, byte[reg_p+i-1])
        SHL(char, (i-1)*8)
        OR(reg_b, char)

    LABEL(labels[0])

    XOR(reg_v1, reg_b)
    sipround(reg_v0, reg_v1)
    XOR(reg_v0, reg_b)

    XOR(reg_v1, 0xff)
    sipround(reg_v0, reg_v1)
    ROL(reg_v1, 32)
    sipround(reg_v0, reg_v1)
    ROL(reg_v1, 32)

    XOR(reg_v0, reg_v1)

    RETURN(reg_v0)
