package main

import (
	"fmt"
)

type CPU struct {
	regs [32]uint8
	pc   int16
	sp   uint16
	sr   uint8
	memory Memory
}

// Actual AVR is 1024 bytes, but my test program is 1800ish.

type Memory [2048]byte

func printMnemonic(label int) {
	ret := fmt.Sprintf("I am %d\n", label)
	for _, op := range(OpCodeLookUpTable) {
		if op.label == label {
			ret = op.mnemonic
		}
	}
	fmt.Println(ret)
}

func (cpu *CPU) Execute(i Instr) {
	fmt.Printf("%.4x:\t", cpu.pc)
	switch i.label {
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		//cpu.pc = i.k32
		return
	case INSN_RJMP:
		// PC <- PC + k + 1
		printMnemonic(i.label)
		//cpu.pc = cpu.pc + i.k16 + 1
		//fmt.Println(i.k16)
		return
	case INSN_ADD:
		// Rd <- Rd + Rr
		printMnemonic(i.label)
		cpu.regs[i.dest] = cpu.regs[i.dest] + cpu.regs[i.source]
		return
	case INSN_ANDI:
		// Rd <- Rd & K
		printMnemonic(i.label)
		cpu.regs[i.dest] = cpu.regs[i.dest] & i.kdata
		return
	case INSN_NOP:
		printMnemonic(i.label)
		return
	case INSN_CLI:
		cpu.sr = 7
		printMnemonic(i.label)
		return
	case INSN_ADC:
		printMnemonic(i.label)
		return
	case INSN_EOR:
		printMnemonic(i.label)
		return
	case INSN_OUT:
		printMnemonic(i.label)
		return
	case INSN_LDI:
		// Rd <- K
		cpu.regs[i.dest] = i.kdata
		printMnemonic(i.label)
		return
	case INSN_RCALL:
		printMnemonic(i.label)
		return
	case INSN_SBI:
		printMnemonic(i.label)
		return
	case INSN_CBI:
		printMnemonic(i.label)
		return
	case INSN_SBIC:
		printMnemonic(i.label)
		return
	case INSN_SBIS:
		printMnemonic(i.label)
		return
	case INSN_BLD:
		printMnemonic(i.label)
		return
	case INSN_BST:
		printMnemonic(i.label)
		return
	case INSN_SBRC:
		printMnemonic(i.label)
		return
	case INSN_STS:
		printMnemonic(i.label)
		return
	case INSN_LDS:
		printMnemonic(i.label)
		return
	case INSN_ADIW:
		printMnemonic(i.label)
		return
	case INSN_SBIW:
		printMnemonic(i.label)
		return
	case INSN_BRCC:
		printMnemonic(i.label)
		return
	case INSN_BRCS:
		printMnemonic(i.label)
		return
	case INSN_BREQ:
		printMnemonic(i.label)
		return
	case INSN_BRGE:
		printMnemonic(i.label)
		return
	case INSN_BRNE:
		printMnemonic(i.label)
		return
	case INSN_BRTC:
		printMnemonic(i.label)
		return
	case INSN_COM:
		printMnemonic(i.label)
		return
	case INSN_CP:
		printMnemonic(i.label)
		return
	case INSN_CPC:
		printMnemonic(i.label)
		return
	case INSN_CPI:
		printMnemonic(i.label)
		return
	case INSN_CPSE:
		printMnemonic(i.label)
		return
	case INSN_DEC:
		printMnemonic(i.label)
		return
	case INSN_IN:
		printMnemonic(i.label)
		return
	case INSN_LDDY:
		printMnemonic(i.label)
		return
	case INSN_LDDZ:
		printMnemonic(i.label)
		return
	case INSN_LDX:
		printMnemonic(i.label)
		return
	case INSN_LDXP:
		printMnemonic(i.label)
		return
	case INSN_LDY:
		printMnemonic(i.label)
		return
	case INSN_LDZ:
		printMnemonic(i.label)
		return
	case INSN_LPMZ:
		printMnemonic(i.label)
		return
	case INSN_LPM:
		printMnemonic(i.label)
		return
	case INSN_LSR:
		printMnemonic(i.label)
		return
	case INSN_MOV:
		printMnemonic(i.label)
		return
	case INSN_MOVW:
		printMnemonic(i.label)
		return
	case INSN_MUL:
		printMnemonic(i.label)
		return
	case INSN_NEG:
		printMnemonic(i.label)
		return
	case INSN_OR:
		printMnemonic(i.label)
		return
	case INSN_ORI:
		printMnemonic(i.label)
		return
	case INSN_POP:
		printMnemonic(i.label)
		return
	case INSN_PUSH:
		printMnemonic(i.label)
		return
	case INSN_RET:
		printMnemonic(i.label)
		return
	case INSN_RETI:
		printMnemonic(i.label)
		return
	case INSN_ROR:
		printMnemonic(i.label)
		return
	case INSN_SBC:
		printMnemonic(i.label)
		return
	case INSN_SBCI:
		printMnemonic(i.label)
		return
	case INSN_SEI:
		printMnemonic(i.label)
		return
	case INSN_STDY:
		printMnemonic(i.label)
		return
	case INSN_STDZ:
		printMnemonic(i.label)
		return
	case INSN_STX:
		printMnemonic(i.label)
		return
	case INSN_STXP:
		printMnemonic(i.label)
		return
	case INSN_STXM:
		printMnemonic(i.label)
		return
	case INSN_STY:
		printMnemonic(i.label)
		return
	case INSN_STZ:
		printMnemonic(i.label)
		return
	case INSN_STZP:
		printMnemonic(i.label)
		return
	case INSN_STZM:
		printMnemonic(i.label)
		return
	case INSN_SUB:
		printMnemonic(i.label)
		return
	case INSN_SUBI:
		printMnemonic(i.label)
		return
	default:
		fmt.Println("I dunno.")
	}
}

func (mem *Memory) Fetch(i int16) []byte {
	ret := mem[cpu.pc:(cpu.pc+i)]
	cpu.pc += i
	return ret
}

func (mem *Memory) LoadProgram(data []byte) {
	for i, b := range(data) {
		mem[i] = b
	}
}
