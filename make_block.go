// +build ignore

package main

import (
	"crypto/aes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	LaneSize  = 4 * aes.BlockSize
	BlockSize = 256
)

var output = flag.String("out", "block_amd64.s", "output filename")

func main() {
	flag.Parse()
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m := NewMeow(f)
	if err := m.Generate(); err != nil {
		log.Fatal(err)
	}
}

// Generator is an interface for assembly generation.
type Generator interface {
	inst(name, format string, args ...interface{})
	data64(name string, x []uint64)
}

// Array represents a byte array based at Offset.
type Array struct {
	Base   string
	Offset int
}

// Addr returns a reference to the given byte index into the array.
func (a Array) Addr(idx int) string {
	return fmt.Sprintf("%d(%s)", a.Offset+idx, a.Base)
}

// Slice returns an array based at position idx into the array.
func (a Array) Slice(idx int) Array {
	return Array{Base: a.Base, Offset: a.Offset + idx}
}

// StackFrame represents the stack frame of a function.
type StackFrame struct {
	Size int
}

// Alloc allocates an Array on the stack frame.
func (s *StackFrame) Alloc(size int) Array {
	a := Array{Base: "SP", Offset: s.Size}
	s.Size += size
	return a
}

// Backend encapsulates the instruction set we are targetting.
type Backend interface {
	// Width is the register width in bits.
	Width() int

	// Zero all streams.
	Zero()

	// AESLoad encrypts stream s with key from memory.
	AESLoad(s int, m Array)

	// AESMerge merges stream s with stream t.
	AESMerge(s, t int)

	// Store stream s to memory location m.
	Store(s int, m Array)
}

// AESNI implements the AES-NI backend.
type AESNI struct {
	g Generator
}

func NewAESNI(g Generator) *AESNI {
	return &AESNI{g: g}
}

func (a AESNI) Width() int { return 128 }

func (a AESNI) Zero() {
	for s := 0; s < 16; s++ {
		a.g.inst("PXOR", "X%d, X%d", s, s)
	}
}

func (a AESNI) AESLoad(s int, m Array) {
	a.g.inst("VAESDEC", "%s, X%d, X%d", m.Addr(0), s, s)
}

func (a AESNI) AESMerge(s, t int) {
	a.g.inst("VAESDEC", "X%d, X%d, X%d", t, s, s)
}

func (a AESNI) Store(i int, m Array) {
	a.g.inst("MOVOU", "X%d, %s", i, m.Addr(0))
}

//func (a *AESNI) StackAlloc(f *StackFrame) {
//	a.spill = f.Alloc(LaneSize)
//}
//
//func (a *AESNI) R0(g Generator) int {
//	a.StoreLane(g, 0, a.spill)
//	for i := 0; i < 4; i++ {
//		a.stream[i] = a.spill.Addr(16 * i)
//		a.stream = append(a.stream, fmt.Sprintf("X%d", i))
//	}
//	return 4
//}
//
//func (a AESNI) LoadLane(g Generator, m Array, i int) {
//	for j := 0; j < 4; j++ {
//		g.inst("MOVOU", "%s, %s", m.Addr(j*aes.BlockSize), a.stream[4*i+j])
//	}
//}
//
//
//
//func (a AESNI) AESMerge(g Generator, i, j int) {
//	for k := 0; k < 4; k++ {
//		g.inst("VAESDEC", "%s, %s, %s", a.stream[4*j+k], a.stream[4*i+k], a.stream[4*i+k])
//	}
//}
//
//func (a AESNI) Rotate(g Generator, i int) {
//	l := a.stream[4*i : 4*i+4]
//	t := l[0]
//	l[0] = l[1]
//	l[1] = l[2]
//	l[2] = l[3]
//	l[3] = t
//}

// Meow writes an assembly implementation of Meow hash components.
type Meow struct {
	w       io.Writer // where to write assembly output
	defines []string  // names of defined macros
	err     error     // saved error from writing
}

// NewMeow builds a new assembly builder writing to w.
func NewMeow(w io.Writer) *Meow {
	return &Meow{
		w: w,
	}
}

// Generate triggers assembly generation.
func (m *Meow) Generate() error {
	m.header()

	m.checksum(NewAESNI(m))

	return m.err
}

