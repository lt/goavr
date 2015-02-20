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
	OpCode{
		Name: "cpc",
		Mask: 0xfc00,
		Value: 0x0400,
	},
	OpCode{
		Name: "push",
		Mask: 0xfe0f,
		Value: 0x920f,
	},
	OpCode{
		Name: "subi",
		Mask: 0xf000,
		Value: 0x5000,
	},
	OpCode{
		Name: "sbci",
		Mask: 0xf000,
		Value: 0x4000,
	},
	OpCode{
		Name: "ori",
		Mask: 0xf000,
		Value: 0x6000,
	},
	OpCode{
		Name: "sbc",
		Mask: 0xfc00,
		Value: 0x0800,
	},
	OpCode{
		Name: "lsr",
		Mask: 0xfe0f,
		Value: 0x9406,
	},
	OpCode{
		Name: "mov",
		Mask: 0xfc00,
		Value: 0x2c00,
	},
	OpCode{
		Name: "ror",
		Mask: 0xfe0f,
		Value: 0x9407,
	},
	OpCode{
		Name: "or",
		Mask: 0xfc00,
		Value: 0x2800,
	},
	OpCode{
		Name: "cbi",
		Mask: 0xff00,
		Value: 0x9800,
	},
	OpCode{
		Name: "pop",
		Mask: 0xfe0f,
		Value: 0x900f,
	},
	OpCode{
		Name: "reti",
		Mask: 0xffff,
		Value: 0x9518,
	},
	OpCode{
		Name: "ret",
		Mask: 0xffff,
		Value: 0x9508,
	},
	OpCode{
		Name: "sbis",
		Mask: 0xff00,
		Value: 0x9b00,
	},
	OpCode{
		Name: "cpse",
		Mask: 0xfc00,
		Value: 0x1000,
	},
	OpCode{
		Name: "brcs",
		Mask: 0xfc07,
		Value: 0xf000,
	},
	OpCode{
		Name: "sbiw",
		Mask: 0xff00,
		Value: 0x9700,
	},
	OpCode{
		Name: "brcc",
		Mask: 0xfc07,
		Value: 0xf400,
	},
	OpCode{
		Name: "cp",
		Mask: 0xfc00,
		Value: 0x1400,
	},
	OpCode{
		Name: "adiw",
		Mask: 0xff00,
		Value: 0x9600,
	},
	OpCode{
		Name: "andi",
		Mask: 0xf000,
		Value: 0x7000,
	},
	OpCode{
		Name: "sbic",
		Mask: 0xfd00,
		Value: 0x9900,
	},
	OpCode{
		Name: "bst",
		Mask: 0xfe08,
		Value: 0xfa00,
	},
	OpCode{
		Name: "bld",
		Mask: 0xfe08,
		Value: 0xf800,
	},
	OpCode{
		Name: "sei",
		Mask: 0xffff,
		Value: 0x9478,
	},
	OpCode{
		Name: "brge",
		Mask: 0xfc07,
		Value: 0xf404,
	},
	OpCode{
		Name: "brtc",
		Mask: 0xfc07,
		Value: 0xf406,
	},
	OpCode{
		Name: "com",
		Mask: 0xfe0f,
		Value: 0x9400,
	},
	OpCode{
		Name: "sbrc",
		Mask: 0xfe08,
		Value: 0xfc00,
	},
	OpCode{
		Name: "sbiw",
		Mask: 0xff00,
		Value: 0x9700,
	},
	OpCode{
		Name: "sbr",
		Mask: 0xfe08,
		Value: 0x6000,
		},
	OpCode{
		Name: "neg",
		Mask: 0xfe0f,
		Value: 0x9401,
	},
	OpCode{
		Name: "sub",
		Mask: 0xfc00,
		Value: 0x1800,
	},
	OpCode{
		Name: "dec",
		Mask: 0xfe0f,
		Value: 0x940a,
	},
	OpCode{
		Name: "mul",
		Mask: 0xfc00,
		Value: 0x9c00,
	},
	// =======
	// Things that work with registers
	// This are tricky. the q values are interpolated into the other bits.
	// But applying the same mask as the other LD w/Z ops gives 0x8001.
	// I'm going to leave it this way until another opcode with that value comes
	// along (and hope that it won't).
	// XXX TODO: This screws up the actual mask value for this opcode. Which
	// Likely means that the way I'm doing all of these opcodes is wrong. Yay.
	// ========
	// LD Rd, X
	OpCode{
		Name: "ldx",
		Mask: 0xfe0f,
		Value: 0x900c,
	},
	// LD Rd, X+
	OpCode{
		Name: "ldx+",
		Mask: 0xfe0f,
		Value: 0x900d,
	},
	// LD Rd, -X
	OpCode{
		Name: "ldx-",
		Mask: 0xfe0f,
		Value: 0x900e,
	},
	// LD Rd, Z
	// LDD Rd, Y+q
	// LDD Rd, Z+q
	OpCode{
		Name: "ldd",
		Mask: 0xde00,
		Value: 0x8000,
	},
	// ST X, Rr
	OpCode{
		Name: "stx",
		Mask: 0xfe0f,
		Value: 0x920c,
	},	
	// ST X+, Rr
	OpCode{
		Name: "stx+",
		Mask: 0xfe0f,
		Value: 0x920d,
	},
	// ST -X, Rr
	OpCode{
		Name: "stx-",
		Mask: 0xfe0f,
		Value: 0x920e,
	},
	// ST Z, Rr
	// STD Y+q, Rr
	// STD Z+q, Rr
	OpCode{
		Name: "std",
		Mask: 0xde00,
		Value: 0x8200,
	},
	// ST Y+, Rr
	OpCode{
		Name: "sty+",
		Mask: 0xfe0f,
		Value: 0x9209,
	},
	// ST -Y, Rr
	OpCode{
		Name: "sty-",
		Mask: 0xfe0f,
		Value: 0x920a,
	},
	// ST Z+, Rr
	OpCode{
		Name: "stz+",
		Mask: 0xfe0f,
		Value: 0x9201,
	},
	// ST -Z, Rr
		OpCode{
		Name: "stz-",
		Mask: 0xfe0f,
		Value: 0x9202,
		},
	// LPM Rd, Z+
	OpCode{
		Name: "lpmz+",
		Mask: 0xfeff,
		Value: 0x9005,
	},
	// =======
	// END things that work with registers
	// ======
	// 32 bit opcodes:
	OpCode{
		Name: "lds",
		Mask: 0xfe0f,
		Value: 0x9000,
	},
	OpCode{
		Name: "sts",
		Mask: 0xfe0f,
		Value: 0x9200,
	},
}
