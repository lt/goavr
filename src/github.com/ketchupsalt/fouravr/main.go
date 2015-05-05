package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var fileName, single string
var dumpMem, dumpProg bool

func init() {
	flag.StringVar(&single, "s", single, "Parse a single instruction")
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
	flag.BoolVar(&dumpMem, "d", false, "Just dump the program memory")
	flag.BoolVar(&dumpProg, "p", false, "Pretty print the whole file")
}

var cpu CPU

func main() {

	flag.Parse()

	if single != "" {
		x, _ := strconv.Atoi(single)
		fmt.Printf("%.4x\n", x)
		op := lookUp(u16lil2byte(uint16(x)))
		fmt.Println(op)
		os.Exit(0)
	}

	if fileName == "" {
		//fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
		// JMP
		// d := []byte{0x0c, 0x94, 0xc5, 0xbb}
		// ANDI
		cpu.imem.LoadProgram([]byte{0x27, 0x7f})
	} else {
		r, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			fmt.Printf("File %s not found.\n", fileName)
			os.Exit(2)
		} else if r.IsDir() {
			fmt.Printf("%s is a directory, dummy.\n", fileName)
			os.Exit(2)
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
	cpu.pc = 0

	// Still don't know how to exit the program.

	for cpu.pc != programEnd {
		cpu.Interactive()
	}
}
