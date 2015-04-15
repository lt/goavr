package main


import 	"encoding/hex"

// 4096 max size of program memory.
// Both data memory and program memory inherit
// from here.

type Memory [2048]byte

// Use Fetch() for grabbing 2 bytes from program memory.
// Increments the program counter.

func (mem *Memory) Fetch() []byte {
	ret := mem[cpu.pc:(cpu.pc+2)]
	cpu.pc += 2
	return ret
}

// Use Read() for reading a single byte from data memory
func (mem *Memory) Read(loc int) byte {
	return mem[loc]
}

// Loads the executable stuff into program memory.

func (mem *Memory) LoadProgram(data []byte) {
	for i, b := range(data) {
		mem[i] = b
	}
}

func (mem *Memory) Dump() string {
	return hex.Dump(mem[0:])
}


// Here but unused.
func (mem *Memory) Store(i int16, b byte) {
	mem[i] = b
}
