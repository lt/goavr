package main

import (
	"debug/elf"
	"flag"
	"fmt"
)

var fileName string
var cSize int = 2
var pc = 0
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
		fmt.Printf(".text:%.4x:\t", (pc))
		dissAssemble(pop(2))
	}
}
