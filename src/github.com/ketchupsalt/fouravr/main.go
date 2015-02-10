package main

// Balls

import (
	"fmt"
	"flag"
	"debug/elf"
	"encoding/hex"
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

func thingsHm(b []byte) string {
	// the second byte passed in contains the opcode (for now)
	// mask off all but the first four bits see if we have a match
	// This doesn't feel right to me.
	m := int(b[1]) & 0xf0
	fmt.Println(m, ": ", b[1])
	switch  m {
	case EOR:
		fmt.Println("EOR")
	case OUT:
		fmt.Println("OUT")
	case RJMP:
		fmt.Println("RJMP")
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
		thingsHm(c)
		fmt.Println(hex.Dump(c))
	}
}
