#define MASK(n) ((1 << (n)) - 1)
#define NBITS(v, off, n) (((v) >> (off)) & MASK(n))

#define SEG_S (1 << 12)
#define SEG_P (1 << 15)
#define SEG_G (1 << 23)
#define SEG_DB (1 << 22)
#define SEG_L (1 << 21)
#define SEG_TYPE(n) ((n) << 8)
#define SEG_DATA SEG_TYPE(2)
#define SEG_EXEC SEG_TYPE(10)
#define SEG_DESC(base, limit, flags) \
    LONG $((NBITS(base, 0, 16) << 16) | NBITS(limit, 0, 16)) \
    LONG $(NBITS(base, 16, 8) | flags | (NBITS(limit, 16, 4) << 16) | (NBITS(base, 24, 8) << 24))
#define SEG_SEL(index, ti, rpl) (((index) << 3) | ((ti) << 2) | (rpl))

#define PG_P (1 << 0)
#define PG_RW (1 << 1)
#define PG_U (1 << 2)
#define PG_PS (1 << 7)

#define CR0_PG (1 << 31)
#define CR4_PAE (1 << 5)
#define CR4_OSFXSR (1 << 9)
#define CR4_OSXMMEXCPT (1 << 10)

#define IA32_EFER 0xC0000080
#define IA32_EFER_LME (1 << 8)