// checksum outputs the entire checksum function.
func (m *Meow) checksum(b Backend) {
	f := &StackFrame{}
	mixer := f.Alloc(aes.BlockSize)
	partial := f.Alloc(aes.BlockSize)

	name := fmt.Sprintf("checksum%d", b.Width())
	m.text(name, f.Size, 56)

	m.arg("seed", 0, "R8")
	m.arg("dst_ptr", 8, "DI")
	m.arg("src_ptr", 32, "SI")
	m.arg("src_len", 40, "AX")

	m.section("Allocate general purpose registers.")
	m.alloc("TOTAL_LEN", "R9")
	m.alloc("MIX0", "R10")
	m.alloc("MIX1", "R11")
	m.alloc("PARTIAL_PTR", "R12")
	m.alloc("TMP", "R13")
	m.alloc("ZERO", "R15")

	m.section("Prepare a zero register.")
	m.inst("XORQ", "ZERO, ZERO")

	m.section("Backup total input length.")
	m.inst("MOVQ", "SRC_LEN, TOTAL_LEN")

	m.section("Prepare Mixer.")
	m.inst("MOVQ", "SEED, MIX0")
	m.inst("SUBQ", "SRC_LEN, MIX0")
	m.inst("MOVQ", "SEED, MIX1")
	m.inst("ADDQ", "SRC_LEN, MIX1")
	m.inst("INCQ", "MIX1")

	for i := 0; i < 2; i++ {
		m.inst("MOVQ", "MIX%d, %s", i, mixer.Addr(8*i))
	}

	m.section("Load zero \"IV\".")
	b.Zero()

	m.section(fmt.Sprintf("Handle full %d-byte blocks.", BlockSize))
	m.label("loop")
	m.inst("CMPQ", "SRC_LEN, $%d", BlockSize)
	m.inst("JB", "sub256")

	m.section("Hash block.")
	src := Array{Base: "SRC_PTR"}
	for i := 0; i < 16; i++ {
		b.AESLoad(i, src.Slice(i*aes.BlockSize))
	}

	m.section("Update source pointer.")
	m.inst("ADDQ", "$%d, SRC_PTR", BlockSize)
	m.inst("SUBQ", "$%d, SRC_LEN", BlockSize)
	m.inst("JMP", "loop")

	m.section(fmt.Sprintf("Handle final sub %d-byte block.", BlockSize))
	m.label("sub256")
	for i := 0; i < BlockSize-aes.BlockSize; i += aes.BlockSize {
		m.inst("CMPQ", "SRC_LEN, $%d", aes.BlockSize)
		m.inst("JB", "sub16")
		b.AESLoad(i/aes.BlockSize, src)
		m.inst("ADDQ", "$%d, SRC_PTR", aes.BlockSize)
		m.inst("SUBQ", "$%d, SRC_LEN", aes.BlockSize)
	}

	m.section("Handle final sub 16-byte block.")
	m.label("sub16")
	m.inst("CMPQ", "SRC_LEN, $0")
	m.inst("JE", "combine")

	m.inst("MOVQ", "ZERO, %s", partial.Addr(0))
	m.inst("MOVQ", "ZERO, %s", partial.Addr(8))
	m.inst("LEAQ", "%s, PARTIAL_PTR", partial.Addr(0))

	m.inst("CMPQ", "TOTAL_LEN, $16")
	m.inst("JB", "byteloop")

	m.inst("LEAQ", "-16(SRC_PTR)(SRC_LEN*1), SRC_PTR")
	m.inst("MOVQ", "$16, SRC_LEN")

	m.label("byteloop")
	m.inst("MOVB", "(SRC_PTR), TMP")
	m.inst("MOVB", "TMP, (PARTIAL_PTR)")
	m.inst("INCQ", "SRC_PTR")
	m.inst("INCQ", "PARTIAL_PTR")
	m.inst("DECQ", "SRC_LEN")
	m.inst("JNE", "byteloop")

	b.AESLoad(15, partial)

	m.section("Combine.")
	m.label("combine")
	m0 := 7
	ordering := []int{10, 4, 5, 12, 8, 0, 1, 9, 13, 2, 6, 14, 3, 11, 15}
	for _, s := range ordering {
		b.AESMerge(m0, s)
	}

	m.section("Mixing.")
	for i := 0; i < 3; i++ {
		b.AESLoad(m0, mixer)
	}

	m.section("Store hash.")
	b.Store(m0, Array{Base: "DST_PTR"})

	m.inst("RET", "")
	m.undefall()
}

// header outputs the file header with code generation warning and standard header includes.
func (m *Meow) header() {
	_, self, _, _ := runtime.Caller(0)
	m.printf("// Code generated by go run %s. DO NOT EDIT.\n\n", filepath.Base(self))
	m.printf("// +build !noasm\n\n")
	m.printf("#include \"textflag.h\"\n")
}

// section marks a section of the code with a comment.
func (m *Meow) section(description string) {
	m.printf("\n\t// %s\n", description)
}

// label defines a label.
func (m *Meow) label(name string) {
	m.printf("\n%s:\n", name)
}

// data64 outputs a DATA section.
func (m *Meow) data64(name string, x []uint64) {
	for i := range x {
		m.printf("\nDATA %s<>+0x%02x(SB)/8, $0x%016x", name, 8*i, x[i])
	}
	m.printf("\nGLOBL %s<>(SB), (NOPTR+RODATA), $%d\n", name, 8*len(x))
}

// text defines a function header.
func (m *Meow) text(name string, frame, args int) {
	m.printf("\nTEXT %s,0,$%d-%d\n", local(name), frame, args)
}

// alloc informally "allocates" a register with a #define statement.
func (m *Meow) alloc(name, reg string) string {
	macro := strings.ToUpper(name)
	m.define(macro, reg)
	return macro
}

// define a macro.
func (m *Meow) define(name, value string) {
	m.printf("#define %s %s\n", name, value)
	m.defines = append(m.defines, name)
}

// undefall undefs all defined macros.
func (m *Meow) undefall() {
	for _, name := range m.defines {
		m.printf("#undef %s\n", name)
	}
	m.defines = nil
}

// arg reads an argument, and allocates a register for it.
func (m *Meow) arg(name string, offset int, reg string) {
	macro := m.alloc(name, reg)
	m.inst("MOVQ", "%s+%d(FP), %s", name, offset, macro)
}

// inst writes an instruction.
func (m *Meow) inst(name, format string, args ...interface{}) {
	args = append([]interface{}{name}, args...)
	m.printf("\t%-8s "+format+"\n", args...)
}

func (m *Meow) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(m.w, format, args...); err != nil {
		m.err = err
	}
}

// local returns a reference to a local symbol (primarily useful for the unicode dot).
func local(name string) string {
	return fmt.Sprintf("\u00b7%s(SB)", name)
}
