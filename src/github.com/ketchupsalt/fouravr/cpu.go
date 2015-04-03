package main

import (
	"fmt"
)

// AVR indirect pointer registers.
// X = (26,27), Y = (28,29, and Z = (30,31)
// XXX TODO(Erin) may be a better way to do this?
// Like maybe make it part of the CPU struct?

const (
	XReg = 26
	YReg = 28
	Zreg = 30
)

type CPU struct {
	regs [32]uint8
	pc   int16
	sp   int16
	sr   int16
	imem Memory
	dmem Memory
}

// Set bits in status register
func (cpu *CPU) set_i() { cpu.sr |= 128 }
func (cpu *CPU) set_t() { cpu.sr |= 64 }
func (cpu *CPU) set_h() { cpu.sr |= 32 }
func (cpu *CPU) set_s() { cpu.sr |= 16 }
func (cpu *CPU) set_v() { cpu.sr |= 8 }
func (cpu *CPU) set_n() { cpu.sr |= 4 }
func (cpu *CPU) set_z() { cpu.sr |= 2 }
func (cpu *CPU) set_c() { cpu.sr |= 1 }

// Clear bits in status regsiter
func (cpu *CPU) clear_i() { cpu.sr &= ^128 }
func (cpu *CPU) clear_t() { cpu.sr &= ^64 }
func (cpu *CPU) clear_h() { cpu.sr &= ^32 }
func (cpu *CPU) clear_s() { cpu.sr &= ^16 }
func (cpu *CPU) clear_v() { cpu.sr &= ^8 }
func (cpu *CPU) clear_n() { cpu.sr &= ^4 }
func (cpu *CPU) clear_z() { cpu.sr &= ^2 }
func (cpu *CPU) clear_c() { cpu.sr &= ^1 }

/*
Golang Logical Operators: (because I'm tired of looking this shit up)
+    ADD
-    SUB
&    bitwise AND
|    bitwise OR
^    bitwise XOR
&^   bit clear (AND NOT)

<<   left shift
>>   right shift
*/

