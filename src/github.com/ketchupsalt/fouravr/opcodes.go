package main

type family int

const (
	Arithmetic family = 0
	Branches = 1
	Transfers = 2
	BitWise = 3
)

type OpCode struct {
	value uint16
	mask uint16
	mnemonic string
	offset uint16
	family family
}

type Instr struct {
	family family
	mnemonic string
	offset uint16
	dest byte
	source byte
	result byte
	kdata byte
	kaddress byte
	ioaddr byte
	iar map[string]uint16
	displacement uint16
	registerBit uint16
	statusBit uint16
}


// Return a parsed opcode from a byte. Not used (yet)
func parseOpCode(b []byte) OpCode {
	o := OpCode{}
	return o
}

// Put all your changes above this line cuz it's hairy
// down here.

var OpCodeLookUpTable = []OpCode{
	OpCode{
		mnemonic: "nop",
		mask: 0xffff,
		value: 0x0000,
		family: BitWise,
	},
	OpCode{
		mnemonic: "movw",
		mask: 0xff00,
		value: 0x0100,
	},
	OpCode{
		mnemonic: "eor",
		mask: 0xfc00,
		value: 0x2400,
		family: Arithmetic, 
	},
	OpCode{
		mnemonic: "add",
		mask: 0xfc00,
		value: 0x0c00,
	},
	OpCode{
		mnemonic: "adc",
		mask: 0xfc00,
		value: 0x1c00,
	},
	OpCode{
		mnemonic: "and",
		mask: 0xfc00,
		value: 0x2000,
	},
	OpCode{
		mnemonic: "in",
		mask: 0xf800,
		value: 0xb000,
	},
	OpCode{
		mnemonic: "out",
		mask: 0xf800,
		value: 0xb800,
	},
	OpCode{
		mnemonic: "rjmp",
		mask: 0xf000,
		value: 0xc000,
	},
	OpCode{
		mnemonic: "rcall",
		mask: 0xf000,
		value: 0xd000,
	},
	OpCode{
		mnemonic: "ldi",
		mask: 0xf000,
		value: 0xe000,
	},
	OpCode{
		mnemonic: "cli",
		mask: 0xffff,
		value: 0x94f8,
	},
	OpCode{
		mnemonic: "sbi",
		mask: 0xff00,
		value: 0x9a00,
	},
	OpCode{
		mnemonic: "cpi",
		mask: 0xf000,
		value: 0x3000,
	},
	OpCode{
		mnemonic: "breq",
		mask: 0xfc07,
		value: 0xf001,
	},
	OpCode{
		mnemonic: "subi",
		mask: 0xf000,
		value: 0x5000,
	},
	OpCode{
		mnemonic: "brne",
		mask: 0xfc07,
		value: 0xf401,
	},
	OpCode{
		mnemonic: "cpc",
		mask: 0xfc00,
		value: 0x0400,
	},
	OpCode{
		mnemonic: "push",
		mask: 0xfe0f,
		value: 0x920f,
	},
	OpCode{
		mnemonic: "subi",
		mask: 0xf000,
		value: 0x5000,
	},
	OpCode{
		mnemonic: "sbci",
		mask: 0xf000,
		value: 0x4000,
	},
	OpCode{
		mnemonic: "ori",
		mask: 0xf000,
		value: 0x6000,
	},
	OpCode{
		mnemonic: "sbc",
		mask: 0xfc00,
		value: 0x0800,
	},
	OpCode{
		mnemonic: "lsr",
		mask: 0xfe0f,
		value: 0x9406,
	},
	OpCode{
		mnemonic: "mov",
		mask: 0xfc00,
		value: 0x2c00,
	},
	OpCode{
		mnemonic: "ror",
		mask: 0xfe0f,
		value: 0x9407,
	},
	OpCode{
		mnemonic: "or",
		mask: 0xfc00,
		value: 0x2800,
	},
	OpCode{
		mnemonic: "cbi",
		mask: 0xff00,
		value: 0x9800,
	},
	OpCode{
		mnemonic: "pop",
		mask: 0xfe0f,
		value: 0x900f,
	},
	OpCode{
		mnemonic: "reti",
		mask: 0xffff,
		value: 0x9518,
	},
	OpCode{
		mnemonic: "ret",
		mask: 0xffff,
		value: 0x9508,
	},
	OpCode{
		mnemonic: "sbis",
		mask: 0xff00,
		value: 0x9b00,
	},
	OpCode{
		mnemonic: "cpse",
		mask: 0xfc00,
		value: 0x1000,
	},
	OpCode{
		mnemonic: "brcs",
		mask: 0xfc07,
		value: 0xf000,
	},
	OpCode{
		mnemonic: "sbiw",
		mask: 0xff00,
		value: 0x9700,
	},
	OpCode{
		mnemonic: "brcc",
		mask: 0xfc07,
		value: 0xf400,
	},
	OpCode{
		mnemonic: "cp",
		mask: 0xfc00,
		value: 0x1400,
	},
	OpCode{
		mnemonic: "adiw",
		mask: 0xff00,
		value: 0x9600,
	},
	OpCode{
		mnemonic: "andi",
		mask: 0xf000,
		value: 0x7000,
	},
	OpCode{
		mnemonic: "sbic",
		mask: 0xfd00,
		value: 0x9900,
	},
	OpCode{
		mnemonic: "bst",
		mask: 0xfe08,
		value: 0xfa00,
	},
	OpCode{
		mnemonic: "bld",
		mask: 0xfe08,
		value: 0xf800,
	},
	OpCode{
		mnemonic: "sei",
		mask: 0xffff,
		value: 0x9478,
	},
	OpCode{
		mnemonic: "brge",
		mask: 0xfc07,
		value: 0xf404,
	},
	OpCode{
		mnemonic: "brtc",
		mask: 0xfc07,
		value: 0xf406,
	},
	OpCode{
		mnemonic: "com",
		mask: 0xfe0f,
		value: 0x9400,
	},
	OpCode{
		mnemonic: "sbrc",
		mask: 0xfe08,
		value: 0xfc00,
	},
	OpCode{
		mnemonic: "sbiw",
		mask: 0xff00,
		value: 0x9700,
	},
	OpCode{
		mnemonic: "sbr",
		mask: 0xfe08,
		value: 0x6000,
		},
	OpCode{
		mnemonic: "neg",
		mask: 0xfe0f,
		value: 0x9401,
	},
	OpCode{
		mnemonic: "sub",
		mask: 0xfc00,
		value: 0x1800,
	},
	OpCode{
		mnemonic: "dec",
		mask: 0xfe0f,
		value: 0x940a,
	},
	OpCode{
		mnemonic: "mul",
		mask: 0xfc00,
		value: 0x9c00,
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
		mnemonic: "ldx",
		mask: 0xfe0f,
		value: 0x900c,
	},
	// LD Rd, X+
	OpCode{
		mnemonic: "ldx+",
		mask: 0xfe0f,
		value: 0x900d,
	},
	// LD Rd, -X
	OpCode{
		mnemonic: "ldx-",
		mask: 0xfe0f,
		value: 0x900e,
	},
	// LD Rd, Z
	// LDD Rd, Y+q
	// LDD Rd, Z+q
	OpCode{
		mnemonic: "ldd",
		mask: 0xde00,
		value: 0x8000,
	},
	// ST X, Rr
	OpCode{
		mnemonic: "stx",
		mask: 0xfe0f,
		value: 0x920c,
	},	
	// ST X+, Rr
	OpCode{
		mnemonic: "stx+",
		mask: 0xfe0f,
		value: 0x920d,
	},
	// ST -X, Rr
	OpCode{
		mnemonic: "stx-",
		mask: 0xfe0f,
		value: 0x920e,
	},
	// ST Z, Rr
	// STD Y+q, Rr
	// STD Z+q, Rr
	OpCode{
		mnemonic: "std",
		mask: 0xde00,
		value: 0x8200,
	},
	// ST Y+, Rr
	OpCode{
		mnemonic: "sty+",
		mask: 0xfe0f,
		value: 0x9209,
	},
	// ST -Y, Rr
	OpCode{
		mnemonic: "sty-",
		mask: 0xfe0f,
		value: 0x920a,
	},
	// ST Z+, Rr
	OpCode{
		mnemonic: "stz+",
		mask: 0xfe0f,
		value: 0x9201,
	},
	// ST -Z, Rr
		OpCode{
		mnemonic: "stz-",
		mask: 0xfe0f,
		value: 0x9202,
		},
	// LPM Rd, Z+
	// LPM Rd, Z
	OpCode{
		mnemonic: "lpmz",
		mask: 0xfeff,
		value: 0x9005,
	},
	OpCode{
		mnemonic: "lpm",
		mask: 0xffff,
		value: 0x95c8,
	},
	// =======
	// END things that work with registers
	// ======
	// 32 bit opcodes:
	OpCode{
		mnemonic: "lds",
		mask: 0xfe0f,
		value: 0x9000,
	},
	OpCode{
		mnemonic: "sts",
		mask: 0xfe0f,
		value: 0x9200,
	},
}
