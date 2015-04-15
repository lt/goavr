package main

import (
	"debug/elf"
	"encoding/hex"
	"fmt"
)

var cSize int = 2
var pc = 0
var data []byte
var programEnd int16

func b2u16big(in []byte) uint16 { return (uint16(in[0]) << 8) | uint16(in[1]) }

func b2i16big(in []byte) int16 { return (int16(in[0]) << 8) | int16(in[1]) }

func b2u16little(in []byte) uint16 { return (uint16(in[1]) << 8) | uint16(in[0]) }

func b2i16little(in []byte) int16 { return (int16(in[1]) << 8) | int16(in[0]) }

func b2u32little(in []byte) uint32 { return (uint32(in[1]) << 8) | uint32(in[0]) }

func b2i32little(in []byte) int32 { return (int32(in[1]) << 8) | int32(in[0]) }

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getStuff(file *elf.File) {
	var x int
	// get executable stuff
	for i, s := range file.Sections {
		if s.SectionHeader.Name == ".text" {
			x = i
		}
	}
	ret, err := file.Sections[x].Data()
	check(err)
	data = append(data, ret...)
	// get the location of the last instruction.
	programEnd = int16(len(data) - 2)
	// get data stuff
	for i, s := range file.Sections {
		if s.SectionHeader.Name == ".data" {
			x = i
		}
	}
	ret, err = file.Sections[x].Data()
	check(err)
	data = append(data, ret...)
}

func dissectExecutable(file *elf.File) {
	for i, s := range file.Sections {
		dd, _ := file.Sections[i].Data()
		fmt.Printf("Section %d (%v)\n", i, s.SectionHeader.Name)
		fmt.Println(hex.Dump(dd))
	}
}

func pop(n int) []byte {
	ret := make([]byte, n)
	copy(ret, data)
	data = append(data[:0], data[n:]...)
	pc += n
	return ret
}

func chunkle(blob []byte, csize int) [][]byte {
	var fin = make([][]byte, 0)
	x := 0

	for i := 0; i < (len(blob) - csize); i += csize {
		fin = append(fin, []byte(blob[i:(x+csize)]))
		x += csize
	}
	fin = append(fin, []byte(blob[x:]))
	return fin
}

func printMnemonic(label int) {
	ret := fmt.Sprintf("I am %d\n", label)
	for _, op := range OpCodeLookUpTable {
		if op.label == label {
			ret = op.mnemonic
		}
	}
	fmt.Printf("%v (%d)\n", ret, label)
}

func printRegs(b [32]uint8) {
	var ret []string
	for i, v := range b {
		ret = append(ret, fmt.Sprintf("r%d[%d]", i, v))
	}
	fmt.Println("Registers:")
	fmt.Println(ret[0:17])
	fmt.Println(ret[17:32])

}
