// +build !noasm

// Code generated by go run make_block.go. DO NOT EDIT.

#include "textflag.h"

TEXT ·checksumAsm(SB),0,$272-56
#define SEED R8
	MOVQ     seed+0(FP), SEED
#define DST_PTR DI
	MOVQ     dst_ptr+8(FP), DST_PTR
#define SRC_PTR SI
	MOVQ     src_ptr+32(FP), SRC_PTR
#define SRC_LEN AX
	MOVQ     src_len+40(FP), SRC_LEN

	// Prepare IV.
#define IV R9
	MOVQ     SEED, IV
	MOVQ     IV, 0(SP)
	ADDQ     SRC_LEN, IV
	INCQ     IV
	MOVQ     IV, 8(SP)

	// Load IV.
	MOVOU    0(SP), X0
	MOVOU    0(SP), X1
	MOVOU    0(SP), X2
	MOVOU    0(SP), X3
	MOVOU    0(SP), X4
	MOVOU    0(SP), X5
	MOVOU    0(SP), X6
	MOVOU    0(SP), X7
	MOVOU    0(SP), X8
	MOVOU    0(SP), X9
	MOVOU    0(SP), X10
	MOVOU    0(SP), X11
	MOVOU    0(SP), X12
	MOVOU    0(SP), X13
	MOVOU    0(SP), X14
	MOVOU    0(SP), X15

loop:
	CMPQ     SRC_LEN, $256
	JL       residual

	// Hash block.
	VAESDEC  0(SRC_PTR), X0, X0
	VAESDEC  16(SRC_PTR), X1, X1
	VAESDEC  32(SRC_PTR), X2, X2
	VAESDEC  48(SRC_PTR), X3, X3
	VAESDEC  64(SRC_PTR), X4, X4
	VAESDEC  80(SRC_PTR), X5, X5
	VAESDEC  96(SRC_PTR), X6, X6
	VAESDEC  112(SRC_PTR), X7, X7
	VAESDEC  128(SRC_PTR), X8, X8
	VAESDEC  144(SRC_PTR), X9, X9
	VAESDEC  160(SRC_PTR), X10, X10
	VAESDEC  176(SRC_PTR), X11, X11
	VAESDEC  192(SRC_PTR), X12, X12
	VAESDEC  208(SRC_PTR), X13, X13
	VAESDEC  224(SRC_PTR), X14, X14
	VAESDEC  240(SRC_PTR), X15, X15

	// Update source pointer.
	ADDQ     $256, SRC_PTR
	SUBQ     $256, SRC_LEN
	JMP      loop

residual:
	CMPQ     SRC_LEN, $0
	JE       finish

	// Duplicate IV.
	MOVQ     0(SP), R10
	MOVQ     8(SP), R11
	MOVQ     R10, 16(SP)
	MOVQ     R11, 24(SP)
	MOVQ     R10, 32(SP)
	MOVQ     R11, 40(SP)
	MOVQ     R10, 48(SP)
	MOVQ     R11, 56(SP)
	MOVQ     R10, 64(SP)
	MOVQ     R11, 72(SP)
	MOVQ     R10, 80(SP)
	MOVQ     R11, 88(SP)
	MOVQ     R10, 96(SP)
	MOVQ     R11, 104(SP)
	MOVQ     R10, 112(SP)
	MOVQ     R11, 120(SP)
	MOVQ     R10, 128(SP)
	MOVQ     R11, 136(SP)
	MOVQ     R10, 144(SP)
	MOVQ     R11, 152(SP)
	MOVQ     R10, 160(SP)
	MOVQ     R11, 168(SP)
	MOVQ     R10, 176(SP)
	MOVQ     R11, 184(SP)
	MOVQ     R10, 192(SP)
	MOVQ     R11, 200(SP)
	MOVQ     R10, 208(SP)
	MOVQ     R11, 216(SP)
	MOVQ     R10, 224(SP)
	MOVQ     R11, 232(SP)
	MOVQ     R10, 240(SP)
	MOVQ     R11, 248(SP)
	MOVQ     R10, 256(SP)
	MOVQ     R11, 264(SP)
#define BLOCK_PTR BX
	LEAQ     16(SP), BLOCK_PTR

byteloop:
	MOVB     (SRC_PTR), R10
	MOVB     R10, (BX)
	INCQ     SRC_PTR
	INCQ     BLOCK_PTR
	DECQ     SRC_LEN
	JNE      byteloop
	VAESDEC  16(SP), X0, X0
	VAESDEC  32(SP), X1, X1
	VAESDEC  48(SP), X2, X2
	VAESDEC  64(SP), X3, X3
	VAESDEC  80(SP), X4, X4
	VAESDEC  96(SP), X5, X5
	VAESDEC  112(SP), X6, X6
	VAESDEC  128(SP), X7, X7
	VAESDEC  144(SP), X8, X8
	VAESDEC  160(SP), X9, X9
	VAESDEC  176(SP), X10, X10
	VAESDEC  192(SP), X11, X11
	VAESDEC  208(SP), X12, X12
	VAESDEC  224(SP), X13, X13
	VAESDEC  240(SP), X14, X14
	VAESDEC  256(SP), X15, X15

