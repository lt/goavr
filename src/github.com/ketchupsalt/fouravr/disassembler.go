package main

import (
	"fmt"
)

func dissAssemble(b []byte) {
	m := lookUp(b)

	switch m.label {

	case INSN_NOP:
		fmt.Printf("%.4x\tnop\n", b2u16big(b))
		//i := Instr{label: m.label, family: m.family}
		//fmt.Println(i)
	case INSN_CLI:
		fmt.Printf("%.4x\tcli\n", b2u16big(b))
		//i := Instr{label: m.label, family: m.family}
		//fmt.Println(i)
	case INSN_ADC:
		// 0001 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tadc\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_EOR:
		// 0010 01rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, source: Rr, dest: Rd, family: m.family}
		fmt.Printf("%.4x\teor\tr%d, r%d\n", b2u16big(b), i.source, i.dest)
	case INSN_OUT:
		// 1011 1AAr rrrr AAAA
		//out := (b[1] >> 3) & 0xff
		AAAA := ((b[1] & 0x06) << 3) | ( b[0] & 0x0f)
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, ioaddr: AAAA, source: Rr}
		fmt.Printf("%.4x\tout\t0x%.2x, r%d\t\t;%d\n", b2u16big(b), i.ioaddr, i.source, i.ioaddr)
	case INSN_RJMP:
		// 1100 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x800) >> 11) == 1 {
			i.kaddress = int16((k + 0xf000) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\trjmp\t.+%d\n", b2u16big(b), i.kaddress)
	case INSN_LDI:
		// 1110 KKKK dddd KKKK
		Kdata := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		i := Instr{label: m.label, family: m.family, dest: Rd, kdata: Kdata}
		fmt.Printf("%.4x\tldi\tr%d, 0x%.2x\t\t;%d\n", b2u16big(b), i.dest, i.kdata, i.dest)
	case INSN_RCALL:
		// 1101 kkkk kkkk kkkk
		k := (uint32(b[1]&0x0f)<<8 | uint32(b[0]))
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x0800) >> 11) == 1 {
			i.kaddress = int16((k + 0xf000) << 1)
			fmt.Printf("%.4x\trcall\t+.%d\n", b2u16big(b), i.kaddress)
		} else {
			i.kaddress = int16(k << 1)
			fmt.Printf("%.4x\trcall\t.%d\n", b2u16big(b), i.kaddress)
		}
	case INSN_SBI:
		// 1001 1010 AAAA Abbb
		AAAA := b[0] >> 3
		bbb := b[0] & 0x7
		i := Instr{label: m.label, family: m.family, ioaddr: AAAA, registerBit: bbb}
		fmt.Printf("%.4x\tsbi\t0x%x, %d\n", b2u16big(b), i.ioaddr, i.registerBit)
	case INSN_CBI:
		//1001 1000 AAAA Abbb
		AAAA := (b[0] & 0xf8) >> 3
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, ioaddr: AAAA, registerBit: bbb}
		fmt.Printf("%.4x\tcbi\t0x%.2x, %d\n", b2u16big(b), i.ioaddr, i.registerBit)
	case INSN_SBIC:
		// 1001 1001 AAAA Abbb
		AAAA := (b[0] & 0xf8) >> 3
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, ioaddr: AAAA, registerBit: bbb}
		fmt.Printf("%.4x\tsbic\t0x%.2x, %d\n", b2u16big(b), i.ioaddr, i.registerBit)
	case INSN_SBIS:
		// 1001 1011 AAAA Abbb
		AAAA := (b[0] & 0xf8) >> 3
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, ioaddr: AAAA, registerBit: bbb}
		fmt.Printf("%.4x\tsbis\t0x%.2x, %d\n", b2u16big(b), i.ioaddr, i.registerBit)
	case INSN_BLD:
		// 1111 100d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, dest: Rd, registerBit: bbb}
		fmt.Printf("%.4x\tbld\tr%d, %d\n", b2u16big(b), i.dest, i.registerBit)
	case INSN_BST:
		// 1111 101d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, dest: Rd, registerBit: bbb}
		fmt.Printf("%.4x\tbst\tr%d, %d\n", b2u16big(b), i.dest, i.registerBit)
	case INSN_SBRC:
		// 1111 110r rrrr 0bbb
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bbb := b[0] & 0x07
		i := Instr{label: m.label, family: m.family, source: Rr, registerBit: bbb}		
		fmt.Printf("%.4x\tsbrc\tr%d, %d\n", b2u16big(b), i.source, i.registerBit)
	case INSN_STS:
		// 1001 001d dddd 0000 kkkk kkkk kkkk kkkk
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		i := Instr{label: m.label, family: m.family, dest: Rd, kaddress: b2i16little(c)}
		fmt.Printf("%.4x\tsts\t0x%.4x, r%d\n", b2u16big(b), i.kaddress, i.dest)
	case INSN_LDS:
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		i := Instr{label: m.label, family: m.family, dest: Rd, kaddress: b2i16little(c)}
		fmt.Printf("%.4x\tlds\tr%d, 0x%.4x\n", b2u16big(b), i.dest, i.kaddress)
	case INSN_JMP:
		// 1001 010k kkkk 110k kkkk kkkk kkkk kkkk
		// XXX ToDo(erin): THIS HAS NOT BEEN TESTED
		k1 := (b[1] & 0x01) << 21
		k2 := ((b[0] & 0xf0) >> 3) | ((b[1] & 0x01) << 3)
		c := pop(2)
		k := k1 | k2 //| c mismatched types.
		fmt.Printf("%.4x\tlds\t0x%.4x\t;%d\n", b2u16big(b), b2u16little(c), k)
		i := Instr{label: m.label, family: m.family}
		fmt.Println(i)
	case INSN_ADD:
		// 0000 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tadd\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_ADIW:
		// 1001 0110 KKdd KKKK
		// 24,26,28,30
		dm := map[byte]byte {
			0: 24,
			1: 26,
			2: 28,
			3: 30,
		}
		Kdata := ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		i := Instr{label: m.label, family: m.family, kdata: Kdata, dest: dm[Rd]}
		fmt.Printf("%.4x\tadiw\tr%d, 0x%.2x\n", b2u16big(b), i.dest , i.kdata)
	case INSN_SBIW:
		// 1001 0111 KKdd KKKK
		dm := map[byte]byte{
			0: 24,
			1: 26,
			2: 28,
			3: 30,
		}
		Kdata := ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		i := Instr{label: m.label, family: m.family, kdata: Kdata, dest: dm[Rd]}
		fmt.Printf("%.4x\tsbiw\tr%d, 0x%.2x\n", b2u16big(b), i.dest, i.kdata)
	case INSN_ANDI:
		//0111 KKKK dddd KKKK
		Kdata := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		i := Instr{label: m.label, family: m.family, kdata: Kdata, dest: Rd}
		fmt.Printf("%.4x\tandi\tr%d, 0x%.2x\n", b2u16big(b), i.dest, i.kdata)
	case INSN_BRCC:
		//1111 01kk kkkk k000
		// 64 ≤ k ≤ +63
		k := (b2u16little(b) & 0x03f8) >> 3
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrcc\t.+%d\n", b2u16big(b), i.kaddress)
	case INSN_BRCS:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k000
		k := (b2u16little(b) & 0x03f8) >> 3
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrcs\t.+%d\n", b2u16big(b), i.kaddress)
	case INSN_BREQ:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbreq\t.%d\n", b2u16big(b), i.kaddress)
	case INSN_BRGE:
		// 1111 01kk kkkk k100
		k := (b2u16little(b) & 0x03f8) >> 3
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrge\t.%d\n", b2u16big(b), i.kaddress)
	case INSN_BRNE:
		// 1111 01kk kkkk k001
		k := (b2u16little(b) & 0x03f8) >> 3
		// check to see if msb of k is 1
		// if it is, the result is negative.
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrne\t.+%d\n", b2u16big(b), i.kaddress)
	case INSN_BRTC:
		// 1111 01kk kkkk k110
		k := (b2u16little(b) & 0x03f8) >> 3
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x40) >> 6) == 1 {
			i.kaddress = int16((k + 0xff80) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\tbrtc\t.+%d\n", b2u16big(b), i.kaddress)
	case INSN_COM:
		// 1001 010d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tcom\tr%d\n", b2u16big(b), i.dest)
	case INSN_CP:
		// 0001 01rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tcp\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_CPC:
		// 0000 01rd dddd rrrr
		Rr := ((b[1] & 0x02) << 3) | b[0]&0x0f
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tcpc\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_CPI:
		// 0011 KKKK dddd KKKK
		KKKK := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		i := Instr{label: m.label, family: m.family, kdata: KKKK, dest: Rd}
		fmt.Printf("%.4x\tcpi\tr%d, 0x%.2x\n", b2u16big(b), i.kdata, i.dest)
	case INSN_CPSE:
		// 0001 00rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tcpse\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_DEC:
		// 1001 010d dddd 1010
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tdec\tr%d\n", b2u16big(b), i.dest)
	case INSN_IN:
		// 1011 0AAd dddd AAAA
		address := ((b[1] & 0x09) << 3) | (b[0] & 0x0f)
		Rd := ((b[1] & 0xf1) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, ioaddr: address, dest: Rd}
		fmt.Printf("%.4x\tin\tr%d, 0x%.2x\n", b2u16big(b), i.dest, i.ioaddr)
	case INSN_LDDY:
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd, offset: m.offset}
		fmt.Printf("%.4x\tldd\tY+%d, r%d\n", b2u16big(b), i.offset, i.dest)
	case INSN_LDDZ:
		// 1001 000d dddd 0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd, offset: m.offset}
		fmt.Printf("%.4x\tldd\tr%d, Z+%d\n", b2u16big(b), i.dest, i.offset)
	case INSN_LDX:
		// 1001 000d dddd 1100
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tld\tr%d, X\n", b2u16big(b), i.dest)
	case INSN_LDXP:
		//1001 000d dddd 1101
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tld\tr%d, X+\n", b2u16big(b), i.dest)
	case INSN_LDY:
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tld\tr%d, Y\n", b2u16big(b), i.dest)
	case INSN_LDZ:
		// 1000 000d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tld\tr%d, Z\n", b2u16big(b), i.dest)
	case INSN_LPMZ:
		//z  1001 000d dddd 0100
		//z+ 1001 000d dddd 0101
		s := b[0] & 0x0f
		// XXX ToDo(Erin) Not sure this works.Rd := (b2u16little(b) & 0x01f0) >> 4
		Rd := (b2u16little(b) & 0x01f0) >> 4
		if s == 6 {
			fmt.Printf("%.4x\tlpm\tr%d, Z\n", b2u16big(b), Rd)
		} else {
			fmt.Printf("%.4x\tlpm\tr%d, Z+\n", b2u16big(b), Rd)
		}
	case INSN_LPM:	
		// 1001 0101 1100 1000
		// XXX ToDo: not tested
		// i := Instr{label: m.label, family: m.family, source: 0x00}
		fmt.Printf("%.4x\tlpm\n", b2u16big(b))
	case INSN_LSR:
		//1001 010d dddd 0110
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tlsr\tr%d\n", b2u16big(b), i.dest)
	case INSN_MOV:
		// 0010 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tmov\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_MOVW:
		// 0000 0001 dddd rrrr
		Rd := (b[0] & 0xf0) >> 3
		Rr := (b[0] & 0x0f) << 1
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tmovw\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_MUL:
		// 1001 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Printf("%.4x\tmul\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_NEG:
		// 1001 010d dddd s0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tneg\tr%d\n", b2u16big(b), i.dest)
	case INSN_OR:
		// 0010 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, dest: Rd, source: Rr}
		fmt.Printf("%.4x\tor\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_ORI:
		// 0110 KKKK dddd KKKK
		Kdata := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		i := Instr{label: m.label, family: m.family, kdata: Kdata, dest: Rd}
		fmt.Printf("%.4x\tori\tr%d, %.2x\n", b2u16big(b), i.dest, i.kdata)
	case INSN_POP:
		// 1001 000d dddd 1111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tpop\t r%d\n", b2u16big(b), i.dest)
	case INSN_PUSH:
		//1001 001d dddd 1111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tpush\tr%d\n", b2u16big(b), i.dest)
	case INSN_RET:
		i := Instr{label: m.label, family: m.family}
		fmt.Printf("%.4x\tret\n", b2u16big(b))
		fmt.Println(i)
	case INSN_RETI:
		i := Instr{label: m.label, family: m.family}
		fmt.Printf("%.4x\treti\n", b2u16big(b))
		fmt.Println(i)
	case INSN_ROR:
		// 1001 010d dddd 0111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, dest: Rd}
		fmt.Printf("%.4x\tror\tr%d\n", b2u16big(b), i.dest)
	case INSN_SBC:
		// 0000 10rd rrrr dddd
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, dest: Rd, source: Rr}
		fmt.Printf("%.4x\tsbc\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_SBCI:
		// 0100 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		Kdata := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		i := Instr{label: m.label, family: m.family, dest: Rd, kdata: Kdata}
		fmt.Printf("%.4x\tsbci\tr%d, 0x%x\n", b2u16big(b), i.dest, i.kdata)
	case INSN_SEI:
		// 1001 0100 0111 1000
		i := Instr{label: m.label, family: m.family}
		fmt.Printf("%.4x\tsei\n", b2u16big(b))
		fmt.Println(i)
	case INSN_STDY:
		// 1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr, offset: m.offset}
		fmt.Printf("%.4x\tstd\tY+%d, r%d\n", b2u16big(b), i.offset, i.source)
	case INSN_STDZ:
		// 10q0 qq1r rrrr 0qqq
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr, offset: m.offset}
		fmt.Printf("%.4x\tstd\tZ+%d, r%d\n", b2u16big(b), i.offset, i.source)
	case INSN_STX:
		// 1001 001r rrrr 1100
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\tX, r%d\n", b2u16big(b), i.source)
	case INSN_STXP:
		// 1001 001r rrrr 1101
		Rr := (b2u16little(b) & 0x01f0) >> 4
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\tX+, r%d\n", b2u16big(b), i.source)
	case INSN_STXM:
		// 1001 001r rrrr 1110
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\t-X, r%d\n", b2u16big(b), i.source)
	case INSN_STY:
		//1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\tY, r%d\n", b2u16big(b), i.source)
	case INSN_STZ:
		// 1000 001r rrrr 0000
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\tZ, r%d\n", b2u16big(b), i.source)
	case INSN_STZP:
		// 1001 001r rrrr 0001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\tZ+, r%d\n", b2u16big(b), i.source)
	case INSN_STZM:
		// 1001 001r rrrr 0010
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		i := Instr{label: m.label, family: m.family, source: Rr}
		fmt.Printf("%.4x\tst\t-Z, r%d\n", b2u16big(b), i.source)
	case INSN_SUB:
		// 0001 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		i := Instr{label: m.label, family: m.family, dest: Rd, source: Rr}
		fmt.Printf("%.4x\tsub\tr%d, r%d\n", b2u16big(b), i.dest, i.source)
	case INSN_SUBI:
		// 0101 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		Kdata := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		i := Instr{label: m.label, family: m.family, dest: Rd, kdata: Kdata}
		fmt.Printf("%.4x\tsubi\tr%d, 0x%x\n", b2u16big(b), i.dest, i.kdata)
	default:
		fmt.Printf("None of the above. Got %s (0x%.4x)\n", m.mnemonic, b2u16big(b))
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
			op.label = INSTN_STDY
		}
	default:
		op.mnemonic = "Unknown"
		op.value = b
	}
	return op
}
