package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
)

var fileName string
var dumpMem, dumpProg bool 



func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
	flag.BoolVar(&dumpMem, "d", false, "Just dump the program memory")
	flag.BoolVar(&dumpProg, "p", false, "Pretty print the whole file")
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

		if dumpProg == true {
			dissectExecutable(file)
			os.Exit(0)
		} else {
			getStuff(file)
			cpu.imem.LoadProgram(data)
		}
	}

	if dumpMem == true {
		fmt.Println(cpu.imem.Dump())
		os.Exit(0)
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

	// Still don't know how to exit the program.
	
	for  cpu.pc != programEnd  {
		fmt.Printf("pc: %.4x\tsr: %.8b\tsp: %.4x\t\n", cpu.pc, cpu.sr, cpu.sp)
		mi := cpu.imem.Fetch()
		cpu.Execute(dissAssemble(mi))
		printRegs(cpu.regs)
		fmt.Println("Stack: ", cpu.dmem[cpu.sp:0x01ff])
		fmt.Println("---------------------------------")
	}
}
