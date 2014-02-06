main:
	stmfd	sp!, {fp, lr}
	add	fp, sp, #4
	sub	sp, sp, #8
	mov	r3, #5
	str	r3, [fp, #-8]
	mov	r3, #2
	str	r3, [fp, #-12]
	mov	r0, #104
	bl	putc
	mov	r0, #101
	bl	putc
	mov	r0, #108
	bl	putc
	mov	r0, #108
	bl	putc
	mov	r0, #111
	bl	putc
	mov	r0, #10
	bl	putc
	ldr	r3, [fp, #-8]
	add	r2, r3, #48
	ldr	r3, [fp, #-12]
	add	r3, r2, r3
	mov	r0, r3
	bl	putc
	mov	r0, #10
	bl	putc
	sub	sp, fp, #4
	@ sp needed
	ldmfd	sp!, {fp, pc}
	.size	main, .-main
	.ident	"GCC: (GNU) 4.8.2 20131219 (prerelease)"
	.section	.note.GNU-stack,"",%progbits
