package main

// Balls

import (
	"debug/elf"
	"flag"
	"fmt"
	//"encoding/hex"
	//"reflect"
)

var fileName string
var chunkSize int = 2
var data =  make([]byte, 0)

func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
}

func main() {

	flag.Parse()

	if fileName == "" {
		fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
	}

	file, _ := elf.Open(fileName)
	getExecutableStuff(file, data)

	for i :=0 ; i < len(data); i++ {
		fmt.Printf(".text:0x%.8x\t", (i << 1))
		//fmt.Printf(hex.Dump(c))
		dissAssemble(pop(data, 2))
	}
}
