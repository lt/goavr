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
var cSize int = 2
var instrCount = 0
var data []byte

func init() {
	flag.StringVar(&fileName, "f", fileName, "File path, yo")
}

func main() {

	flag.Parse()

	if fileName == "" {
		fileName = "/Users/erin/codebase/fouravr/Demo/firmware/main.elf"
	}

	file, _ := elf.Open(fileName)
	getExecutableStuff(file)

	for len(data) > 0 {
		fmt.Printf(".text:%.4x:\t", (instrCount))
		//fmt.Printf(hex.Dump(c))
		dissAssemble(pop(2))
	}
}
