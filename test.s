	.arch armv7-a
	.eabi_attribute 27, 3
	.eabi_attribute 28, 1
	.fpu vfpv3-d16
	.eabi_attribute 20, 1
	.eabi_attribute 21, 1
	.eabi_attribute 23, 3
	.eabi_attribute 24, 1
	.eabi_attribute 25, 1
	.eabi_attribute 26, 2
	.eabi_attribute 30, 6
	.eabi_attribute 34, 1
	.eabi_attribute 18, 4
	.file	"test.c"
	.text
	.align	2
	.global	printd
	.type	printd, %function
printd:
	@ args = 0, pretend = 0, frame = 8
	@ frame_needed = 1, uses_anonymous_args = 0
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	sub	sp, sp, #8
	str	r0, [fp, #-8]
	ldr	r3, [fp, #-8]
	cmp	r3, #0
	bne	.L2
	mov	r0, #48
	bl	putc
	b	.L1
.L2:
	ldr	r2, [fp, #-8]
	movw	r3, #26215
	movt	r3, 26214
	smull	r1, r3, r3, r2
	mov	r1, r3, asr #2
	mov	r3, r2, asr #31
	rsb	r3, r3, r1
	mov	r0, r3
	bl	printd
	ldr	r1, [fp, #-8]
	movw	r3, #26215
	movt	r3, 26214
	smull	r2, r3, r3, r1
	mov	r2, r3, asr #2
	mov	r3, r1, asr #31
	rsb	r2, r3, r2
	mov	r3, r2
	mov	r3, r3, asl #2
	add	r3, r3, r2
	mov	r3, r3, asl #1
	rsb	r2, r3, r1
	add	r3, r2, #48
	mov	r0, r3
	bl	putc
.L1:
	sub	sp, fp, #4
	@ sp needed
	ldmfd	sp!, {fp, pc}
	.size	printd, .-printd
	.align	2
	.global	sort
	.type	sort, %function
sort:
	@ args = 0, pretend = 0, frame = 16
	@ frame_needed = 1, uses_anonymous_args = 0
	@ link register save eliminated.
	str	fp, [sp, #-4]!
	add	fp, sp, #0
	sub	sp, sp, #20
	str	r0, [fp, #-16]
	str	r1, [fp, #-20]
	mov	r3, #0
	str	r3, [fp, #-8]
	mov	r3, #0
	str	r3, [fp, #-12]
	mov	r3, #0
	str	r3, [fp, #-8]
	b	.L5
.L9:
	ldr	r3, [fp, #-8]
	add	r3, r3, #1
	str	r3, [fp, #-12]
	b	.L6
.L8:
	ldr	r3, [fp, #-8]
	mov	r3, r3, asl #2
	ldr	r2, [fp, #-16]
	add	r3, r2, r3
	ldr	r2, [r3]
	ldr	r3, [fp, #-12]
	mov	r3, r3, asl #2
	ldr	r1, [fp, #-16]
	add	r3, r1, r3
	ldr	r3, [r3]
	cmp	r2, r3
	ble	.L7
	ldr	r3, [fp, #-8]
	mov	r3, r3, asl #2
	ldr	r2, [fp, #-16]
	add	r3, r2, r3
	ldr	r2, [fp, #-8]
	mov	r2, r2, asl #2
	ldr	r1, [fp, #-16]
	add	r2, r1, r2
	ldr	r1, [r2]
	ldr	r2, [fp, #-12]
	mov	r2, r2, asl #2
	ldr	r0, [fp, #-16]
	add	r2, r0, r2
	ldr	r2, [r2]
	eor	r2, r1, r2
	str	r2, [r3]
	ldr	r3, [fp, #-12]
	mov	r3, r3, asl #2
	ldr	r2, [fp, #-16]
	add	r3, r2, r3
	ldr	r2, [fp, #-12]
	mov	r2, r2, asl #2
	ldr	r1, [fp, #-16]
	add	r2, r1, r2
	ldr	r1, [r2]
	ldr	r2, [fp, #-8]
	mov	r2, r2, asl #2
	ldr	r0, [fp, #-16]
	add	r2, r0, r2
	ldr	r2, [r2]
	eor	r2, r1, r2
	str	r2, [r3]
	ldr	r3, [fp, #-8]
	mov	r3, r3, asl #2
	ldr	r2, [fp, #-16]
	add	r3, r2, r3
	ldr	r2, [fp, #-8]
	mov	r2, r2, asl #2
	ldr	r1, [fp, #-16]
	add	r2, r1, r2
	ldr	r1, [r2]
	ldr	r2, [fp, #-12]
	mov	r2, r2, asl #2
	ldr	r0, [fp, #-16]
	add	r2, r0, r2
	ldr	r2, [r2]
	eor	r2, r1, r2
	str	r2, [r3]
.L7:
	ldr	r3, [fp, #-12]
	add	r3, r3, #1
	str	r3, [fp, #-12]
.L6:
	ldr	r2, [fp, #-12]
	ldr	r3, [fp, #-20]
	cmp	r2, r3
	blt	.L8
	ldr	r3, [fp, #-8]
	add	r3, r3, #1
	str	r3, [fp, #-8]
.L5:
	ldr	r2, [fp, #-8]
	ldr	r3, [fp, #-20]
	cmp	r2, r3
	blt	.L9
	sub	sp, fp, #0
	@ sp needed
	ldr	fp, [sp], #4
	bx	lr
	.size	sort, .-sort
	.section	.rodata
	.align	2
.LC0:
	.word	7
	.word	2
	.word	9
	.word	3
	.word	4
	.word	12
	.word	2
	.word	15
	.word	3
	.word	21
	.text
	.align	2
	.global	main
	.type	main, %function
main:
	@ args = 0, pretend = 0, frame = 48
	@ frame_needed = 1, uses_anonymous_args = 0
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	sub	sp, sp, #48
	movw	r3, #:lower16:.LC0
	movt	r3, #:upper16:.LC0
	sub	ip, fp, #48
	mov	lr, r3
	ldmia	lr!, {r0, r1, r2, r3}
	stmia	ip!, {r0, r1, r2, r3}
	ldmia	lr!, {r0, r1, r2, r3}
	stmia	ip!, {r0, r1, r2, r3}
	ldmia	lr, {r0, r1}
	stmia	ip, {r0, r1}
	sub	r3, fp, #48
	mov	r0, r3
	mov	r1, #10
	bl	sort
	mov	r3, #0
	str	r3, [fp, #-8]
	b	.L11
.L12:
	ldr	r2, [fp, #-8]
	mvn	r3, #43
	mov	r2, r2, asl #2
	sub	r1, fp, #4
	add	r2, r1, r2
	add	r3, r2, r3
	ldr	r3, [r3]
	mov	r0, r3
	bl	printd
	ldr	r3, [fp, #-8]
	add	r3, r3, #1
	str	r3, [fp, #-8]
.L11:
	ldr	r3, [fp, #-8]
	cmp	r3, #9
	ble	.L12
	sub	sp, fp, #4
	@ sp needed
	ldmfd	sp!, {fp, pc}
	.size	main, .-main
	.ident	"GCC: (GNU) 4.8.2 20131219 (prerelease)"
	.section	.note.GNU-stack,"",%progbits
