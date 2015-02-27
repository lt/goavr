package main

import (
	"fmt"
)

func dissAssemble(b []byte) {
	m := LookUp(b)

	switch m.Name {

	case "nop":
		fmt.Printf("%.4x\tnop\n", b2u16big(b))
	case "adc":
		// 0001 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "eor":
		r := b[1] & 0x02
		d := b[1] & 0x01
		rrrr := b[0] & 0x0f
		dddd := (b[0] >> 4) & 0x0f
		fmt.Printf("%.4x\teor\tr%d, r%d\t;??%b %b\n", b2u16big(b), rrrr, dddd, r, d)
	case "out":
		//out := (b[1] >> 3) & 0xff
		AA1 := (b[1] & 0x06) >> 1
		AA2 := b[0] & 0x0f
		address := AA1<<4 | AA2
		i := b2u16little(b)
		Rr := (i & 0x01f0) >> 4
		fmt.Printf("%.4x\tout\t0x%.2x, r%d\t\t;%d\n", b2u16big(b), address, Rr, address)
	case "cli":
		fmt.Printf("%.4x\tcli\n", b2u16big(b))
	case "jmp":
		// 1001 010k kkkk 110k kkkk kkkk kkkk kkkk
		// XXX ToDo(erin): not tested.
		k1 := (b[1] & 0x01) << 21
		k2 := ((b[0] & 0xf0) >> 3) | ((b[1] & 0x01) << 3)
		c := pop(2)
		k := k1 | k2 //| c mismatched types.
		fmt.Printf("%.4x\tlds\t0x%.4x\t;%d\n", b2u16big(b), b2u16little(c), k)
	case "rjmp":
		// 1100 KKKK dddd KKKK
		k := (b2u16little(b) & 0x0fff) << 1
		if ((k & 0x800) >> 11) == 1 {
			i := -b2i16little(b)
			nk := (i & 0x0fff) << 1
			fmt.Printf("%.4x\trjmp\t.-%d\n", b2u16big(b), nk)
		} else {
			fmt.Printf("%.4x\trjmp\t.+%d\n", b2u16big(b), k)
		}
	case "ldi":
		K1 := b[1] & 0x0f
		K2 := b[0] & 0x0f
		KKKK := (K1 << 4) | K2
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tldi\tr%d, 0x%.2x\t\t;%d\n", b2u16big(b), Rd, KKKK, Rd)
	case "rcall":
		// 1101 kkkk kkkk kkkk
		k := (b2u16little(b) & 0x0fff) << 1
		if k&0x0400 == 0 { //?yes?no?maybe?
			fmt.Printf("%.4x\trcall\t.+%d\n", b2u16big(b), k)
		} else {
			i := -b2i16little(b)
			nk := (i & 0x0fff) << 1
			fmt.Printf("%.4x\trcall\t.-%d\n", b2u16big(b), nk)
		}
	case "sbi":
		AAAA := b[0] >> 3
		bbb := b[0] & 0x7
		fmt.Printf("%.4x\tsbi\t0x%x, %d\n", b2u16big(b), AAAA, bbb)
	case "sts":
		// 1001 001d dddd 0000 kkkk kkkk kkkk kkkk
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		fmt.Printf("%.4x\tsts\t0x%.4x, r%d\n", b2u16big(b), b2u16little(c), Rd)
	case "lds":
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		c := pop(2)
		fmt.Printf("%.4x\tlds\t0x%.4x, r%d\n", b2u16big(b), b2u16little(c), Rd)
	case "add":
		// 0000 11rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tadd\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "adiw":
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
	case "andi":
		//0111 KKKK dddd KKKK
		K := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tandi\tr%d, 0x%.2x\n", b2u16big(b), Rd, K)
	case "bld":
		// 1111 100d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tbld\tr%d, %d\n", b2u16big(b), Rd, bit)
	case "brcc":
		//1111 01kk kkkk k000
		// 64 ≤ k ≤ +63
		k := ((b[1] & 0x03) << 5) | ((b[0]&0xf8)>>3)<<1
		fmt.Printf("%.4x\tbrcc\t.+%d\n", b2u16big(b), k)
	case "brcs":
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k000
		k := (((b[1] & 0x03) << 5) | ((b[0] & 0xf8) >> 3)) << 1
		fmt.Printf("%.4x\tbrcs\t.+%d\n", b2u16big(b), k)
	case "breq":
		// Supposed to be -64<k<+63, but avr-objdump doesn't display
		// these values this way.
		// 1111 00kk kkkk k001
		k := (((b[1] & 0x03) << 5) | ((b[0] & 0xf8) >> 3)) << 1
		fmt.Printf("%.4x\tbreq\t.%d\n", b2u16big(b), k)
	case "brge":
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
	case "brne":
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
	case "brtc":
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
	case "bst":
		// 1111 101d dddd 0bbb
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tbst\tr%d, %d\n", b2u16big(b), Rd, bit)
	case "cbi":
		//1001 1000 AAAA Abbb
		address := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tcbi\t0x%.2x, %d\n", b2u16big(b), address, bit)
	case "com":
		// 1001 010d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcom\tr%d\n", b2u16big(b), Rd)
	case "cp":
		// 0001 01rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcp\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "cpc":
		Rr := ((b[1] & 0x02) << 3) | b[0]&0x0f
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tcpc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "cpi":
		K1 := b[1] & 0x0f
		K2 := b[0] & 0x0f
		KKKK := (K1 << 4) | K2
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tcpi\tr%d, 0x%.2x\n", b2u16big(b), Rd, KKKK)
	case "cpse":
		// 0001 00rd dddd rrrr
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tcpse\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "dec":
		// 1001 010d dddd 1010
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tdec\tr%d\n", b2u16big(b), Rd)
	case "in":
		// 1011 0AAd dddd AAAA
		address := ((b[1] & 0x09) << 3) | (b[0] & 0x0f)
		Rd := ((b[1] & 0xf1) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tin\tr%d, 0x%.2x\n", b2u16big(b), Rd, address)
	case "lddy+":
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tldd\tY+%d, r%d\n", b2u16big(b), m.Offset, Rd)
	case "lddz+":
		// 1001 000d dddd 0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tldd\tr%d, Z+%d\n", b2u16big(b), Rd, m.Offset)
	case "ldx":
		// 1001 000d dddd 1100
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X\n", b2u16big(b), Rd)
	case "ldx+":
		//1001 000d dddd 1101
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, X+\n", b2u16big(b), Rd)
	case "ldy":
		// 1001 000d dddd 1001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Y\n", b2u16big(b), Rd)
	case "ldz":
		// 1000 000d dddd 0000
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tld\tr%d, Z\n", b2u16big(b), Rd)
	case "lpmz":
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
	case "lpm":
		// 1001 0101 1100 1000
		fmt.Printf("%.4x\tlpm\n", b2u16big(b))
	case "lsr":
		//1001 010d dddd 0110
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tlsr\tr%d\n", b2u16big(b), Rd)
	case "mov":
		// 0010 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmov\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "movw":
		// 0000 0001 dddd rrrr
		Rd := (b[0] & 0xf0) >> 3
		Rr := (b[0] & 0x0f) << 1
		fmt.Printf("%.4x\tmovw\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "mul":
		// 1001 11rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tmul\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "neg":
		// 1001 010d dddd s0001
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tneg\tr%d\n", b2u16big(b), Rd)
	case "or":
		// 0010 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tor\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "ori":
		// 0110 KKKK dddd KKKK
		K := ((b[1] & 0x0f) << 4) | (b[0] & 0x0f)
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		fmt.Printf("%.4x\tori\tr%d, %.2x\n", b2u16big(b), Rd, K)
	case "pop":
		// 1001 000d dddd 1111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tpop\t r%d\n", b2u16big(b), Rd)
	case "push":
		//1001 001d dddd 1111
		Rd := (b2u16little(b) & 0x01f0) >> 4
		fmt.Printf("%.4x\tpush\tr%d\n", b2u16big(b), Rd)
	case "ret":
		fmt.Printf("%.4x\tret\n", b2u16big(b))
	case "reti":
		fmt.Printf("%.4x\treti\n", b2u16big(b))
	case "ror":
		// 1001 010d dddd 0111
		Rd := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tror\tr%d\n", b2u16big(b), Rd)
	case "sbc":
		// 0000 10rd rrrr dddd
		Rr := (((b[1]&0x02)>>1)<<4 | (b[0] & 0x0f))
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsbc\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "sbci":
		// 0100 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		KKKK := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsbci\tr%d, 0x%x\n", b2u16big(b), Rd, KKKK)
	case "sbic":
		// 1001 1001 AAAA Abbb
		A := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tsbic\t0x%.2x, %d\n", b2u16big(b), A, bit)
	case "sbis":
		// 1001 1011 AAAA Abbb
		address := (b[0] & 0xf8) >> 3
		bit := b[0] & 0x07
		fmt.Printf("%.4x\tsbis\t0x%.2x, %d\n", b2u16big(b), address, bit)
	case "sbiw":
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
	case "sbrc":
		fmt.Printf("%.4x\tsbrc\n", b2u16big(b))
	case "sei":
		fmt.Printf("%.4x\tsei\n", b2u16big(b))
	case "stdy+":
		// 1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tY+%d, r%d\n", b2u16big(b), m.Offset, Rr)
	case "stdz+":
		// 10q0 qq1r rrrr 0qqq
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tstd\tZ+%d, r%d\n", b2u16big(b), m.Offset, Rr)
	case "stx":
		// 1001 001r rrrr 1100
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tX, r%d\n", b2u16big(b), Rr)
	case "stx+":
		// 1001 001r rrrr 1101
		Rr := (b2u16little(b) & 0x01f0) >> 4
		fmt.Printf("%.4x\tst\tX+, r%d\n", b2u16big(b), Rr)
	case "stx-":
		// 1001 001r rrrr 1110
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-X, r%d\n", b2u16big(b), Rr)
	case "sty":
		//1001 001r rrrr 1001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tY, r%d\n", b2u16big(b), Rr)
	case "stz":
		// 1000 001r rrrr 0000
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ, r%d\n", b2u16big(b), Rr)
	case "stz+":
		// 1001 001r rrrr 0001
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\tZ+, r%d\n", b2u16big(b), Rr)
	case "stz-":
		// 1001 001r rrrr 0010
		Rr := ((b[1] & 0x01) << 4) | ((b[0] & 0xf0) >> 4)
		fmt.Printf("%.4x\tst\t-Z, r%d\n", b2u16big(b), Rr)
	case "sub":
		// 0001 10rd dddd rrrr
		Rr := (((b[1] & 0x02) >> 1) << 4) | (b[0] & 0x0f)
		Rd := ((b[1]&0x01)<<4 | ((b[0] & 0xf0) >> 4))
		fmt.Printf("%.4x\tsub\tr%d, r%d\n", b2u16big(b), Rd, Rr)
	case "subi":
		// 0101 KKKK dddd KKKK
		Rd := ((b[0] & 0xf0) >> 4) + 0x10
		KKKK := (b[1]&0x0f)<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tsubi\tr%d, 0x%x\n", b2u16big(b), Rd, KKKK)
	default:
		fmt.Printf("None of the above. Got %s (0x%.4x)\n", m.Name, b2u16big(b))
	}
}

func LookUp(raw []byte) OpCode {
	var op OpCode
	b := b2u16little(raw)
	for _, entry := range OpCodeLookUpTable {
		v := b & entry.Mask
		if v == entry.Value {
			op = entry
			switch entry.Name {
			case "std":
				return deConvoluter(b, op)
			case "ldd":
				return deConvoluter(b, op)
			}
			return op
		} else {
			op = OpCode{Name: "Unknown", Value: b}
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
			op.Name = "ldz"
		} else {
			op.Name = "lddz+"
			op.Offset = offset
		}
	case 0x8008:
		if offset == 0 {
			op.Name = "ldy"
		} else {
			op.Name = "lddy+"
			op.Offset = offset
		}
	case 0x8200:
		if offset == 0 {
			op.Name = "stz"
		} else {
			op.Name = "stdz+"
			op.Offset = offset
		}
	case 0x8208:
		if offset == 0 {
			op.Name = "sty"
		} else {
			op.Name = "stdy+"
			op.Offset = offset
		}
	default:
		op.Name = "Unknown"
		op.Value = b
	}
	return op
}
