package main

import (
	"fmt"
)

type CPU struct {
	regs [32]uint8
	pc   uint16
	sp   uint16
	sr   uint8
}

type Memory [0x045f]byte

func (cpu *CPU) Execute(i Instr) {
	fmt.Println(i.label)
}

func (cpu *CPU) Pc() {

}