finish:
	MOVOU    X0, 16(SP)
	MOVOU    X1, 32(SP)
	MOVOU    X2, 48(SP)
	MOVOU    X3, 64(SP)
	MOVOU    0(SP), X0
	MOVOU    0(SP), X1
	MOVOU    0(SP), X2
	MOVOU    0(SP), X3

	// Rotation block 0.
	VAESDEC  16(SP), X0, X0
	VAESDEC  32(SP), X1, X1
	VAESDEC  48(SP), X2, X2
	VAESDEC  64(SP), X3, X3
	VAESDEC  X4, X0, X0
	VAESDEC  X5, X1, X1
	VAESDEC  X6, X2, X2
	VAESDEC  X7, X3, X3
	VAESDEC  X8, X0, X0
	VAESDEC  X9, X1, X1
	VAESDEC  X10, X2, X2
	VAESDEC  X11, X3, X3
	VAESDEC  X12, X0, X0
	VAESDEC  X13, X1, X1
	VAESDEC  X14, X2, X2
	VAESDEC  X15, X3, X3

	// Rotation block 1.
	VAESDEC  32(SP), X0, X0
	VAESDEC  48(SP), X1, X1
	VAESDEC  64(SP), X2, X2
	VAESDEC  16(SP), X3, X3
	VAESDEC  X5, X0, X0
	VAESDEC  X6, X1, X1
	VAESDEC  X7, X2, X2
	VAESDEC  X4, X3, X3
	VAESDEC  X9, X0, X0
	VAESDEC  X10, X1, X1
	VAESDEC  X11, X2, X2
	VAESDEC  X8, X3, X3
	VAESDEC  X13, X0, X0
	VAESDEC  X14, X1, X1
	VAESDEC  X15, X2, X2
	VAESDEC  X12, X3, X3

	// Rotation block 2.
	VAESDEC  48(SP), X0, X0
	VAESDEC  64(SP), X1, X1
	VAESDEC  16(SP), X2, X2
	VAESDEC  32(SP), X3, X3
	VAESDEC  X6, X0, X0
	VAESDEC  X7, X1, X1
	VAESDEC  X4, X2, X2
	VAESDEC  X5, X3, X3
	VAESDEC  X10, X0, X0
	VAESDEC  X11, X1, X1
	VAESDEC  X8, X2, X2
	VAESDEC  X9, X3, X3
	VAESDEC  X14, X0, X0
	VAESDEC  X15, X1, X1
	VAESDEC  X12, X2, X2
	VAESDEC  X13, X3, X3

	// Rotation block 3.
	VAESDEC  64(SP), X0, X0
	VAESDEC  16(SP), X1, X1
	VAESDEC  32(SP), X2, X2
	VAESDEC  48(SP), X3, X3
	VAESDEC  X7, X0, X0
	VAESDEC  X4, X1, X1
	VAESDEC  X5, X2, X2
	VAESDEC  X6, X3, X3
	VAESDEC  X11, X0, X0
	VAESDEC  X8, X1, X1
	VAESDEC  X9, X2, X2
	VAESDEC  X10, X3, X3
	VAESDEC  X15, X0, X0
	VAESDEC  X12, X1, X1
	VAESDEC  X13, X2, X2
	VAESDEC  X14, X3, X3

	// Final merge.
	VAESDEC  0(SP), X0, X0
	VAESDEC  0(SP), X1, X1
	VAESDEC  0(SP), X2, X2
	VAESDEC  0(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  0(SP), X1, X1
	VAESDEC  0(SP), X2, X2
	VAESDEC  0(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  0(SP), X1, X1
	VAESDEC  0(SP), X2, X2
	VAESDEC  0(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  0(SP), X1, X1
	VAESDEC  0(SP), X2, X2
	VAESDEC  0(SP), X3, X3
	VAESDEC  0(SP), X0, X0
	VAESDEC  0(SP), X1, X1
	VAESDEC  0(SP), X2, X2
	VAESDEC  0(SP), X3, X3

	// Store hash.
	MOVOU    X0, 0(DST_PTR)
	MOVOU    X1, 16(DST_PTR)
	MOVOU    X2, 32(DST_PTR)
	MOVOU    X3, 48(DST_PTR)
	RET      
