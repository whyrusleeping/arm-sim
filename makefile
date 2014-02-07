all:
	go build parse.go machine.go

asm:
	gcc -O0 -S test.c
