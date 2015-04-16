package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (cpu *CPU) Step() {
	fmt.Printf("pc: %.4x\tsr: %.8b\tsp: %.4x\t\n", cpu.pc, cpu.sr, cpu.sp)
	mi := cpu.imem.Fetch()
	cpu.Execute(dissAssemble(mi))
	printRegs(cpu.regs)
	fmt.Println("Stack: ", cpu.dmem[cpu.sp:0x01ff])
	fmt.Println("---------------------------------")
}

func (cpu *CPU) Interactive() {
		fmt.Println("Type ? for help.")
	for {
		prompt := bufio.NewReader(os.Stdin)
		fmt.Print("$> ")

		response, err := prompt.ReadString('\n')

		check(err)
		// Ugh.
		r := strings.Split(response, "\n")

		switch r[0] {
		case "?":
			fmt.Println("g to run the whole program")
			fmt.Println("q to quit")
			fmt.Println("s to single step")
			fmt.Println("j prompts for a pc (in hex)")
			fmt.Println("d dumps the data memory")
			fmt.Println("p dumps the program mempory")
			fmt.Println("return jumps 5 instructions")
			fmt.Println("any number /n/ jumps /n/ instructions")
		case "g":
			for {
				cpu.Step()
				if cpu.pc == programEnd {
					break
				}
			}
		case "q":
			os.Exit(0)
		case "s":
			var b string
			for {
				cpu.Step()
				fmt.Scanf("%s", &b)
				if b == "x" {
					break
				}
				if cpu.pc == programEnd { break }
			}
		case "b":
			cpu.pc -= 2
		case "r":
			cpu.pc = 0x0026
		case "d":
			fmt.Println(cpu.dmem.Dump())
		case "p":
			fmt.Println(cpu.imem.Dump())
		case "j":
			var o int16
			fmt.Println("Enter pc:")
			fmt.Scanf("%x", &o)
			for cpu.pc < (o + 2) {
				cpu.Step()
				if cpu.pc == programEnd { break }
			}
		default:
			var n int
			// default case: step n times
			if r[0] == "" {
				n = 5
			} else {
				n, err = strconv.Atoi(r[0])
				if err != nil {
					fmt.Println("Command not recognized.")
					break
				}
			}
			for i := 1; i < (n + 1); i++ {
				cpu.Step()
				if cpu.pc == programEnd {
					break
				}
			}
		}
	}
}

