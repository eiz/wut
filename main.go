package main

import (
	"unsafe"
)

var (
	statusInitializing = "Initializing..."
)

var (
	consoleRow = 0
	consoleCol = 0
)

var strBuffer [128]byte
var consoleBuffer []uint16
var biosMemory []byte
var rsdp *acpiRsdp
var mbhdr *multibootHeader

type sliceHeader struct {
	ptr      uintptr
	length   int
	capacity int
}

type acpiRsdp struct {
	signature   uint64
	checksum    byte
	oemid       [6]byte
	revision    byte
	rsdtAddress uint32
	length      uint32
	xsdtAddress uint64
	xchecksum   byte
}

type multibootHeader struct {
	flags           uint32
	memlower        uint32
	memupper        uint32
	bootdevice      uint32
	cmdline         uint32
	modscount       uint32
	modsaddr        uint32
	syms            [4]uint32
	mmaplength      uint32
	mmapaddr        uint32
	driveslength    uint32
	drivesaddr      uint32
	configtable     uint32
	bootloadername  uint32
	apmtable        uint32
	vbecontrolinfo  uint32
	vbemodeinfo     uint32
	vbemode         uint16
	vbeinterfaceseg uint16
	vbeinterfaceoff uint16
	vbeinterfacelen uint16
}

type e820Entry struct {
	size     uint32
	addrlo   uint32
	addrhi   uint32
	lengthlo uint32
	lengthhi uint32
	memtype  uint32
}

//go:nosplit
func out(foo ...interface{}) {
	for _, v := range foo {
		switch tv := v.(type) {
		case string:
			consoleWrite(tv)
		case []byte:
			consoleWriteBytes(tv)
		case uint32:
			consoleWriteUint64(uint64(tv))
		case uint64:
			consoleWriteUint64(tv)
		default:
			consoleWrite("<value>")
		}
	}
}

//go:nosplit
func formatUint64(v uint64) []byte {
	i := len(strBuffer)

	if v == 0 {
		i--
		strBuffer[i] = '0'
	}

	for v != 0 {
		i--
		n := byte(v & 0xF)

		if n >= 10 {
			strBuffer[i] = 'A' + (n - 10)
		} else {
			strBuffer[i] = '0' + n
		}

		v >>= 4
	}

	return strBuffer[i:]
}

//go:nosplit
func consoleInit() {
	var sl sliceHeader

	sl.ptr = 0xB8000
	sl.length = 80 * 25
	sl.capacity = 80 * 25
	consoleBuffer = *(*[]uint16)(unsafe.Pointer(&sl))
}

//go:nosplit
func findAcpi() {
	var sl sliceHeader

	sl.ptr = 0xE0000
	sl.length = 0x20000
	sl.capacity = 0x20000
	biosMemory = *(*[]byte)(unsafe.Pointer(&sl))

	for i := 0; i < len(biosMemory); i += 16 {
		if biosMemory[i] == 'R' && biosMemory[i+1] == 'S' &&
			biosMemory[i+2] == 'D' && biosMemory[i+3] == ' ' &&
			biosMemory[i+4] == 'P' && biosMemory[i+5] == 'T' &&
			biosMemory[i+6] == 'R' && biosMemory[i+7] == ' ' {

			ck := byte(0)

			for j := i; j < i+20; j++ {
				ck += biosMemory[j]
			}

			if ck == 0 {
				out("Found ACPI root @ ", uint64(0xE0000)+uint64(i), "\r\n")
				rsdp = (*acpiRsdp)(unsafe.Pointer(uintptr(0xE0000 + i)))
				out("XSDT: ", rsdp.xsdtAddress, "\r\n")
				out("RSDT: ", uint64(rsdp.rsdtAddress), "\r\n")
				out("Length: ", uint64(rsdp.length), "\r\n")
				out("Signature: ", rsdp.signature, "\r\n")
			}
		}
	}
}

//go:nosplit
func consoleClear() {
	l := len(consoleBuffer)

	for i := 0; i < l; i++ {
		consoleBuffer[i] = uint16(0x1700)
	}
}

//go:nosplit
func consoleWriteChar(chr byte) {
	if consoleCol == 80 {
		consoleRow++
		consoleCol = 0
	}

	if consoleRow == 25 {
		consoleRow = 0
	}

	if chr == '\n' {
		consoleRow++
		return
	} else if chr == '\r' {
		consoleCol = 0
		return
	}

	consoleBuffer[consoleRow*80+consoleCol] =
		uint16(0x1700) | uint16(chr)
	consoleCol++
}

//go:nosplit
func consoleWrite(str string) {
	l := len(str)

	for i := 0; i < l; i++ {
		consoleWriteChar(str[i])
	}
}

//go:nosplit
func consoleWriteBytes(b []byte) {
	l := len(b)

	for i := 0; i < l; i++ {
		consoleWriteChar(b[i])
	}
}

//go:nosplit
func consoleWriteUint64(v uint64) {
	consoleWriteBytes(formatUint64(v))
}

//go:nosplit
func dumpe820() {
	for i := uint32(0); i <= mbhdr.mmaplength-4; {
		entry := (*e820Entry)(unsafe.Pointer(uintptr(mbhdr.mmapaddr + i)))

		if i+entry.size > mbhdr.mmaplength || entry.size < 0x14 {
			break
		}

		out(" size ", uint64(entry.size),
			" addr ", uint64(entry.addrlo)|(uint64(entry.addrhi)<<32),
			" length ", uint64(entry.lengthlo)|(uint64(entry.lengthhi)<<32),
			" type ", uint64(entry.memtype), "\r\n")
		i += entry.size + 4
	}
}

//go:nosplit
//go:nowritebarrierrec
func earlyMain() {
	consoleInit()
	consoleClear()
	out("mbhdr flags: ", mbhdr.flags, "\r\n")
	out("mmap addr: ", mbhdr.mmapaddr, " length: ", mbhdr.mmaplength, "\r\n")
	dumpe820()
	findAcpi()
}
