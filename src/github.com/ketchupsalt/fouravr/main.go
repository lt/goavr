package main

import (
	"debug/elf"
	"flag"
	"fmt"
)

var fileName string
var cSize int = 2
var pc = 0
var data []byte


func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
}

var cpu CPU

func main() {

	flag.Parse()
	if fileName == "" {
		//fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
		// JMP
		// d := []byte{0x0c, 0x94, 0xc5, 0xbb}
		// ANDI
		cpu.memory.LoadProgram([]byte{0x27, 0x7f})		
	} else {
		file, _ := elf.Open(fileName)
		getExecutableStuff(file)
		cpu.memory.LoadProgram(data)
	}

	//mi := dissAssemble(pop(2))
	// sample JMP instruction 94 0c bb c5
	// 1001 0100 0000 1100 1011 1011 1100 0101
	//mi := Instr{label: INSN_JMP, family: Branches, }

	cpu.pc = 0
	
	for x := 0; x < 500; x++ {
		fmt.Printf("Before:\t%v\n", cpu.regs)
		mi := cpu.memory.Fetch(2)
		//fmt.Println(dissAssemble(mi))
		cpu.Execute(dissAssemble(mi))
		fmt.Printf("After:\t%v\n", cpu.regs)
	}

}
