package main

import (
	"flag"
	"gonerics/parser"
	"log"
)


func main() {
	var input string
	var output string
	var execute bool

	flag.StringVar(&input, "in", "", "-in <input.go>")
	flag.StringVar(&output, "out", "", "-out <output.go>")
	flag.BoolVar(&execute, "run", false, "-run")

	flag.Parse()

	if input == "" || output == "" {
		log.Fatal("A valid input and output file must be provided")
	}

	if err := parser.Parse(input, output, execute); err != nil {
		log.Fatal(err.Error())
	}
}