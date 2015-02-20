package main

import (
	"fmt"
)


func dissAssemble(b []byte) {
	m := LookUp(b)

	switch m {

	case "nop":
		fmt.Println("nop")
	case "adc":
		// CPSE, ADC
		fmt.Println("adc")

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
		c := pop(data, 2)
		fmt.Printf("sts\t0x%.4x\n", bigEndianConcat(c))
	case "lds":
		c := pop(data, 2)
		fmt.Printf("lds\t0x%.4x\n", bigEndianConcat(c))
	case "add":
		fmt.Printf("add\n")
	case "adiw":
		fmt.Printf("adiw\n")
	case "andi":
		fmt.Printf("andi\n")
	case "bld":
		fmt.Printf("bld\n")
	case "brcc":
		fmt.Printf("brcc\n")
	case "brcs":
		fmt.Printf("brcs\n")
	case "breq":
		fmt.Printf("breq\n")
	case "brge":
		fmt.Printf("brge\n")
	case "brne":
		fmt.Printf("brne\n")
	case "brtc":
		fmt.Printf("brtc\n")
	case "bst":
		fmt.Printf("bst\n")
	case "cbi":
		fmt.Printf("cbi\n")
	case "com":
		fmt.Printf("com\n")
	case "cp":
		fmt.Printf("cp\n")
	case "cpc":
		fmt.Printf("cpc\n")
	case "cpi":
		fmt.Printf("cpi\n")
	case "cpse":
		fmt.Printf("cpse\n")
	case "dec":
		fmt.Printf("dec\n")
	case "in":
		fmt.Printf("in\n")
	case "lddy+1":
		fmt.Printf("lddy+1\n")
	case "lddy+2":
		fmt.Printf("lddy+2\n")
	case "lddy+3":
		fmt.Printf("lddy+3\n")
	case "lddz+1":
		fmt.Printf("lddz+1\n")
	case "lddz+2":
		fmt.Printf("lddz+2\n")
	case "lddz+3":
		fmt.Printf("lddz+3\n")
	case "ldx":
		fmt.Printf("ldx\n")
	case "ldx+":
		fmt.Printf("ldx+\n")
	case "ldy":
		fmt.Printf("ldy\n")
	case "ldz":
		fmt.Printf("ldz\n")
	case "lpmz+":
		fmt.Printf("lpmz+\n")
	case "lsr":
		fmt.Printf("lsr\n")
	case "mov":
		fmt.Printf("mov\n")
	case "movw":
		fmt.Printf("movw\n")
	case "mul":
		fmt.Printf("mul\n")
	case "neg":
		fmt.Printf("neg\n")
	case "or":
		fmt.Printf("or\n")
	case "ori":
		fmt.Printf("ori\n")
	case "pop":
		fmt.Printf("pop\n")
	case "push":
		fmt.Printf("push\n")
	case "ret":
		fmt.Printf("ret\n")
	case "reti":
		fmt.Printf("reti\n")
	case "ror":
		fmt.Printf("ror\n")
	case "sbc":
		fmt.Printf("sbc\n")
	case "sbci":
		fmt.Printf("sbci\n")
	case "sbic":
		fmt.Printf("sbic\n")
	case "sbis":
		fmt.Printf("sbis\n")
	case "sbiw":
		fmt.Printf("sbiw\n")
	case "sbrc":
		fmt.Printf("sbrc\n")
	case "sei":
		fmt.Printf("sei\n")
	case "stdy+1":
		fmt.Printf("stdy+1\n")
	case "stdy+2":
		fmt.Printf("stdy+2\n")
	case "stdz+1":
		fmt.Printf("stdz+1\n")
	case "stdz+2":
		fmt.Printf("stdz+2\n")
	case "stdz+3":
		fmt.Printf("stdz+3\n")
	case "stx":
		fmt.Printf("stx\n")
	case "stx+":
		fmt.Printf("stx+\n")
	case "stx-":
		fmt.Printf("stx-\n")
	case "sty":
		fmt.Printf("sty\n")
	case "stz":
		fmt.Printf("stz\n")
	case "stz+":
		fmt.Printf("stz+\n")
	case "sub":
		fmt.Printf("sub\n")
	case "subi":
		fmt.Printf("subi\n")
	default:
		fmt.Printf("None of the above. Got %s\n", m)
	}
}

func LookUp(raw []byte) string {
	var ret string
	b := littleEndianConcat(raw)
	for _, entry := range OpCodeLookUpTable {
		v := b & entry.Mask
		if v == entry.Value {
			switch entry.Name {
			case "std":
				return deConvoluter(b)
			case "ldd":
				return deConvoluter(b)
			}
			return entry.Name
		} else {
			ret = fmt.Sprintf("%.16b (0x%.4x)", b, b)
		}
	}
	return ret
}

func deConvoluter(b uint16) string {
	x := b  & 0xd208
	offset := b & 0x2c07
	switch x {
	case 0x8000:
		if offset == 0 {
			return "ldz"
		} else {
			return fmt.Sprintf("lddz+%d", offset)
		}
	case 0x8008:
		if offset == 0 {
			return "ldy"
		} else {
			return fmt.Sprintf("lddy+%d", offset)
		}
	case 0x8200:
		if offset == 0 {
			return "stz"
		} else {
			return fmt.Sprintf("stdz+%d", offset)
		}
	case 0x8208:
		if offset == 0 {
			return "sty"
		} else {
			return fmt.Sprintf("stdy+%d", offset)
		}
	default:
		return "idunno"
	}
}
