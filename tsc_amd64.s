#include "textflag.h"

TEXT rdtscp(SB),NOSPLIT,$0-8
    RDTSCP
    SHLQ    $32, DX
    ADDQ    DX, AX
    MOVQ    AX, ret+0(FP)
    MOVQ    CX, ret+8(FP)
    RET