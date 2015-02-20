package main

import (
	"debug/elf"
	//"fmt"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getExecutableStuff(file *elf.File) {
	var x int
	for i, s := range file.Sections {
		if s.SectionHeader.Name == ".text" {
			x = i
		}
	}
	ret, err := file.Sections[x].Data()
	check(err)
	data = append(data, ret...)
}

func pop(a []byte, n int) []byte {
	ret := make([]byte, n)
	copy(ret, a)
	a = append(a[:0], a[n:]...)
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

func bigEndianConcat(b []byte) uint16 {
	ret := (uint16(b[0]) << 8) | uint16(b[1])
	return ret
}

func littleEndianConcat(b []byte) uint16 {
	ret := (uint16(b[1]) << 8) | uint16(b[0])
	return ret
}
