package main

import (
	"fmt"
)


func dissAssemble(b []byte) {
	m := LookUp(b)

	switch m.Name {

	case "nop":
		fmt.Println("nop")
	case "adc":
		// CPSE, ADC
		fmt.Printf("%.4x\tadc\n",  bigEndianConcat(b))
	case "eor":
		r := b[1] & 0x02
		d := b[1] & 0x01

		rrrr := b[0] & 0x0f
		dddd := (b[0] >> 4) & 0x0f

		fmt.Printf("%.4x\teor\tr%d,r%d\t\t;??%b %b\n", bigEndianConcat(b), rrrr, dddd, r, d)
	case "out":
		//out := (b[1] >> 3) & 0xff
		AA := (b[1] & 0x06) >> 1
		i := littleEndianConcat(b)
		Rr := (i & 0x01f0) >> 4
		address := AA<<4 | (b[0] & 0x0f)
		fmt.Printf("%.4x\tout\t0x%.2x,r%d\t\t;%d\n", bigEndianConcat(b), address, Rr, address)

	case "cli":
		fmt.Printf("%.4x\tcli\n", bigEndianConcat(b))
		
	case "rjmp":
		i := littleEndianConcat(b)
		// Something about two's complement to make negative something something goes here.
		k := (i & 0x0fff) << 1
		fmt.Printf("%.4x\trjmp\t.+%d\n", bigEndianConcat(b), k)
	case "ldi":
		// this does not work
		K1 := b[1] & 0x0f
		K0 := b[0] & 0x0f
		cdata := (uint16(K1) << 4) | uint16(K0)
		Rd := (b[0] & 0xf0) >> 1
		fmt.Printf("%.4x\tldi\tr%x,0x%.2x\t\t;%d\n", bigEndianConcat(b), Rd, cdata, Rd)
	case "rcall":
		fmt.Printf("%.4x\trcall\n", bigEndianConcat(b))
	case "sbi":
		AAAA := b[0] >> 3
		bbb := b[0] & 0x7
		fmt.Printf("%.4x\tsbi\t0x%x,%d\n", bigEndianConcat(b),AAAA, bbb)
	case "sts":
		c := pop(2)
		fmt.Printf("%.4x\tsts\t0x%.4x\n", bigEndianConcat(b), bigEndianConcat(c))
	case "lds":
		c := pop(2)
		fmt.Printf("%.4x\tlds\t0x%.4x\n", bigEndianConcat(b), bigEndianConcat(c))
	case "add":
		fmt.Printf("%.4x\tadd\n",  bigEndianConcat(b))
	case "adiw":
		fmt.Printf("%.4x\tadiw\n",  bigEndianConcat(b))
	case "andi":
		fmt.Printf("%.4x\tandi\n",  bigEndianConcat(b))
	case "bld":
		fmt.Printf("%.4x\tbld\n",  bigEndianConcat(b))
	case "brcc":
		fmt.Printf("%.4x\tbrcc\n",  bigEndianConcat(b))
	case "brcs":
		fmt.Printf("%.4x\tbrcs\n",  bigEndianConcat(b))
	case "breq":
		fmt.Printf("%.4x\tbreq\n",  bigEndianConcat(b))
	case "brge":
		fmt.Printf("%.4x\tbrge\n",  bigEndianConcat(b))
	case "brne":
		fmt.Printf("%.4x\tbrne\n",  bigEndianConcat(b))
	case "brtc":
		fmt.Printf("%.4x\tbrtc\n",  bigEndianConcat(b))
	case "bst":
		fmt.Printf("%.4x\tbst\n",  bigEndianConcat(b))
	case "cbi":
		fmt.Printf("%.4x\tcbi\n",  bigEndianConcat(b))
	case "com":
		fmt.Printf("%.4x\tcom\n",  bigEndianConcat(b))
	case "cp":
		fmt.Printf("%.4x\tcp\n",  bigEndianConcat(b))
	case "cpc":
		fmt.Printf("%.4x\tcpc\n",  bigEndianConcat(b))
	case "cpi":
		fmt.Printf("%.4x\tcpi\n",  bigEndianConcat(b))
	case "cpse":
		fmt.Printf("%.4x\tcpse\n",  bigEndianConcat(b))
	case "dec":
		fmt.Printf("%.4x\tdec\n",  bigEndianConcat(b))
	case "in":
		fmt.Printf("%.4x\tin\n",  bigEndianConcat(b))
	case "lddy+":
		fmt.Printf("%.4x\tldd Y+%d\n",  bigEndianConcat(b), m.Offset)
	case "lddz+":
		fmt.Printf("%.4x\tldd Z+%d\n",  bigEndianConcat(b), m.Offset)
	case "ldx":
		fmt.Printf("%.4x\tldx\n",  bigEndianConcat(b))
	case "ldx+":
		fmt.Printf("%.4x\tldx+\n",  bigEndianConcat(b))
	case "ldy":
		fmt.Printf("%.4x\tldy\n",  bigEndianConcat(b))
	case "ldz":
		fmt.Printf("%.4x\tldz\n",  bigEndianConcat(b))
	case "lpmz+":
		fmt.Printf("%.4x\tlpmz+\n",  bigEndianConcat(b))
	case "lsr":
		fmt.Printf("%.4x\tlsr\n",  bigEndianConcat(b))
	case "mov":
		fmt.Printf("%.4x\tmov\n",  bigEndianConcat(b))
	case "movw":
		fmt.Printf("%.4x\tmovw\n",  bigEndianConcat(b))
	case "mul":
		fmt.Printf("%.4x\tmul\n",  bigEndianConcat(b))
	case "neg":
		fmt.Printf("%.4x\tneg\n",  bigEndianConcat(b))
	case "or":
		fmt.Printf("%.4x\tor\n",  bigEndianConcat(b))
	case "ori":
		fmt.Printf("%.4x\tori\n",  bigEndianConcat(b))
	case "pop":
		fmt.Printf("%.4x\tpop\n",  bigEndianConcat(b))
	case "push":
		fmt.Printf("%.4x\tpush\n",  bigEndianConcat(b))
	case "ret":
		fmt.Printf("%.4x\tret\n",  bigEndianConcat(b))
	case "reti":
		fmt.Printf("%.4x\treti\n",  bigEndianConcat(b))
	case "ror":
		fmt.Printf("%.4x\tror\n",  bigEndianConcat(b))
	case "sbc":
		fmt.Printf("%.4x\tsbc\n",  bigEndianConcat(b))
	case "sbci":
		fmt.Printf("%.4x\tsbci\n",  bigEndianConcat(b))
	case "sbic":
		fmt.Printf("%.4x\tsbic\n",  bigEndianConcat(b))
	case "sbis":
		fmt.Printf("%.4x\tsbis\n",  bigEndianConcat(b))
	case "sbiw":
		fmt.Printf("%.4x\tsbiw\n",  bigEndianConcat(b))
	case "sbrc":
		fmt.Printf("%.4x\tsbrc\n",  bigEndianConcat(b))
	case "sei":
		fmt.Printf("%.4x\tsei\n",  bigEndianConcat(b))
	case "stdy+":
		fmt.Printf("%.4x\tstd Y+%d\n",  bigEndianConcat(b), m.Offset)
	case "stdz+":
		fmt.Printf("%.4x\tstd Z+%d\n",  bigEndianConcat(b), m.Offset)
	case "stx":
		fmt.Printf("%.4x\tstx\n",  bigEndianConcat(b))
	case "stx+":
		r := (littleEndianConcat(b) & 0x01f0) >> 4
		fmt.Printf("%.4x\tst X+, r%d\n",  bigEndianConcat(b), r)
	case "stx-":
		fmt.Printf("%.4x\tstx-\n",  bigEndianConcat(b))
	case "sty":
		fmt.Printf("%.4x\tsty\n",  bigEndianConcat(b))
	case "stz":
		fmt.Printf("%.4x\tstz\n",  bigEndianConcat(b))
	case "stz+":
		fmt.Printf("%.4x\tstz+\n",  bigEndianConcat(b))
	case "stz-":
		fmt.Printf("%.4x\tstz-\n",  bigEndianConcat(b))
	case "sub":
		fmt.Printf("%.4x\tsub\n",  bigEndianConcat(b))
	case "subi":
		fmt.Printf("%.4x\tsubi\n",  bigEndianConcat(b))
	default:
		fmt.Printf("None of the above. Got %s (0x%.4x)\n", m.Name, bigEndianConcat(b) )
	}
}

func LookUp(raw []byte) OpCode {
	var op OpCode
	b := littleEndianConcat(raw)
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
	x := b  & 0xd208
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
