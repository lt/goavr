package main

import (
	"fmt"
)

func dissAssemble(b []byte) {
	m := lookUp(b)

	switch m.label {

	case INSN_NOP:
		fmt.Printf("%.4x\tnop\n", b2u16big(b))
		i := Instr{label: m.label, family: m.family}
		fmt.Println(i)
	case INSN_ADC:
		// 0001 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
		i := Instr{label: m.label, family: m.family, source: Rr, dest: Rd}
		fmt.Println(i)
	case INSN_EOR:
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\teor\tr%d, r%d\n", b2u16big(b), Rr, Rd)
		i := Instr{label: m.label, source: Rr, dest: Rd, family: m.family}
		fmt.Println(i)
	case INSN_OUT:
		//out := (b[1] >> 3) & 0xff
		AA1 := (b[1] & 0x06) >> 1
		AA2 := b[0] & 0x0f
		address := AA1<<4 | AA2
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tout\t0x%.2x, r%d\t\t;%d\n", b2u16big(b), address, Rr, address)
		i := Instr{label: m.label, family: m.family, ioaddr: address, source: Rr}
		fmt.Println(i)
	case INSN_CLI:
		fmt.Printf("%.4x\tcli\n", b2u16big(b))
		i := Instr{label: m.label, family: m.family}
		fmt.Println(i)
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
	case INSN_RJMP:
		// 1100 kkkk kkkk kkkk
		k := (uint32(b[1] & 0x0f) << 8 | uint32(b[0]))
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x800) >> 11) == 1 {
			i.kaddress = int16((k + 0xf000) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\trjmp\t.+%d\n", b2u16big(b), i.address)
		fmt.Println(i)
	case INSN_LDI:
		K1 := b[1] & 0x0f
		K2 := b[0] & 0x0f
		KKKK := (K1 << 4) | K2
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tldi\tr%d, 0x%.2x\t\t;%d\n", b2u16big(b), Rd, KKKK, Rd)
	case INSN_RCALL:
		// 1101 kkkk kkkk kkkk
		k := (uint32(b[1] & 0x0f) << 8 | uint32(b[0]))
		i := Instr{label: m.label, family: m.family}
		if ((k & 0x0800) >> 11) == 1 { 
			i.kaddress = int16((k + 0xf000) << 1)
		} else {
			i.kaddress = int16(k << 1)
		}
		fmt.Printf("%.4x\trcall\t.%d\n", b2u16big(b), i.kaddress)
		fmt.Println(i)
	case INSN_SBI:
		AAAA := b[0] >> 3
		bbb := b[0] & 0x7
		fmt.Printf("%.4x\tsbi\t0x%x, %d\n", b2u16big(b), AAAA, bbb)
	case INSN_STS:
		// 1001 001d dddd 0000 kkkk kkkk kkkk kkkk
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		fmt.Printf("%.4x\tsts\t0x%.4x, r%d\n", b2u16big(b), b2u16little(c), Rd)
	case INSN_LDS:
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		fmt.Printf("%.4x\tlds\t0x%.4x, r%d\n", b2u16big(b), b2u16little(c), Rd)
	case INSN_ADD:
		// 0000 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadd\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_ADIW:
		// 1001 0110 KKdd KKKK
		// 24,26,28,30
		dm := map[byte]int{
			0: 24,
			1: 26,
			2: 28,
			3: 30,
		}
		k := ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		fmt.Printf("%.4x\tadiw\tr%d, 0x%.2x\n", b2u16big(b), dm[Rd], k)
	case INSN_ANDI:
		//0111 KKKK dddd KKKK
		K := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tandi\tr%d, 0x%.2x\n", b2u16big(b), Rd, K)
	case INSN_BLD:
		// 1111 100d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tbld\tr%d, %d\n", b2u16big(b), Rd, bit)
	case INSN_BRCC:
		//1111 01kk kkkk k000
		// 64 ≤ k ≤ +63
		k := ((b[1] & 0x03) << 5) | ((b[0]&0xf8)>>3)<<1
		fmt.Printf("%.4x\tbrcc\t.+%d\n", b2u16big(b), k)
	case INSN_BRCS:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k000
		k := (((b[1] & 0x03) << 5) | ((b[0] & 0xf8) >> 3)) << 1
		fmt.Printf("%.4x\tbrcs\t.+%d\n", b2u16big(b), k)
	case INSN_BREQ:
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k001
		k := (((b[1] & 0x03) << 5) | ((b[0] & 0xf8) >> 3)) << 1
		fmt.Printf("%.4x\tbreq\t.%d\n", b2u16big(b), k)
	case INSN_BRGE:
		// 1111 01kk kkkk k100
		k1 := ((b[1] & 0x03) << 5)
		k2 := ((b[0] & 0xf8) >> 3) << 1
		// XXX TODO(erin) -- this might need to be << 1 -- haven't found a positive one yet.
		k := k1 | k2
		if k <= 64 {
			fmt.Printf("%.4x\tbrge\t.%d\n", b2u16big(b), k)
		} else {
			i := -b2i16little(b)
			nk := (((i & 0x03f8) >> 3) + 1) << 1
			fmt.Printf("%.4x\tbrge\t.-%d\n", b2u16big(b), nk)
		}
	case INSN_BRNE:
		k := (b2u16little(b) & 0x03f8) >> 2
		// check to see if msb of k is 1
		// if it is, the result is negative.
		if ((b2u16little(b) & 0x0200) >> 9) == 1 {
			nk := ^(k) & 0x00ff
			nk += 1
			fmt.Printf("%.4x\tbrne\t.-%d\n", b2u16big(b), nk)
		} else {
			fmt.Printf("%.4x\tbrne\t.+%d\n", b2u16big(b), k)
		}
	case INSN_BRTC:
		// 1111 01kk kkkk k110
		k1 := ((b[1] & 0x03) << 5)
		k2 := ((b[0] & 0xf8) >> 3) << 1
		// XXX TODO(erin) -- this might need to be << 1 -- haven't found a positive one yet.
		k := k1 | k2
		if k <= 64 {
			fmt.Printf("%.4x\tbrtc\t.+%d\n", b2u16big(b), k)
		} else {
			i := -b2i16little(b)
			nk := (((i & 0x03f8) >> 3) + 1) << 1
			fmt.Printf("%.4x\tbrtc\t.-%d\n", b2u16big(b), nk)
		}
	case INSN_BST:
		// 1111 101d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tbst\tr%d, %d\n", b2u16big(b), Rd, bit)
	case INSN_CBI:
		//1001 1000 AAAA Abbb
		address := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tcbi\t0x%.2x, %d\n", b2u16big(b), address, bit)
	case INSN_COM:
		// 1001 010d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcom\tr%d\n", b2u16big(b), Rd)
	case INSN_CP:
		// 0001 01rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcp\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_CPC:
		Rr := ((b[1] & 0x02) << 3) | b[0]&0x0f
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcpc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_CPI:
		K1 := b[1] & 0x0f
		K2 := b[0] & 0x0f
		KKKK := (K1 << 4) | K2
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tcpi\tr%d, 0x%.2x\n", b2u16big(b), Rd, KKKK)
	case INSN_CPSE:
		// 0001 00rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcpse\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_DEC:
		// 1001 010d dddd 1010
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tdec\tr%d\n", b2u16big(b), Rd)
	case INSN_IN:
		// 1011 0AAd dddd AAAA
		address := ((b[1] & 0x09) << 3) | (b[0] & 0x0f)
		Rd := ((b[1] & 0xf1) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tin\tr%d, 0x%.2x\n", b2u16big(b), Rd, address)
	case INSN_LDDY:
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tldd\tY+%d, r%d\n", b2u16big(b), m.offset, Rd)
	case INSN_LDDZ:
		// 1001 000d dddd 0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tldd\tr%d, Z+%d\n", b2u16big(b), Rd, m.offset)
	case INSN_LDX:
		// 1001 000d dddd 1100
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X\n", b2u16big(b), Rd)
	case INSN_LDXP:
		//1001 000d dddd 1101
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X+\n", b2u16big(b), Rd)
	case INSN_LDY:
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Y\n", b2u16big(b), Rd)
	case INSN_LDZ:
		// 1000 000d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Z\n", b2u16big(b), Rd)
	case INSN_LPMZ:
		//z  1001 000d dddd 0100
		//z+ 1001 000d dddd 0101
		s := b[0] & 0x0f
		// XXX ToDo(Erin) Not sure this works.
		Rd := (b2u16little(b) & 0x01f0) >> 4
		if s == 6 {
			fmt.Printf("%.4x\tlpm\tr%d, Z\n", b2u16big(b), Rd)
		} else {
			fmt.Printf("%.4x\tlpm\tr%d, Z+\n", b2u16big(b), Rd)
		}
	case INSN_LPM:
		// 1001 0101 1100 1000
		fmt.Printf("%.4x\tlpm\n", b2u16big(b))
	case INSN_LSR:
		//1001 010d dddd 0110
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tlsr\tr%d\n", b2u16big(b), Rd)
	case INSN_MOV:
		// 0010 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmov\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_MOVW:
		// 0000 0001 dddd rrrr
		Rd := (b[0] & 0xf0) >> 3
		Rr := (b[0] & 0x0f) << 1
		fmt.Printf("%.4x\tmovw\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_MUL:
		// 1001 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmul\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_NEG:
		// 1001 010d dddd s0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tneg\tr%d\n", b2u16big(b), Rd)
	case INSN_OR:
		// 0010 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tor\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_ORI:
		// 0110 KKKK dddd KKKK
		K := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tori\tr%d, %.2x\n", b2u16big(b), Rd, K)
	case INSN_POP:
		// 1001 000d dddd 1111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tpop\t r%d\n", b2u16big(b), Rd)
	case INSN_PUSH:
		//1001 001d dddd 1111
		Rd := (b2u16little(b) & 0x01f0) >> 4
		fmt.Printf("%.4x\tpush\tr%d\n", b2u16big(b), Rd)
	case INSN_RET:
		fmt.Printf("%.4x\tret\n", b2u16big(b))
	case INSN_RETI:
		fmt.Printf("%.4x\treti\n", b2u16big(b))
	case INSN_ROR:
		// 1001 010d dddd 0111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tror\tr%d\n", b2u16big(b), Rd)
	case INSN_SBC:
		// 0000 10rd rrrr dddd
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsbc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_SBCI:
		// 0100 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		KKKK := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsbci\tr%d, 0x%x\n", b2u16big(b), Rd, KKKK)
	case INSN_SBIC:
		// 1001 1001 AAAA Abbb
		A := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tsbic\t0x%.2x, %d\n", b2u16big(b), A, bit)
	case INSN_SBIS:
		// 1001 1011 AAAA Abbb
		address := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tsbis\t0x%.2x, %d\n", b2u16big(b), address, bit)
	case INSN_SBIW:
		// 1001 0111 KKdd KKKK
		dm := map[byte]int{
			0: 24,
			1: 26,
			2: 28,
			3: 30,
		}
		k := ((b[0] & 0xc0) >> 2) | (b[0] & 0x0f)
		Rd := (b[0] & 0x30) >> 4
		fmt.Printf("%.4x\tsbiw\tr%d, 0x%.2x\n", b2u16big(b), dm[Rd], k)
	case INSN_SBRC:
		fmt.Printf("%.4x\tsbrc\n", b2u16big(b))
	case INSN_SEI:
		fmt.Printf("%.4x\tsei\n", b2u16big(b))
	case INSN_STDY:
		// 1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tY+%d, r%d\n", b2u16big(b), m.offset, Rr)
	case INSN_STDZ:
		// 10q0 qq1r rrrr 0qqq
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tZ+%d, r%d\n", b2u16big(b), m.offset, Rr)
	case INSN_STX:
		// 1001 001r rrrr 1100
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tX, r%d\n", b2u16big(b), Rr)
	case INSN_STXP:
		// 1001 001r rrrr 1101
		Rr := (b2u16little(b) & 0x01f0) >> 4
		fmt.Printf("%.4x\tst\tX+, r%d\n", b2u16big(b), Rr)
	case INSN_STXM:
		// 1001 001r rrrr 1110
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-X, r%d\n", b2u16big(b), Rr)
	case INSN_STY:
		//1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tY, r%d\n", b2u16big(b), Rr)
	case INSN_STZ:
		// 1000 001r rrrr 0000
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ, r%d\n", b2u16big(b), Rr)
	case INSN_STZP:
		// 1001 001r rrrr 0001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ+, r%d\n", b2u16big(b), Rr)
	case INSN_STZM:
		// 1001 001r rrrr 0010
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-Z, r%d\n", b2u16big(b), Rr)
	case INSN_SUB:
		// 0001 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsub\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case INSN_SUBI:
		// 0101 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		KKKK := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsubi\tr%d, 0x%x\n", b2u16big(b), Rd, KKKK)
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
