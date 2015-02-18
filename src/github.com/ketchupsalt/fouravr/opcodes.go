package main

// Bloody heck.

const (
	EOR = 0x20
	ADD = 0x30
	ADC = 0x70
	AND = 0x80
	IN = 0xb0
	OUT = 0xb8
	RJMP = 0xc0
	RCALL = 0xd0
	LDI = 0xe0
	CLI = 0x90
	SBI = 0x9a
/* not used yet
	ADIW = 0x96
	ASR = 0x74
	BCLR = 0x148
	PUSH = 0x73
	POP = 0x72
*/
)


type OpCode struct {
	Value uint16
	Mask uint16
	Op int
	Name string
}


type HookUp struct {
	Set []OpCode
}


// Look up the mask, return a parsed opcode
func parseOpCode(b []byte) OpCode {
	o := OpCode{}
	return o
}

// Put all your changes above this line cuz it's hairy
// down here.

var OpCodeLookUpTable = []OpCode{
	OpCode{
		Name: "nop",
		Mask: 0xffff,
		Value: 0x0000,
	},
	OpCode{
		Name: "movw",
		Mask: 0xff00,
		Value: 0x0100,
	},
	OpCode{
		Name: "eor",
		Mask: 0xfc00,
		Value: 0x2400,
	},
	OpCode{
		Name: "add",
		Mask: 0xfc00,
		Value: 0x0c00,
	},
	OpCode{
		Name: "adc",
		Mask: 0xfc00,
		Value: 0x1c00,
	},
	OpCode{
		Name: "and",
		Mask: 0xfc00,
		Value: 0x2000,
	},
	OpCode{
		Name: "in",
		Mask: 0xf800,
		Value: 0xb000,
	},
	OpCode{
		Name: "out",
		Mask: 0xf800,
		Value: 0xb800,
	},
	OpCode{
		Name: "rjmp",
		Mask: 0xf000,
		Value: 0xc000,
	},
	OpCode{
		Name: "rcall",
		Mask: 0xf000,
		Value: 0xd000,
	},
	OpCode{
		Name: "ldi",
		Mask: 0xf000,
		Value: 0xe000,
	},
	OpCode{
		Name: "cli",
		Mask: 0xffff,
		Value: 0x94f8,
	},
	OpCode{
		Name: "sbi",
		Mask: 0xff00,
		Value: 0x9a00,
	},
	OpCode{
		Name: "cpi",
		Mask: 0xf000,
		Value: 0x3000,
	},
	OpCode{
		Name: "breq",
		Mask: 0xfc07,
		Value: 0xf001,
	},
	OpCode{
		Name: "subi",
		Mask: 0xf000,
		Value: 0x5000,
	},
	OpCode{
		Name: "brne",
		Mask: 0xfc07,
		Value: 0xf401,
	},
}
