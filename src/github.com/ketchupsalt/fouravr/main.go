package main

import (
	"debug/elf"
	"flag"
//	"fmt"
)

var fileName string
var cSize int = 2
var pc = 0
var data []byte

func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
}

func main() {

	flag.Parse()

	if fileName == "" {
		//fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
		// JMP
		// d := []byte{0x0c, 0x94, 0xc5, 0xbb}
		// RJMP
		d := []byte{0x27, 0xc0}
		for _, z := range(d) {
			//fmt.Printf("%.8b\n", z)
			data = append(data, z)
		}
		
	} else {
		file, _ := elf.Open(fileName)
		getExecutableStuff(file)
	}

	//mi := dissAssemble(pop(2))
	// sample JMP instruction 94 0c bb c5
	// 1001 0100 0000 1100 1011 1011 1100 0101
	//mi := Instr{label: INSN_JMP, family: Branches, }
	mi := dissAssemble(pop(2))
	cpu := CPU{}
	cpu.pc = 0
	cpu.Execute(mi)
	
/*
	for len(data) > 0 {
		fmt.Printf(".text:%.4x:\t", (pc))
		f := dissAssemble(pop(2))
		fmt.Println(f)
	}
*/
}
