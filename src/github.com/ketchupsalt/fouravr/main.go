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

	i := 0
	for len(data) > 0 {
		fmt.Printf(".text:0x%.8x\t", (i << 1))
		//fmt.Printf(hex.Dump(c))
		dissAssemble(pop(2))
		i++
	}
}
