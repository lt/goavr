package main

// Balls

import (
	"debug/elf"
	"flag"
	"fmt"
	//"encoding/hex"
	//"reflect"
)

var fileName string

func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getExecutableStuff(file *elf.File) []byte {
	var x int
	for i, s := range file.Sections {
		if s.SectionHeader.Name == ".text" {
			x = i
		}
	}
	ret, _ := file.Sections[x].Data()
	return ret
}

func chunkle(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, 0)
	x := 0

	for i := 0; i < (len(blob) - csize); i += csize {
		fin = append(fin, []byte(blob[i:(x+csize)]))
		x += csize
	}
	fin = append(fin, []byte(blob[x:]))
	return fin
}

func bigEndianConcat(b []byte) uint16 {
	ret := (uint16(b[0]) << 8) | uint16(b[1])
	return ret
}

func littleEndianConcat(b []byte) uint16 {
	ret := (uint16(b[1]) << 8) | uint16(b[0])
	return ret
}

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

func main() {

	flag.Parse()

	if fileName == "" {
		fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
	}

	file, _ := elf.Open(fileName)
	t := getExecutableStuff(file)

	for i, c := range chunkle(t, 2) {
		fmt.Printf(".text:0x%.8x\t", (i << 1))
		//fmt.Printf(hex.Dump(c))
		dissAssemble(c)
	}
}
