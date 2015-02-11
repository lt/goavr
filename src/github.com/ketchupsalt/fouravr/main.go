package main

// Balls

import (
	"fmt"
	"flag"
	"debug/elf"
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
	for i, s := range(file.Sections) {
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

	for i :=0; i < (len(blob) - csize); i+=csize {
		fin = append(fin, []byte(blob[i:(x+csize)]))
		x += csize
	}
	fin = append(fin, []byte(blob[x:]))
	return fin
}

func bigEndianConcat(b []byte) uint16 {
	ret := (uint16(b[0]) << 8 ) | uint16(b[1])
	return ret
}


func littleEndianConcat(b []byte) uint16 {
	ret := (uint16(b[1]) << 8 ) | uint16(b[0])
	return ret
}


func dissAssemble(b []byte) string {
	// the second byte passed in contains the opcode
	// (for now -- I might pass them in in little endian order later)
	// mask off all but the first four bits see if we have a match
	// This doesn't feel right to me.
	m := b[1] & 0xf0
	//fmt.Println(m, ": ", b[1])
	switch  m {
	case EOR:
		r := b[1] & 0x02
		d := b[1] & 0x01

		rrrr := b[0] & 0x0f
		dddd := (b[0] >> 4) & 0x0f
		
		fmt.Printf("EOR\t%b %b\tr%d,r%d\n", r, d, rrrr, dddd)
	case OUT:
		//out := (b[1] >> 3) & 0xff
		AA := (b[1]  & 0x06) >> 1
		i := littleEndianConcat(b)
		register := (i & 0x01f0) >> 4
		address := AA << 4 | (b[0] & 0x0f)
		fmt.Printf("out\t0x%.2x,r%d\t;%d\n", address, register, address)
	case RJMP:
		i := littleEndianConcat(b)
		// I DON'T KNOW WHY THIS WORKS!!!
		// left shift 1 produces the same output objdump does.
		// find out why?
		k := (i & 0x0fff) << 1
		fmt.Printf("RJMP\t.+%d\n", k)
	case CLI:
		fmt.Println("CLI")
	case LDI:
		fmt.Println("LDI")
	case RCALL:
		fmt.Println("RCALL")
	default:
		fmt.Println("None of the above")
	}
	return "Test"
}

func main() {

	if fileName == "" {
		fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
	}
	
	file, _ := elf.Open(fileName)
	t := getExecutableStuff(file)
	//fmt.Println(hex.Dump(t))

	//fmt.Println(chunkle(t,2))
	
	for _, c := range(chunkle(t, 2)) {
		//fmt.Printf(hex.Dump(c))
		dissAssemble(c)
	}
}
