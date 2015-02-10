package main

// Bloody heck.

const (
	EOR = 0x20
	ADD = 0x30
	ADC = 0x70
	AND = 0x80
	OUT = 0xb0
	RJMP = 0xc0
	RCALL = 0xd0
	LDI = 0xe0
/* not used yet
	ADIW = 0x96
	ASR = 0x74
	BCLR = 0x148
	PUSH = 0x73
	POP = 0x72
*/
)

/*
type OpCode struct {
	Code int64
	Instruction string
}

type OpCodes struct {
	Foo []OpCode
}
*/
