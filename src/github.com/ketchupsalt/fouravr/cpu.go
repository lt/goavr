package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var programEnd int16

type CPU struct {
	pc   int16
	sp   StackPointer
	sr   byte
	imem Memory
	dmem Memory
}

// Set bits in status register
func (cpu *CPU) set_i() { cpu.dmem[cpu.sr] |= 128 }
func (cpu *CPU) set_t() { cpu.dmem[cpu.sr] |= 64 }
func (cpu *CPU) set_h() { cpu.dmem[cpu.sr] |= 32 }
func (cpu *CPU) set_s() { cpu.dmem[cpu.sr] |= 16 }
func (cpu *CPU) set_v() { cpu.dmem[cpu.sr] |= 8 }
func (cpu *CPU) set_n() { cpu.dmem[cpu.sr] |= 4 }
func (cpu *CPU) set_z() { cpu.dmem[cpu.sr] |= 2 }
func (cpu *CPU) set_c() { cpu.dmem[cpu.sr] |= 1 }

// Clear bits in status register
func (cpu *CPU) clear_i() { cpu.dmem[cpu.sr] &= 127 }
func (cpu *CPU) clear_t() { cpu.dmem[cpu.sr] &= 191 }
func (cpu *CPU) clear_h() { cpu.dmem[cpu.sr] &= 223 }
func (cpu *CPU) clear_s() { cpu.dmem[cpu.sr] &= 239 }
func (cpu *CPU) clear_v() { cpu.dmem[cpu.sr] &= 247 }
func (cpu *CPU) clear_n() { cpu.dmem[cpu.sr] &= 251 }
func (cpu *CPU) clear_z() { cpu.dmem[cpu.sr] &= 253 }
func (cpu *CPU) clear_c() { cpu.dmem[cpu.sr] &= 254 }

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
	fmt.Printf("pc: %.4x\tsr: %.8b\tsp: %.4x\t\n", cpu.pc, cpu.dmem[cpu.sr], cpu.sp.current())
	//defer handlePanic()
	cpu.imem.Fetch()
	cpu.Execute(dissAssemble(current))
	cpu.dmem.printRegs()
	cpu.dmem.printStack()
	fmt.Println("---------------------------------")
}