func (cpu *CPU) Execute(i Instr) {
	fmt.Printf("pc: %.4x\tsr: %.8b\tsp: %.4x\t\n", cpu.pc, cpu.sr, cpu.sp)
	switch i.label {
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		//cpu.pc = i.k32
		return
	case INSN_RJMP:
		// PC <- PC + k + 1
		cpu.pc = cpu.pc + i.k16
		return
	case INSN_ADD:
		// Rd <- Rd + Rr
		cpu.regs[i.dest] = cpu.regs[i.dest] + cpu.regs[i.source]
		return
	case INSN_ANDI:
		// Rd <- Rd & K
		cpu.regs[i.dest] = cpu.regs[i.dest] & i.kdata
		return
	case INSN_NOP:
		return
	case INSN_CLI:
		cpu.sr = 7
		return
	case INSN_ADC:
		// Rd <- Rd + Rr + C
		c := uint8(cpu.sr & 0x01)
		cpu.regs[i.dest] = cpu.regs[i.dest] + cpu.regs[i.source] + c
		return
	case INSN_EOR:
		// Rd <- Rd^Rr
		cpu.regs[i.dest] = cpu.regs[i.dest] ^ cpu.regs[i.source]
		return
	case INSN_OUT:
		// I/O port <- Rr
		return
	case INSN_LDI:
		// Rd <- K
		cpu.regs[i.dest] = i.kdata
		return
	case INSN_RCALL:
		// PC <- PC + k + 1
		// says +1, but that generates the wrong value.
		cpu.pc = cpu.pc + i.k16 // + 1
		return
	case INSN_SBI:

		return
	case INSN_CBI:

		return
	case INSN_SBIC:

		return
	case INSN_SBIS:

		return
	case INSN_BLD:

		return
	case INSN_BST:

		return
	case INSN_SBRC:

		return
	case INSN_STS:

		return
	case INSN_LDS:

		return
	case INSN_ADIW:

		return
	case INSN_SBIW:

		return
	case INSN_BRCC:

		return
	case INSN_BRCS:

		return
	case INSN_BREQ:

		return
	case INSN_BRGE:

		return
	case INSN_BRNE:
		// if (Z = 0) then PC <-  PC + k + 1
		if (cpu.sr & 0x02) == 0 {
			cpu.pc += i.k16
		}
		return
	case INSN_BRTC:

		return
	case INSN_COM:
		// Rd <- ^Rd
		cpu.regs[i.dest] = ^cpu.regs[i.dest]
		return
	case INSN_CP:

		return
	case INSN_CPC:
		// Rd - Rr - C
		c := byte(cpu.sr & 0x01)
		r := cpu.regs[i.dest] - cpu.regs[i.source] - c
		if r != 0 {
			cpu.clear_z()
		}
		if (cpu.regs[i.source] + c) > cpu.regs[i.dest] {
			cpu.clear_c()
		} else {
			cpu.set_c()
		}
		return
	case INSN_CPI:
		// Rd - K
		if (cpu.regs[i.dest] - i.kdata) == 0 {
			cpu.set_z()
		}
		if i.kdata > cpu.regs[i.dest] {
			cpu.set_c()
		}
		return
	case INSN_CPSE:

		return
	case INSN_DEC:
		// Rd <- Rd - 1
		cpu.regs[i.dest] -= 1
		return
	case INSN_IN:

		return
	case INSN_LDDY:

		return
	case INSN_LDDZ:

		return
	case INSN_LDX:

		return
	case INSN_LDXP:

		return
	case INSN_LDY:

		return
	case INSN_LDZ:

		return
	case INSN_LPMZ:

		return
	case INSN_LPMZP:
		// Rd <- (Z), Z <- Z + 1
		z := (cpu.regs[31] << 8) | cpu.regs[30]
		cpu.regs[i.dest] = cpu.imem[z]
		cpu.regs[30] += 1
		return
	case INSN_LPM:

		return
	case INSN_LSR:

		return
	case INSN_MOV:
		// Rd <- Rr
		cpu.regs[i.dest] = cpu.regs[i.source]
		return
	case INSN_MOVW:
		// Rd+1:Rd <- Rr+1:Rr
		cpu.regs[i.dest+1] = cpu.regs[i.source+1]
		cpu.regs[i.dest] = cpu.regs[i.source]
		return
	case INSN_MUL:

		return
	case INSN_NEG:

		return
	case INSN_OR:

		return
	case INSN_ORI:

		return
	case INSN_POP:

		return
	case INSN_PUSH:
		// STACK <- Rr
		cpu.sp -= 1
		return
	case INSN_RET:

		return
	case INSN_RETI:

		return
	case INSN_ROR:

		return
	case INSN_SBC:

		return
	case INSN_SUBI:
		// Rd <- Rd - K
		cpu.regs[i.dest] = cpu.regs[i.dest] - i.kdata

		return
	case INSN_SBCI:
		// Rd <- Rd - K - C
		c := cpu.sr & 0x01
		cpu.regs[i.dest] = cpu.regs[i.dest] - i.kdata - uint8(c)
		return
	case INSN_SEI:

		return
	case INSN_STDY:

		return
	case INSN_STDZ:

		return
	case INSN_STX:

		return
	case INSN_STXP:
		// (X) <- Rr, X <- X + 1
		// 26 = low byte, 27 = high byte
		x := int16(cpu.regs[27])<<8 | int16(cpu.regs[26])
		fmt.Printf("dmem: %.4x\t", x)
		cpu.dmem[x] = cpu.regs[i.source]
		cpu.regs[26] += 1

		return
	case INSN_STXM:

		return
	case INSN_STY:

		return
	case INSN_STZ:

		return
	case INSN_STZP:

		return
	case INSN_STZM:

		return
	case INSN_SUB:
		// Rd <- Rd - Rr
		cpu.regs[i.dest] = cpu.regs[i.dest] - cpu.regs[i.source]
		return
	default:
		fmt.Println("I dunno.")
	}
}
