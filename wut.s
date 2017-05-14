#include "textflag.h"
#include "x86.h"
#define GDT_ENTRIES 3

TEXT _load_kernel(SB),NOSPLIT,$0
    MOVL BX, ·mbhdr(SB)
    MOVL gdt32ptr<>(SB), GDTR
    MOVL $SEG_SEL(1, 0, 0), AX
    MOVW AX, SS
    MOVW AX, DS
    MOVW AX, ES
    MOVW AX, FS
    MOVW AX, GS

// JMP $SEG_SEL(2, 0, 0):_load_kernel32(SB)
    BYTE $0xEA
    LONG $_load_kernel32(SB)
    WORD $SEG_SEL(2, 0, 0)

TEXT _load_kernel32(SB),NOSPLIT,$0
// Setup bootstrap page table. We use 4 level page tables with a single large
// page identity-mapping the first 1GB of memory. This is only supported on
// fairly modern implementations of x86_64.

// Load root page table into CR3
    LEAL pgl4<>(SB), AX
    MOVL AX, CR3

// Add a single entry to level 4 table pointing to level 3
    LEAL pgl3<>(SB), AX
    ORL $(PG_P|PG_RW|PG_U), AX
    MOVL AX, pgl4<>(SB)

// Add a single entry to level 3 table mapping 0x00000000-0x3FFFFFFF
    MOVL $(PG_P|PG_RW|PG_PS|PG_U), AX
    MOVL AX, pgl3<>(SB)

// Enable PAE
    MOVL CR4, AX
    ORL $CR4_PAE, AX
    MOVL AX, CR4

// Enable long mode
    MOVL $IA32_EFER, CX
    RDMSR
    ORL $IA32_EFER_LME, AX
    WRMSR

// Enable paging
    MOVL CR0, AX
    ORL $CR0_PG, AX
    MOVL AX, CR0

// Disable compatibility mode and jump to 64bit entry point
    MOVL gdt64ptr<>(SB), GDTR
    MOVL $SEG_SEL(1, 0, 0), AX
    MOVW AX, FS
    MOVW AX, GS

//  JMP $SEG_SEL(2, 0, 0):_load_kernel64(SB)
    BYTE $0xEA
    LONG $_load_kernel64(SB)
    WORD $SEG_SEL(2, 0, 0)

TEXT gdt32ptr<>(SB),NOSPLIT,$0
    WORD $((GDT_ENTRIES * 8) - 1)
    LONG $gdt32<>(SB)

TEXT gdt32<>(SB),NOSPLIT,$0
    SEG_DESC(0, 0, 0)
    SEG_DESC(0, 0xFFFFFFFF, SEG_S|SEG_P|SEG_G|SEG_DB|SEG_DATA)
    SEG_DESC(0, 0xFFFFFFFF, SEG_S|SEG_P|SEG_G|SEG_DB|SEG_EXEC)

TEXT gdt64ptr<>(SB),NOSPLIT,$0
    WORD $((GDT_ENTRIES * 8) - 1)
    LONG $gdt64<>(SB)

TEXT gdt64<>(SB),NOSPLIT,$0
    SEG_DESC(0, 0, 0)
    SEG_DESC(0, 0xFFFFFFFF, SEG_S|SEG_P|SEG_G|SEG_DB|SEG_DATA)
    SEG_DESC(0, 0xFFFFFFFF, SEG_S|SEG_P|SEG_G|SEG_L|SEG_EXEC)

GLOBL pgl4<>(SB),ALIGN4K|NOPTR,$4096
GLOBL pgl3<>(SB),ALIGN4K|NOPTR,$4096
