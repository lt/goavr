package main

import (
	"fmt"
)

func dissAssemble(b []byte) Instr {
	dm := map[byte]byte{
		0: 24,
		1: 26,
		2: 28,
		3: 30,
	}
	
	fmt.Printf("pc: %.4x\t", (cpu.pc))
	m := lookUp(b)
	inst := Instr{family: m.family, label: m.label}
	switch m.label {
	case INSN_NOP:
		fmt.Printf("%.4x\tnop\n", b2u16big(b))
		return inst
	case INSN_CLI:
		fmt.Printf("%.4x\tcli\n", b2u16big(b))
		return inst
	case INSN_ADC:
		// 0001 11rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_EOR:
		// 0010 01rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\teor\tr%d, r%d\n", b2u16big(b), inst.source, inst.dest)
		return inst
	case INSN_OUT:
		// 1011 1AAr rrrr AAAA
		//out := (b[1] >> 3) & 0xff
		inst.ioaddr = ((b[1] & 0x06) << 3) | (b[0] & 0x0f)
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tout\t0x%.2x, r%d\t\t;%d\n", b2u16big(b), inst.ioaddr, inst.source, inst.ioaddr)
		return inst
	case INSN_RJMP:
		// 1100 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		if ((k & 0x800) >> 11) == 1 {
			inst.k16 = int16((k + 0xf000) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\trjmp\t.+%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_LDI:
		// 1110 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tldi\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_RCALL:
		// 1101 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		if ((k & 0x0800) >> 11) == 1 {
			inst.k16 = int16((k + 0xf000) << 1)
			fmt.Printf("%.4x\trcall\t.%d\t;%.4x\n", b2u16big(b), inst.k16, (inst.k16 + cpu.pc))
		} else {
			inst.k16 = int16(k << 1)
			fmt.Printf("%.4x\trcall\t.+%d\t;%.4x\n", b2u16big(b), inst.k16, (inst.k16 + cpu.pc))
		}
		return inst
	case INSN_SBI:
		// 1001 1010 AAAA Abbb
		inst.ioaddr = b[0] >> 3
		inst.registerBit = b[0] & 0x7
		fmt.Printf("%.4x\tsbi\t0x%x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_CBI:
		//1001 1000 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tcbi\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_SBIC:
		// 1001 1001 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tsbic\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_SBIS:
		// 1001 1011 AAAA Abbb
		inst.ioaddr = (b[0] & 0xf8) >> 3
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tsbis\t0x%.2x, %d\n", b2u16big(b), inst.ioaddr, inst.registerBit)
		return inst
	case INSN_BLD:
		// 1111 100d dddd 0bbb
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tbld\tr%d, %d\n", b2u16big(b), inst.dest, inst.registerBit)
		return inst
	case INSN_BST:
		// 1111 101d dddd 0bbb
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tbst\tr%d, %d\n", b2u16big(b), inst.dest, inst.registerBit)
		return inst
	case INSN_SBRC:
		// 1111 110r rrrr 0bbb
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.registerBit = b[0] & 0x07
		fmt.Printf("%.4x\tsbrc\tr%d, %d\n", b2u16big(b), inst.source, inst.registerBit)
		return inst
	case INSN_STS:
		// 1001 001d dddd 0000 kkkk kkkk kkkk kkkk
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		inst.k16 = b2i16little(c)
		fmt.Printf("%.4x\tsts\t0x%.4x, r%d\n", b2u16big(b), inst.k16, inst.dest)
		return inst
	case INSN_LDS:
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := cpu.imem.Fetch()
		inst.k16 = b2i16little(c)
		fmt.Printf("%.4x\tlds\tr%d, 0x%.4x\n", b2u16big(b), inst.dest, inst.k16)
		return inst
	case INSN_JMP:
		// 1001 010k kkkk 110k kkkk kkkk kkkk kkkk
		var k1, k2, k3 uint32
		k1 = uint32(b[1] & 0x01)<< 20
		k2 = uint32(b[0] & 0xf0)<< 12
		c := cpu.imem.Fetch()
		k3 = uint32(c[1]) << 8 | uint32(c[0])
		inst.k32 = k1 | k2 | k3
		fmt.Printf("%.4x\tjmp\t0x%.8x\t;%d\n", b2u16big(b), inst.k32, inst.k32)
		return inst
	case INSN_ADD:
		// 0000 11rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadd\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_ADIW:
		// 1001 0110 KKdd KKKK
		// 24,26,28,30
		inst.kdata = ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		inst.dest = dm[Rd]
		fmt.Printf("%.4x\tadiw\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_SBIW:
		// 1001 0111 KKdd KKKK
		inst.kdata = ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		inst.dest = dm[Rd]
		fmt.Printf("%.4x\tsbiw\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_ANDI:
		//0111 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tandi\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_BRCC:
		//1111 01kk kkkk k000
		// 64 ≤ k ≤ +63
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrcc\t.+%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_BRCS:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k000
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrcs\t.+%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_BREQ:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\tbreq\t.%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_BRGE:
		// 1111 01kk kkkk k100
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrge\t.%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_BRNE:
		// 1111 01kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		// check to see if msb of k is 1
		// if it is, the result is negative.
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		fmt.Printf("%.4x\tbrne\t.%d\n", b2u16big(b), inst.k16)
		} else {
			inst.k16 = int16(k << 1)
		fmt.Printf("%.4x\tbrne\t.+%d\n", b2u16big(b), inst.k16)
		}
		return inst
	case INSN_BRTC:
		// 1111 01kk kkkk k110
		k := (b2u16little(b) & 0x03f8) >> 3
		if ((k & 0x40) >> 6) == 1 {
			inst.k16 = int16((k + 0xff80) << 1)
		} else {
			inst.k16 = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrtc\t.+%d\n", b2u16big(b), inst.k16)
		return inst
	case INSN_COM:
		// 1001 010d dddd 0000
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcom\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_CP:
		// 0001 01rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcp\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_CPC:
		// 0000 01rd dddd rrrr
		inst.source = ((b[1] & 0x02) << 3) | b[0]&0x0f
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcpc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_CPI:
		// 0011 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tcpi\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_CPSE:
		// 0001 00rd dddd rrrr
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcpse\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_DEC:
		// 1001 010d dddd 1010
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tdec\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_IN:
		// 1011 0AAd dddd AAAA
		inst.ioaddr = ((b[1] & 0x09) << 3) | (b[0] & 0x0f)
		inst.dest = ((b[1] & 0xf1) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tin\tr%d, 0x%.2x\n", b2u16big(b), inst.dest, inst.ioaddr)
		return inst
	case INSN_LDDY:
		// 1001 000d dddd 1001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		fmt.Printf("%.4x\tldd\tY+%d, r%d\n", b2u16big(b), inst.offset, inst.dest)
		return inst
	case INSN_LDDZ:
		// 1001 000d dddd 0001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		inst.offset = m.offset
		fmt.Printf("%.4x\tldd\tr%d, Z+%d\n", b2u16big(b), inst.dest, inst.offset)
		return inst
	case INSN_LDX:
		// 1001 000d dddd 1100
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDXP:
		//1001 000d dddd 1101
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDY:
		// 1001 000d dddd 1001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Y\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LDZ:
		// 1000 000d dddd 0000
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Z\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPMZ:
		//z  1001 000d dddd 0100
		// XXX ToDo not tested
		// XXX ToDo(Erin) Not sure this works.
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tlpm\tr%d, Z\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPMZP:
		//z+ 1001 000d dddd 0101
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tlpm\tr%d, Z+\n", b2u16big(b), inst.dest)
		return inst
	case INSN_LPM:
		// 1001 0101 1100 1000
		// XXX ToDo: not tested
		// i := Instr{label: m.label, family: m.family, source: 0x00}
		fmt.Printf("%.4x\tlpm\n", b2u16big(b))
		return inst
	case INSN_LSR:
		//1001 010d dddd 0110
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tlsr\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_MOV:
		// 0010 11rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmov\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_MOVW:
		// 0000 0001 dddd rrrr
		inst.dest = (b[0] & 0xf0) >> 3
		inst.source = (b[0] & 0x0f) << 1
		fmt.Printf("%.4x\tmovw\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_MUL:
		// 1001 11rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmul\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_NEG:
		// 1001 010d dddd s0001
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tneg\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_OR:
		// 0010 10rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tor\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_ORI:
		// 0110 KKKK dddd KKKK
		inst.kdata = ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tori\tr%d, %.2x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_POP:
		// 1001 000d dddd 1111
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tpop\t r%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_PUSH:
		//1001 001d dddd 1111
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tpush\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_RET:
		fmt.Printf("%.4x\tret\n", b2u16big(b))
		return inst
	case INSN_RETI:
		fmt.Printf("%.4x\treti\n", b2u16big(b))
		return inst
	case INSN_ROR:
		// 1001 010d dddd 0111
		inst.dest = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tror\tr%d\n", b2u16big(b), inst.dest)
		return inst
	case INSN_SBC:
		// 0000 10rd rrrr dddd
		inst.source = (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsbc\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_SBCI:
		// 0100 KKKK dddd KKKK
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.kdata = (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsbci\tr%d, 0x%x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	case INSN_SEI:
		// 1001 0100 0111 1000
		fmt.Printf("%.4x\tsei\n", b2u16big(b))
		return inst
	case INSN_STDY:
		// 1001 001r rrrr 1001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tY+%d, r%d\n", b2u16big(b), inst.offset, inst.source)
		return inst
	case INSN_STDZ:
		// 10q0 qq1r rrrr 0qqq
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tZ+%d, r%d\n", b2u16big(b), inst.offset, inst.source)
		return inst
	case INSN_STX:
		// 1001 001r rrrr 1100
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tX, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STXP:
		// 1001 001r rrrr 1101
		//inst.source = (b2u16little(b) & 0x01f0) >> 4
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tX+, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STXM:
		// 1001 001r rrrr 1110
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-X, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STY:
		//1001 001r rrrr 1001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tY, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZ:
		// 1000 001r rrrr 0000
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZP:
		// 1001 001r rrrr 0001
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ+, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_STZM:
		// 1001 001r rrrr 0010
		inst.source = ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-Z, r%d\n", b2u16big(b), inst.source)
		return inst
	case INSN_SUB:
		// 0001 10rd dddd rrrr
		inst.source = (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		inst.dest = ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsub\tr%d, r%d\n", b2u16big(b), inst.dest, inst.source)
		return inst
	case INSN_SUBI:
		// 0101 KKKK dddd KKKK
		inst.dest = ((b[0] & 0xf0) >> 4) + 0x10
		inst.kdata = (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsubi\tr%d, 0x%x\n", b2u16big(b), inst.dest, inst.kdata)
		return inst
	default:
		fmt.Printf("None of the above. Got %s (0x%.4x)\n", m.mnemonic, b2u16big(b))
		return inst
	}

}

func lookUp(raw []byte) OpCode {
	var op OpCode
	b := b2u16little(raw)
	for _, entry := range OpCodeLookUpTable {
		v := b & entry.mask
		if v == entry.value {
			op = entry
			switch entry.mnemonic {
			case "std":
				return deConvoluter(b, op)
			case "ldd":
				return deConvoluter(b, op)
			}
			return op
		} else {
			op = OpCode{mnemonic: "unknown", value: b}
		}
	}
	return op
}

func deConvoluter(b uint16, op OpCode) OpCode {
	x := b & 0xd208
	offset := b & 0x2c07
	switch x {
	case 0x8000:
		if offset == 0 {
			op.mnemonic = "ldz"
			op.label = INSN_LDZ
		} else {
			op.mnemonic = "lddz"
			op.offset = offset
			op.label = INSN_LDDZ
		}
	case 0x8008:
		if offset == 0 {
			op.mnemonic = "ldy"
			op.label = INSN_LDY
		} else {
			op.mnemonic = "lddy"
			op.offset = offset
			op.label = INSN_LDDY
		}
	case 0x8200:
		if offset == 0 {
			op.mnemonic = "stz"
			op.label = INSN_STZ
		} else {
			op.mnemonic = "stdz"
			op.offset = offset
			op.label = INSN_STDZ
		}
	case 0x8208:
		if offset == 0 {
			op.mnemonic = "sty"
			op.label = INSN_STY
		} else {
			op.mnemonic = "stdy"
			op.offset = offset
			op.label = INSN_STDY
		}
	default:
		op.mnemonic = "Unknown"
		op.value = b
	}
	return op
}
