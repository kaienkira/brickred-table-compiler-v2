package main

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/kaienkira/brickred-table-compiler-v2/compiler/internal"
	flag "github.com/spf13/pflag"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, ""+
		"brickred table cutter\n"+
		"usage: %s "+
		"-f <define_file> "+
		"-r <reader> "+
		"-i <input_dir> "+
		"-o <output_dir>",
		filepath.Base(os.Args[0]))
}

func main() {
	// parse command line options
	var optDefineFilePath string
	var optReader string
	var optInputDir string
	var optOutputDir string

	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	flagSet.BoolP("help", "h", false, "")
	flagSet.StringVarP(&optDefineFilePath, "-define_file_path", "f", "", "")
	flagSet.StringVarP(&optReader, "-reader", "r", "", "")
	flagSet.StringVarP(&optInputDir, "-input_dir", "i", "", "")
	flagSet.StringVarP(&optOutputDir, "-output_dir", "o", "", "")
	flagSet.Parse(os.Args[1:])

	// check command line options
	// -- required options
	if optDefineFilePath == "" ||
		optReader == "" ||
		optInputDir == "" ||
		optOutputDir == "" {
		printUsage()
		os.Exit(1)
	}

	// -- check option define_file_path
	if UtilCheckFileExists(optDefineFilePath) == false {
		fmt.Fprintf(os.Stderr,
			"error: can not find define file `%s`\n",
			optDefineFilePath)
		os.Exit(1)
	}
}
