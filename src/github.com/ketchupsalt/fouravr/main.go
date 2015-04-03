package main

import (
	"debug/elf"
	"flag"
	"fmt"
//	"os"
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
		cpu.imem.LoadProgram([]byte{0x27, 0x7f})		
	} else {
		file, _ := elf.Open(fileName)
		getExecutableStuff(file)
		cpu.imem.LoadProgram(data)
	}

	//mi := dissAssemble(pop(2))
	// sample JMP instruction 94 0c bb c5
	// 1001 0100 0000 1100 1011 1011 1100 0101
	//mi := Instr{label: INSN_JMP, family: Branches, }


	// Manually setting the program counter to the start of
	// the stuff I actually want to step through.
	cpu.pc = 0x0026
	
	// RAMEND is typically 0x1ff. Compiler leaves a preamble
	// in the decompiled code that is supposed to initialize
	// the stack pointer, but I can't yet figure out where,
	// so I am setting this manually.
	cpu.sp = 0x1ff
	

	// I made the CLI instruction set the status register to 7
	// which is the global interrupt enable. This ends the
	// program, for now, even though there's an rjmp that jumps
	// back to it before exiting.

	for cpu.sr != 7 {
		fmt.Println("")
		mi := cpu.imem.Fetch()
		cpu.Execute(dissAssemble(mi))
		printRegs(cpu.regs)
	}
}