func (cpu *CPU) Interactive() {
	cpu.sr = 0x3f
	cpu.sp.high = 0x3e
	cpu.sp.low = 0x3d
	
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
				if cpu.pc == programEnd {
					break
				}
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
				if cpu.pc == programEnd {
					break
				}
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

	bitMasks := []byte{1, 2, 4, 8, 16, 32, 64, 128}

	switch i.label {
	case INSN_JMP:
		// we all know this doesn't work because
		// this version doesn't have a 22bit pc
		//cpu.pc = i.k32
		return
	case INSN_IJMP:
		// PC <- Z(15:0)
		z := b2i16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		cpu.pc = z << 1
		return
	case INSN_RJMP:
		// PC <- PC + k + 1
		cpu.pc = cpu.pc + i.k16
		return
	case INSN_ADD:
		// Rd <- Rd + Rr
		r := uint16(cpu.dmem[i.dest]) + uint16(cpu.dmem[i.source])
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
		cpu.dmem[i.dest] = byte(r)
		return
	case INSN_AND:
		// Rd <- Rd & Rr
		cpu.clear_v()
		r := cpu.dmem[i.dest] & cpu.dmem[i.source]
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x80 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		cpu.dmem[i.dest] = r
		return
	case INSN_ANDI:
		// Rd <- Rd & K
		cpu.dmem[i.dest] = cpu.dmem[i.dest] & i.kdata
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
		c := uint8(cpu.dmem[cpu.sr] & 0x01)
		r := cpu.dmem[i.dest] + cpu.dmem[i.source] + c
		cpu.dmem[i.dest] = r
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
		cpu.dmem[i.dest] = cpu.dmem[i.dest] ^ cpu.dmem[i.source]
		return
	case INSN_IN:
		// Rd <- I/O(A)
		cpu.dmem[i.dest] = cpu.dmem[i.ioaddr]
		return
	case INSN_OUT:
		// I/O(A) <- Rr
		cpu.dmem[i.ioaddr] = cpu.dmem[i.source]
		return
	case INSN_LDI:
		// Rd <- K
		cpu.dmem[i.dest] = i.kdata
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
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> i.registerBit
		if r == 0 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_SBIS:
		// If I/O(A,b) = 1 then PC <- PC + 2
		r := cpu.dmem[i.ioaddr] & bitMasks[i.registerBit] >> i.registerBit
		if r == 1 {
			// instructions are 1 word
			cpu.pc += 2
		}
		return
	case INSN_BLD:
		// Rd(b) <- T
		// Copies the T Flag in the SREG (Status Register) to bit b in register Rd.
		t := (cpu.dmem[cpu.sr] & bitMasks[7])
		cpu.dmem[i.dest] = byte(t)
		return
	case INSN_BST:
		// T <- Rd(b)
		// Stores bit b from Rd to the T Flag in SREG (Status Register).
		t := cpu.dmem[i.dest] & bitMasks[i.registerBit] >> i.registerBit
		cpu.dmem[cpu.sr] &= (t << 7)
		return
	case INSN_SBRC:
		// if Rr(b) = 0 then PC += 2
		r := (cpu.dmem[i.source] & bitMasks[i.registerBit]) >> i.registerBit
		if r == 0 {
			cpu.pc += 2
		}
		return
	case INSN_SBRS:
		// if Rr(b) = 1 then PC += 2
		r := (cpu.dmem[i.source] & bitMasks[i.registerBit]) >> i.registerBit
		if r == 1 {
			cpu.pc += 2
		}
		return
	case INSN_STS:
		// (k) <- Rr
		cpu.dmem[i.k16] = cpu.dmem[i.source]
		return
	case INSN_LDS:
		// Rd <- (k)
		cpu.dmem[i.dest] = cpu.dmem[i.k16]
		return
	case INSN_ADIW:
		// Rd+1:Rd <- Rd+1:Rd + K
		// low byte
		x := uint16(cpu.dmem[i.dest])
		// high byte
		y := uint16(cpu.dmem[i.dest+1])
		r := ((y << 8) | x) + uint16(i.kdata)
		cpu.dmem[i.dest] = uint8(r & 0x00ff)
		cpu.dmem[i.dest+1] = uint8(r >> 8)
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
		if i.kdata > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		// low byte
		x := uint16(cpu.dmem[i.dest])
		// high byte
		y := uint16(cpu.dmem[i.dest+1])
		r := ((y << 8) | x) - uint16(i.kdata)
		cpu.dmem[i.dest] = uint8(r & 0x00ff)
		cpu.dmem[i.dest+1] = uint8(r >> 8)
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
		c := cpu.dmem[cpu.sr] & 0x01
		if c == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRCS:
		// Branch if carry set
		c := cpu.dmem[cpu.sr] & 0x01
		if c == 1 {
			cpu.pc += i.k16
		}
		return
	case INSN_BREQ:
		//if Rd = Rr(Z=1) then PC <- PC + k + 1
		r := (cpu.dmem[cpu.sr] & bitMasks[7]) >> 7
		if r == 1 {
			cpu.pc += i.k16
		}
		return
	case INSN_BRGE:
		// if Rd >= Rr then PC += k
		n := cpu.dmem[cpu.sr] & bitMasks[2]
		v := cpu.dmem[cpu.sr] & bitMasks[3]
		if (n ^ v) == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRNE:
		// if (Z = 0) then PC <-  PC + k + 1
		z := (cpu.dmem[cpu.sr] & bitMasks[1]) >> 1
		if z == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_BRTC:
		// if T = 0 then PC <- PC + k + 1
		t := (cpu.dmem[cpu.sr] & bitMasks[6]) >> 6
		if t == 0 {
			cpu.pc += i.k16 //+1
		}
		return
	case INSN_COM:
		// Rd <- ^Rd
		r := ^cpu.dmem[i.dest]
		cpu.dmem[i.dest] = r
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
		r := uint16(cpu.dmem[i.dest]) - uint16(cpu.dmem[i.source])
		// XXX ToDo(ERIN) -- not so sure about the logic for C
		// "if the absolute value of the contents of Rr is larger than the absolute value of Rd"
		// means not zero, right?
		if r != 0 {
			cpu.clear_z()
			cpu.set_c()
		} else {
			cpu.set_z()
			cpu.clear_c()
		}
		if ((r & 0x0080) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0x00ff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		return
	case INSN_CPC:
		// Rd - Rr - C
		d := uint16(cpu.dmem[i.dest])
		s := uint16(cpu.dmem[i.source])
		c := uint16(cpu.dmem[cpu.sr] & 0x01)
		if (s + c) > d {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := d - s - c
		if r != 0 {
			cpu.clear_z()
		}
		if ((r & 0x0080) >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		if r > 0x00ff {
			cpu.set_v()
		} else {
			cpu.clear_v()
		}
		return
	case INSN_CPI:
		// Rd - K
		// I can't tell from the doc, but I think this check
		// has to happen before K is subtracted. We'll see.
		if i.kdata > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.dmem[i.dest] - i.kdata
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
		if cpu.dmem[i.dest] == cpu.dmem[i.source] {
			cpu.pc += 2
		}
		return
	case INSN_DEC:
		// Rd <- Rd - 1
		r := cpu.dmem[i.dest] - 1
		cpu.dmem[i.dest] = r
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
		y := b2u16little([]byte{cpu.dmem[28], cpu.dmem[29]}) + i.offset
		fmt.Printf("%.4x\t%.4x\n",y, cpu.dmem[y])
		cpu.dmem[i.dest] = cpu.dmem[y]
		return
	case INSN_LDDZ:
		// Rd <- (Z + q)
		z := b2u16little([]byte{cpu.dmem[30], cpu.dmem[31]}) + i.offset
		fmt.Printf("%.4x\t%.4x\n",z, cpu.dmem[z])
		cpu.dmem[i.dest] = cpu.dmem[z]
		return
	case INSN_LDX:
		// Rd <- (X)
		x := b2u16little([]byte{cpu.dmem[26], cpu.dmem[27]})
		cpu.dmem[i.dest] = cpu.dmem[x]
		return
	case INSN_LDXP:
		// Rd <- (X), X <- X + 1
		x := b2u16little([]byte{cpu.dmem[26], cpu.dmem[27]})
		cpu.dmem[i.dest] = cpu.dmem[x]
		cpu.dmem[26] += 1
		return
	case INSN_LDY:
		// Rd <- (Y)
		y := b2u16little([]byte{cpu.dmem[28], cpu.dmem[29]})
		cpu.dmem[i.dest] = cpu.dmem[y]
		return
	case INSN_LDYP:
		// Rd <- (Y), Y <- Y + 1
		y := b2u16little([]byte{cpu.dmem[28], cpu.dmem[29]})
		cpu.dmem[i.dest] = cpu.dmem[y]
		// XXX TODO(ERIN) this could overflow into the high
		// byte someday.
		cpu.dmem[28] += 1
		return
	case INSN_LDYM:
		// Rd <- (Y), Y <- Y - 1
		y := b2u16little([]byte{cpu.dmem[28], cpu.dmem[29]})
		// pre-decrement
		cpu.dmem[28] -= 1
		cpu.dmem[i.dest] = cpu.dmem[y]
		return
	case INSN_LDZ:
		// Rd <- (Z) (dmem)
		z := b2u16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		cpu.dmem[i.dest] = cpu.dmem[z]
		return
	case INSN_LDZP:
		// Rd <- (Z) (dmem), Z <- Z - 1
		z := b2u16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		cpu.dmem[i.dest] = cpu.dmem[z]
		// post-decrement
		// XXX TODO(ERIN) this could overflow into the high
		// byte some
		cpu.dmem[30] += 1
		return
	case INSN_LDZM:
		// Rd <- (Z) (dmem), Z <- Z - 1
		z := b2u16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		// pre-decrement
		cpu.dmem[30] -= 1
		cpu.dmem[i.dest] = cpu.dmem[z]
		return
	case INSN_LPMZ:
		// Rd <- (Z) (imem)
		z := b2i16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		cpu.dmem[i.dest] = cpu.imem[z]
		return
	case INSN_LPMZP:
		// Rd <- (Z), Z <- Z + 1 (imem)
		//z := int16(cpu.dmem[31] << 8) | int16(cpu.dmem[30])
		z := b2i16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		fmt.Printf("char found at %.4x:\t%.4x\n", z, cpu.imem[z])
		cpu.dmem[i.dest] = cpu.imem[z]
		// post-increment
		cpu.dmem[30] += 1
		return
	case INSN_LPM:
		// R0 <- (Z)
		z := b2i16little([]byte{cpu.dmem[30], cpu.dmem[31]})
		cpu.dmem[0] = cpu.imem[z]
		return
	case INSN_LSR:
		// logical shift right Rd
		r := cpu.dmem[i.dest] >> 1
		cpu.dmem[i.dest] = r
		return
	case INSN_MOV:
		// Rd <- Rr
		cpu.dmem[i.dest] = cpu.dmem[i.source]
		return
	case INSN_MOVW:
		// Rd+1:Rd <- Rr+1:Rr
		cpu.dmem[i.dest+1] = cpu.dmem[i.source+1]
		cpu.dmem[i.dest] = cpu.dmem[i.source]
		return
	case INSN_NEG:
		// Replaces the contents of register Rd with its two's complement
		r := ^cpu.dmem[i.dest] + 1
		cpu.dmem[i.dest] = r
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
		r := cpu.dmem[i.dest] | cpu.dmem[i.source]
		cpu.dmem[i.dest] = r
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
		r := cpu.dmem[i.dest] | uint8(i.kdata)
		cpu.dmem[i.dest] = r
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
		cpu.dmem[i.dest] = cpu.dmem[cpu.sp.current() + 1]
		cpu.sp.inc(1)
		return
	case INSN_PUSH:
		// STACK <- Rr
		cpu.dmem[cpu.sp.current()] = cpu.dmem[i.source]
		cpu.sp.dec(1)
		return
	case INSN_RCALL:
		// PC <- PC + k + 1, STACK <- PC + 1, SP - 2
		// push the current PC onto the stack because
		// it is automaticaly incremented elsewhere.
		// low byte
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc & 0x00ff)
		cpu.sp.dec(1)
		// high byte
		cpu.dmem[cpu.sp.current()] = byte(cpu.pc >> 8)
		cpu.sp.dec(1)
		// says +1, but that generates the wrong value
		// because the PC is incremented automaticaly anyway
		cpu.pc = cpu.pc + i.k16 //+ 1
		return
	case INSN_RET:
		// PC <- Stack
		// r29
		h := int16(cpu.dmem[cpu.sp.current() + 1])
		cpu.sp.inc(1)
		// r 28
		l := int16(cpu.dmem[cpu.sp.current() + 1])
		cpu.sp.inc(1)
		cpu.pc = ((h << 8) | l)
		fmt.Printf("%.4x\n", cpu.pc)
		return
	case INSN_RETI:
		// PC <- Stack, enable interrupts
		low := cpu.dmem[cpu.sp.current() - 1]
		high := cpu.dmem[cpu.sp.current()]
		cpu.pc = b2i16little([]byte{high, low})
		cpu.sp.inc(2)
		cpu.set_i()
		return
	case INSN_ROR:
		// Shifts all bits in Rd one place to the right.
		// The C Flag is shifted into bit 7 of Rd.
		// Bit 0 is shifted into the C Flag.
		c := (cpu.dmem[i.dest] & 0x80) >> 7
		r := (cpu.dmem[i.dest] >> 1) | uint8((cpu.dmem[cpu.sr]&0x01)<<7)
		cpu.dmem[i.dest] = r
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
		c := uint8(cpu.dmem[cpu.sr] & 0x01)
		if (cpu.dmem[i.source] + c) > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.dmem[i.dest] - cpu.dmem[i.source] - c
		cpu.dmem[i.dest] = r
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
		if i.kdata > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.dmem[i.dest] - i.kdata
		cpu.dmem[i.dest] = r
		if r == 0 {
			cpu.set_z()
		} else {
			cpu.clear_z()
		}
		if (r & 0x80 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_SBCI:
		// Rd <- Rd - K - C
		c := uint8(cpu.dmem[cpu.sr] & 0x01)
		if (i.kdata + c) > cpu.dmem[i.dest] {
			cpu.set_c()
		} else {
			cpu.clear_c()
		}
		r := cpu.dmem[i.dest] - i.kdata - c
		cpu.dmem[i.dest] = r
		if r == 0 {
			cpu.clear_z()
		}
		if (r & 0x80 >> 7) == 1 {
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
		r := cpu.dmem[i.dest] - cpu.dmem[i.source]
		cpu.dmem[i.dest] = r
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
		if (r & 0x80 >> 7) == 1 {
			cpu.set_n()
		} else {
			cpu.clear_n()
		}
		return
	case INSN_MUL:
		// R1h:R0l <- Rx x Rr
		r := uint16(cpu.dmem[i.dest]) * uint16(cpu.dmem[i.source])
		fmt.Printf("%.4x\n", r)
		cpu.dmem[1] = uint8(r & 0xff00 >> 8)
		cpu.dmem[0] = uint8(r & 0x00ff)
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
	case INSN_STX:
		x := uint16(cpu.dmem[27])<<8 | uint16(cpu.dmem[26])
		cpu.dmem[x] = cpu.dmem[i.source]
		return
	case INSN_STXP:
		// (X) <- Rr, X <- X + 1
		// 26 = low byte, 27 = high byte
		x :=  uint16(cpu.dmem[27])<<8 | uint16(cpu.dmem[26])
		cpu.dmem[x] = cpu.dmem[i.source]
		// post-increment
		cpu.dmem[26] += 1
		return
	case INSN_STXM:
		// pre-decrement
		cpu.dmem[26] -= 1
		x := uint16(cpu.dmem[27])<<8 | uint16(cpu.dmem[26])
		cpu.dmem[x] = cpu.dmem[i.source]
		return
	case INSN_STY:
		y := uint16(cpu.dmem[29])<<8 | uint16(cpu.dmem[28])
		cpu.dmem[y] = cpu.dmem[i.source]
		return
	case INSN_STYP:
		y := uint16(cpu.dmem[28])<<8 | uint16(cpu.dmem[28])
		cpu.dmem[y] = cpu.dmem[i.source]
		// post-increment
		cpu.dmem[28] += 1
		return
	case INSN_STYM:
		// pre-decrement
		cpu.dmem[28] -= 1
		y := uint16(cpu.dmem[29])<<8 | uint16(cpu.dmem[28])
		cpu.dmem[y] = cpu.dmem[i.source]
		return
	case INSN_STDY:
		// (Y) <- Rr
		y := uint16(cpu.dmem[29]) << 8 | uint16(cpu.dmem[28])
		fmt.Printf("%.4x\n", (y+i.offset))
		cpu.dmem[y+i.offset] = cpu.dmem[i.source]
		return
	case INSN_STZ:
		z := uint16(cpu.dmem[31])<<8 | uint16(cpu.dmem[30])
		cpu.dmem[z] = cpu.dmem[i.source]
		return
	case INSN_STZP:
		z := uint16(cpu.dmem[31])<<8 | uint16(cpu.dmem[30])
		cpu.dmem[z] = cpu.dmem[i.source]
		// post-increment
		cpu.dmem[26] += 1
		return
	case INSN_STZM:
		// pre-decrement
		cpu.dmem[30] -= 1
		z := uint16(cpu.dmem[31])<<8 | uint16(cpu.dmem[30])
		cpu.dmem[z] = cpu.dmem[i.source]
		return
	case INSN_STDZ:
		z := uint16(cpu.dmem[31]) <<8 | uint16(cpu.dmem[30])
		cpu.dmem[z+i.offset] = cpu.dmem[i.source]
		return
	default:
		fmt.Println("I dunno.")
		return
	}
}
