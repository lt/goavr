package main

import (
	"fmt"
)

type CPU struct {
	regs [32]uint8
	pc   int16
	sp   uint16
	sr   uint8
}

type Memory [0x045f]byte

func (cpu *CPU) Execute(i Instr) {
	switch i.label {
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		//cpu.pc = i.k32
		return
	case INSN_RJMP:
		cpu.pc = cpu.pc + i.k16 + 1
		return
	default:
		fmt.Println("I dunno.")
	}
}

func (cpu *CPU) Pc() {

}
