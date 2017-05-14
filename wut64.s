#include "textflag.h"
#include "x86.h"

TEXT _load_kernel64(SB),NOSPLIT,$0
    LEAQ bootstack<>(SB), SP
    ADDQ $4096-8, SP
    MOVQ CR4, AX
    ORQ $(CR4_OSFXSR|CR4_OSXMMEXCPT), AX
    MOVQ AX, CR4
    CALL ·earlyMain(SB)
idle:
    HLT
    JMP idle

GLOBL bootstack<>(SB),ALIGN4K|NOPTR,$4096
