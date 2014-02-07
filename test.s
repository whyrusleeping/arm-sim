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
	.global	prints
	.type	prints, %function
prints:
	@ args = 0, pretend = 0, frame = 8
	@ frame_needed = 1, uses_anonymous_args = 0
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	sub	sp, sp, #8
	str	r0, [fp, #-8]
	b	.L2
.L3:
	ldr	r3, [fp, #-8]
	add	r2, r3, #1
	str	r2, [fp, #-8]
	ldrb	r3, [r3]	@ zero_extendqisi2
	mov	r0, r3
	bl	putc
.L2:
	ldr	r3, [fp, #-8]
	ldrb	r3, [r3]	@ zero_extendqisi2
	cmp	r3, #0
	bne	.L3
	sub	sp, fp, #4
	@ sp needed
	ldmfd	sp!, {fp, pc}
	.size	prints, .-prints
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
	bne	.L5
	mov	r0, #48
	bl	putc
	b	.L4
.L5:
	ldr	r2, [fp, #-8]
	movw	r3, #26215
	movt	r3, 26214
	smull	r1, r3, r3, r2
	mov	r1, r3, asr #2
	mov	r3, r2, asr #31
	rsb	r3, r3, r1
	mov	r0, r3
	bl	printd
.L4:
	sub	sp, fp, #4
	@ sp needed
	ldmfd	sp!, {fp, pc}
	.size	printd, .-printd
	.section	.rodata
	.align	2
.LC0:
	.ascii	"Hello World!\012\000"
	.text
	.align	2
	.global	main
	.type	main, %function
main:
	@ args = 0, pretend = 0, frame = 0
	@ frame_needed = 1, uses_anonymous_args = 0
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	movw	r0, #:lower16:.LC0
	movt	r0, #:upper16:.LC0
	bl	prints
	movw	r0, #39784
	movt	r0, 1
	bl	printd
	ldmfd	sp!, {fp, pc}
	.size	main, .-main
	.ident	"GCC: (GNU) 4.8.2 20131219 (prerelease)"
	.section	.note.GNU-stack,"",%progbits