func (cpu *CPU) Execute(i Instr) {

	bitMasks := []byte{0, 1, 2, 4, 8, 16, 32, 64, 128}

	switch i.label {
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		//cpu.pc = i.k32
		return
	case INSN_IJMP:
		// PC <- Z(15:0)
		z := b2i16little([]byte{cpu.regs[30], cpu.regs[31]})
		cpu.pc = z - 1
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
	case INSN_AND:
		// Rd <- Rd & Rr
		r := cpu.regs[i.dest] & cpu.regs[i.source]
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x80 >> 8) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
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
	case INSN_IN:
		// Rd <- I/O(A)
		cpu.regs[i.dest] = cpu.dmem[i.ioaddr]
		return
	case INSN_OUT:
		// I/O(A) <- Rr
		cpu.dmem[i.ioaddr] = cpu.regs[i.source]
		return
	case INSN_LDI:
		// Rd <- K
		cpu.regs[i.dest] = i.kdata
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
	case INSN_SBRS:
		// if Rr(b) = 1 then PC += 2
		s := i.registerBit - 1
		r := (cpu.regs[i.source] & bitMasks[i.registerBit]) >> s
		if r == 1 {
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
		// low byte
		x := uint16(cpu.regs[i.dest])
		// high byte
		y := uint16(cpu.regs[i.dest+1])
		r := ((y << 8) | x) + uint16(i.kdata)
		cpu.regs[i.dest] = uint8(r & 0x00ff)
		cpu.regs[i.dest+1] = uint8(r >> 8)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r >> 15) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SBIW:
		// Rd+1:Rd <- Rd+1:Rd - K
		if i.kdata > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		// low byte
		x := uint16(cpu.regs[i.dest])
		// high byte
		y := uint16(cpu.regs[i.dest+1])
		r := ((y << 8) | x) - uint16(i.kdata)
		cpu.regs[i.dest] = uint8(r & 0x00ff)
		cpu.regs[i.dest+1] = uint8(r >> 8)
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r >> 15) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
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
		r := ^cpu.regs[i.dest]
		cpu.regs[i.dest] = r
		cpu.clear_v()
		cpu.set_c()
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
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
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
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
	case INSN_LDDY:
		// Rd <- (Y + q)
		y := b2u16little([]byte{cpu.regs[28], cpu.regs[29]})
		cpu.regs[i.dest] = cpu.dmem[y+i.offset]
		return
	case INSN_LDDZ:
		// Rd <- (Z + q)
		z := b2u16little([]byte{cpu.regs[30], cpu.regs[31]})
		cpu.regs[i.dest] = cpu.dmem[z+i.offset]
		return
	case INSN_LDX:
		// Rd <- (X)
		x := b2u16little([]byte{cpu.regs[26], cpu.regs[27]})
		cpu.regs[i.dest] = cpu.dmem[x]
		return
	case INSN_LDXP:
		// Rd <- (X), X <- X + 1
		x := b2u16little([]byte{cpu.regs[26], cpu.regs[27]})
		cpu.regs[i.dest] = cpu.dmem[x]
		cpu.regs[26] += 1
		return
	case INSN_LDY:
		// Rd <- (Y)
		y := b2u16little([]byte{cpu.regs[28], cpu.regs[29]})
		cpu.regs[i.dest] = cpu.dmem[y]
		return
	case INSN_LDYP:
		// Rd <- (Y), Y <- Y + 1
		y := b2u16little([]byte{cpu.regs[28], cpu.regs[29]})
		cpu.regs[i.dest] = cpu.dmem[y]
		// XXX TODO(ERIN) this could overflow into the high
		// byte someday.
		cpu.regs[28] += 1
		return
	case INSN_LDYM:
		// Rd <- (Y), Y <- Y - 1
		y := b2u16little([]byte{cpu.regs[28], cpu.regs[29]})
		// pre-decrement
		cpu.regs[28] -= 1
		cpu.regs[i.dest] = cpu.dmem[y]
		return
	case INSN_LDZ:
		// Rd <- (Z) (dmem)
		z := b2u16little([]byte{cpu.regs[30], cpu.regs[31]})
		cpu.regs[i.dest] = cpu.dmem[z]
		return
	case INSN_LDZP:
		// Rd <- (Z) (dmem), Z <- Z - 1
		z := b2u16little([]byte{cpu.regs[30], cpu.regs[31]})
		cpu.regs[i.dest] = cpu.dmem[z]
		// post-decrement
		// XXX TODO(ERIN) this could overflow into the high
		// byte some
		cpu.regs[30] += 1
		return
	case INSN_LDZM:
		// Rd <- (Z) (dmem), Z <- Z - 1
		z := b2u16little([]byte{cpu.regs[30], cpu.regs[31]})
		// pre-decrement
		cpu.regs[30] -= 1
		cpu.regs[i.dest] = cpu.dmem[z]
		return
	case INSN_LPMZ:
		// Rd <- (Z) (imem)
		z := b2i16little([]byte{cpu.regs[30], cpu.regs[31]})
		cpu.regs[i.dest] = cpu.imem[z]
		return
	case INSN_LPMZP:
		// Rd <- (Z), Z <- Z + 1 (imem)
		//z := int16(cpu.regs[31] << 8) | int16(cpu.regs[30])
		z := b2i16little([]byte{cpu.regs[30], cpu.regs[31]})
		fmt.Printf("char found at %.4x:\t%.4x\n", z, cpu.imem[z])
		cpu.regs[i.dest] = cpu.imem[z]
		// post-increment
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
	case INSN_NEG:
		// Replaces the contents of register Rd with its two's complement
		r := ^cpu.regs[i.dest] + 1
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
			cpu.clear_c()
		} else {
			cpu.clear_z()
			cpu.set_c()
		}
		if ((r & 128) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_OR:
		// Rd <- Rd | Rr
		r := cpu.regs[i.dest] | cpu.regs[i.source]
		cpu.regs[i.dest] = r
		cpu.clear_v()
		if ((r & 12) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_ORI:
		// Rd <- Rd | K
		r := cpu.regs[i.dest] | uint8(i.kdata)
		cpu.regs[i.dest] = r
		if ((r & 12) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_POP:
		// Rd <- Stack
		cpu.regs[i.dest] = cpu.dmem[cpu.sp]
		cpu.sp += 1
		return
	case INSN_PUSH:
		// STACK <- Rr
		cpu.dmem[cpu.sp-1] = cpu.regs[i.source]
		cpu.sp -= 1
		return
	case INSN_RCALL:
		// PC <- PC + k + 1, STACK <- PC + 1, SP - 2
		// push the current PC onto the stack because
		// it is automaticaly incremented elsewhere.
		// low byte
		cpu.dmem[cpu.sp-1] = byte(cpu.pc & 0x00ff)
		// high byte
		cpu.dmem[cpu.sp-2] = byte(cpu.pc >> 8)
		// says +1, but that generates the wrong value
		// because the PC is incremented automaticaly anyway
		cpu.pc = cpu.pc + i.k16 //+ 1
		cpu.sp -= 2
		return
	case INSN_RET:
		// PC <- Stack
		low := cpu.dmem[cpu.sp]
		high := cpu.dmem[cpu.sp+1]
		cpu.pc = b2i16little([]byte{high, low})
		cpu.sp += 2
		return
	case INSN_RETI:
		// PC <- Stack, enable interrupts
		low := cpu.dmem[cpu.sp-1]
		high := cpu.dmem[cpu.sp]
		cpu.pc = b2i16little([]byte{high, low})
		cpu.sp += 2
		cpu.set_i()
		return
	case INSN_ROR:
		// Shifts all bits in Rd one place to the right.
		// The C Flag is shifted into bit 7 of Rd.
		// Bit 0 is shifted into the C Flag.
		c := (cpu.regs[i.dest] & 0x80) >> 7
		r := (cpu.regs[i.dest] >> 1) | uint8((cpu.sr&0x01)<<7)
		cpu.regs[i.dest] = r
		if c == 0 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
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
	case INSN_SBC:
		// Rd = Rd - Rr - C
		c := uint8(cpu.sr & 0x01)
		if (cpu.regs[i.source] + c) > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.regs[i.dest] - cpu.regs[i.source] - c
		cpu.regs[i.dest] = r
		if r != 0 {
			cpu.clear_z()
		}
		if ((r & 0x80) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SUBI:
		// Rd <- Rd - K
		r := cpu.regs[i.dest] - i.kdata
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if i.kdata > r {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if (r & 0x80 >> 8) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SBCI:
		// Rd <- Rd - K - C
		c := uint8(cpu.sr & 0x01)
		if (i.kdata + c) > cpu.regs[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.regs[i.dest] - i.kdata - c
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.clear_z()
		}
		if (r & 0x80 >> 8) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SEI:
		// set global interrupt flag
		cpu.set_i()
		return
	case INSN_SUB:
		// Rd <- Rd - Rr
		r := cpu.regs[i.dest] - cpu.regs[i.source]
		cpu.regs[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if i.kdata > r {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if (r & 0x80 >> 8) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_STXP:
		// (X) <- Rr, X <- X + 1
		// 26 = low byte, 27 = high byte
		x := b2i16little([]byte{cpu.regs[26], cpu.regs[27]})
		cpu.dmem[x] = cpu.regs[i.source]
		cpu.regs[26] += 1
		return
	case INSN_MUL:
		// R1h:R0l <- Rx x Rr
		r := uint16(cpu.regs[i.dest]) * uint16(cpu.regs[i.source])
		cpu.regs[1] = uint8(r & 0xff00 >> 8)
		cpu.regs[0] = uint8(r & 0x00ff)
		if (r & 0x8000 >> 15) == 1 {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		return
	case INSN_STDY:
		// (Y) <- Rr
		y := uint16(cpu.regs[29]<<8) | uint16(cpu.regs[28])
		cpu.dmem[y+i.offset] = cpu.regs[i.source]
		return
	case INSN_STDZ:
		z := uint16(cpu.regs[31]<<8) | uint16(cpu.regs[30])
		cpu.dmem[z+i.offset] = cpu.regs[i.source]
		return
	case INSN_STX:
		x := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[x] = cpu.regs[i.source]
		return
	case INSN_STXM:
		x := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[x-1] = cpu.regs[i.source]
		return
	case INSN_STY:
		y := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[y] = cpu.regs[i.source]
		return
	case INSN_STZ:
		z := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[z] = cpu.regs[i.source]
		return
	case INSN_STZP:
		z := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[z] = cpu.regs[i.source]
		cpu.regs[26] += 1
		return
	case INSN_STZM:
		z := uint16(cpu.regs[27]<<8) | uint16(cpu.regs[26])
		cpu.dmem[z-1] = cpu.regs[i.source]
		return
	default:
		fmt.Println("I dunno.")
		return
	}
}
