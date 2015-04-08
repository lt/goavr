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

	bitMasks := map[byte]byte{
		1: 1,
		2: 2,
		3: 4,
		4: 8,
		5: 16,
		6: 32,
		7: 64,
		8: 128,
	}
	
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
		r := cpu.regs[i.dest] + cpu.regs[i.source]
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if r > 0xff {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		return
	case INSN_ANDI:
		// Rd <- Rd & K
		cpu.regs[i.dest] = cpu.regs[i.dest] & i.kdata
		return
	case INSN_NOP:
		// Duh.
		return
	case INSN_CLI:
		// Clear global interrupt
		cpu.clear_i()
		return
	case INSN_ADC:
		// Rd <- Rd + Rr + C
		c := uint8(cpu.sr & 0x01)
		r := cpu.regs[i.dest] + cpu.regs[i.source] + c
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if r > 0xff {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
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
		// I/O(A,b) <- 1
		fmt.Printf("%.8b\n", cpu.dmem[i.ioaddr])
		// bitMasks is a map that looks up the mask to isolate
		// the necessary bit
		cpu.dmem[i.ioaddr] &= bitMasks[i.registerBit]
		fmt.Printf("%.8b\n", cpu.dmem[i.ioaddr])
		return
	case INSN_CBI:
		// I/O(A,b) <- 0
		cpu.dmem[i.ioaddr] ^= bitMasks[i.registerBit]
		fmt.Printf("%.8b\n", cpu.dmem[i.ioaddr])
		return
	case INSN_SBIC:
		// If I/O(A,b) = 0 then PC <- PC + 2 (or 3) else PC <- PC + 1
		s := i.registerBit - 1
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> s
		if r == 0 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_SBIS:
		// If I/O(A,b) = 1 then PC <- PC + 2
		s := i.registerBit - 1
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> s
		if r == 1 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_BLD:
		// Rd(b) <- T
		// Copies the T Flag in the SREG (Status Register) to bit b in register Rd.
		t := (cpu.sr & int16(bitMasks[7])) >> 7
		cpu.dmem[i.dest] = byte(t)
		return
	case INSN_BST:
		// T <- Rd(b)
		// Stores bit b from Rd to the T Flag in SREG (Status Register).
		s := i.registerBit - 1
		t := cpu.regs[i.dest] & bitMasks[i.registerBit] >> s
		cpu.sr &= int16(t << 7)
		return
	case INSN_SBRC:
		// if Rr(b) = 0 then PC += 2
		s := i.registerBit - 1
		r := (cpu.regs[i.source] & bitMasks[i.registerBit]) >> s
		if r == 0 {
			cpu.pc += 2
		}
		return
	case INSN_STS:
		// (k) <- Rr
		cpu.dmem[i.k16] = cpu.regs[i.source]
		return
	case INSN_LDS:
		// Rd <- (k)
		cpu.regs[i.dest] = cpu.dmem[i.k16]
		return
	case INSN_ADIW:
		// Rd+1:Rd <- Rd+1:Rd + K
		kl := i.kdata & 0x0f
		kh := (i.kdata & 0xf0) >> 4
		// I think I may have swapped the byte order here.
		// XXX TODO(ERIN)
		cpu.regs[i.dest] = cpu.regs[i.dest] + kl
		cpu.regs[i.dest+1] = cpu.regs[i.dest+1] + kh
		return
	case INSN_SBIW:
		fmt.Println(i)
		// Rd+1:Rd <- Rd+1:Rd - K
		kl := i.kdata & 0x0f
		kh := (i.kdata & 0xf0) >> 4
		// XXX TODO(ERIN) See ADIW
		cpu.regs[i.dest] = cpu.regs[i.dest] - kl
		cpu.regs[i.dest+1] = cpu.regs[i.dest+1] - kh
		r := b2u16little([]byte{cpu.regs[i.dest+1],cpu.regs[i.dest]} )
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_BRCC:
		// Branch if carry cleared
		c := cpu.sr & 0x01
		if c == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRCS:
		// Branch if carry set
		c := cpu.sr & 0x01
		if c == 1 {
			cpu.pc += i.k16
		}
		return
	case INSN_BREQ:
		//if Rd = Rr(Z=1) then PC <- PC + k + 1
		r := (cpu.sr & int16(bitMasks[7])) >> 7
		if r == 1 {
			cpu.pc += i.k16
		}
		return
	case INSN_BRGE:
		// if Rd >= Rr then PC += k
		if cpu.regs[i.dest] >= cpu.regs[i.source] {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRNE:
		// if (Z = 0) then PC <-  PC + k + 1
		z := (cpu.sr & 0x02) >> 1
		if z == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRTC:
		// if T = 0 then PC <- PC + k + 1
		t := (cpu.sr & 64) >> 6
		if t == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_COM:
		// Rd <- ^Rd
		cpu.regs[i.dest] = ^cpu.regs[i.dest]
		return
	case INSN_CP:
		// Rd - Rr
		if cpu.regs[i.source] > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.regs[i.dest] - cpu.regs[i.source]
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_CPC:
		// Rd - Rr - C
		c := uint8(cpu.sr & 0x01)
		if (cpu.regs[i.source] + c) > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.regs[i.dest] - cpu.regs[i.source] - c
		if r != 0 {
			cpu.clear_z()
		}
		return
	case INSN_CPI:
		// Rd - K
		// I can't tell from the doc, but I think this check
		// has to happen before K is subtracted. We'll see.
		if i.kdata > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.regs[i.dest] - i.kdata

		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}

		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_CPSE:
		// if Rd = Rr then PC <- PC + 2
		if cpu.regs[i.dest] == cpu.regs[i.source] {
			cpu.pc += 2
		}
		return
	case INSN_DEC:
		// Rd <- Rd - 1
		r := cpu.regs[i.dest] - 1
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_IN:
		// Rd <- I/O(A)
		cpu.regs[i.dest] = cpu.dmem[i.ioaddr]
		return
	case INSN_LDDY:
		// Rd <- (Y + q)
		y := b2u16little([]byte{cpu.regs[29], cpu.regs[28]})
		cpu.regs[i.dest] = cpu.dmem[y + i.offset]
		return
	case INSN_LDDZ:
		// Rd <- (Z + q)
		z := b2u16little([]byte{cpu.regs[31], cpu.regs[30]})
		cpu.regs[i.dest] = cpu.dmem[z + i.offset]
		return
	case INSN_LDX:
		// Rd <- (X)
		x := b2u16little([]byte{cpu.regs[27], cpu.regs[26]})
		cpu.regs[i.dest] = cpu.dmem[x]
		return
	case INSN_LDXP:
		// Rd <- (X), X <- X + 1
		x := b2u16little([]byte{cpu.regs[27], cpu.regs[26]})
		cpu.regs[i.dest] = cpu.dmem[x]
		cpu.regs[26] += 1
		return
	case INSN_LDY:
		// Rd <- (Y)
		y := b2u16little([]byte{cpu.regs[29], cpu.regs[28]})
		cpu.regs[i.dest] = cpu.dmem[y]
		return
	case INSN_LDZ:
		// Rd <- (Z) (dmem)
		z := b2u16little([]byte{cpu.regs[31], cpu.regs[30]})
		cpu.regs[i.dest] = cpu.dmem[z]
		return
	case INSN_LPMZ:
		// Rd <- (Z) (imem)
		z := b2u16little([]byte{cpu.regs[31], cpu.regs[30]})
		cpu.regs[i.dest] = cpu.imem[z]
		return
	case INSN_LPMZP: 
		// Rd <- (Z), Z <- Z + 1 (imem)
		z := (cpu.regs[31] << 8) | cpu.regs[30]
		cpu.regs[i.dest] = cpu.imem[z]
		cpu.regs[30] += 1
		return
	case INSN_LPM:
		// R0 <- (Z)
		z := (cpu.regs[31] << 8) | cpu.regs[30]
		cpu.regs[0] = cpu.imem[z]
		return
	case INSN_LSR:
		// logical shift right Rd
		r := cpu.regs[i.dest] >> 1
		cpu.regs[i.dest] = r
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
		r := cpu.regs[i.dest] - i.kdata
		cpu.regs[i.dest] = r
		if r != 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if i.kdata > r {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		return
	case INSN_SBCI:
		// Rd <- Rd - K - C
		c := uint8(cpu.sr & 0x01)
		r := cpu.regs[i.dest] - i.kdata - c
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.clear_z()
		}
		if (i.kdata + c) > r {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		return
	case INSN_SEI:
		// set global interrupt flag
		cpu.set_i()
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
		return
	}
}
