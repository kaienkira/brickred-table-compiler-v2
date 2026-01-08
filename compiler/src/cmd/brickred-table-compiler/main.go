package main

import (
	"os"

	. "github.com/kaienkira/brickred-table-compiler-v2/compiler/internal"
)

func main() {

	// create parser
	parser := NewTableParser()
	defer parser.Close()

	os.Exit(0)
}
