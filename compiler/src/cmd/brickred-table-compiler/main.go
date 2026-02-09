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
		"brickred table compiler\n"+
		"usage: %s "+
		"-f <define_file> "+
		"-l <language> "+
		"-r <reader>"+
		"\n"+
		"    [-o <output_dir>]\n"+
		"    [-n <new_line_type>] (unix|dos) default is unix\n"+
		"language supported: cpp csharp\n",
		filepath.Base(os.Args[0]))
}

func run() int {
	// parse command line options
	var optHelp bool
	var optDefineFilePath string
	var optLanguage string
	var optReader string
	var optOutputDir string
	var optNewLineType string

	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	flagSet.BoolVarP(&optHelp, "help", "h", false, "")
	flagSet.StringVarP(&optDefineFilePath, "-define_file_path", "f", "", "")
	flagSet.StringVarP(&optLanguage, "-language", "l", "", "")
	flagSet.StringVarP(&optReader, "-reader", "r", "", "")
	flagSet.StringVarP(&optOutputDir, "-output_dir", "o", "", "")
	flagSet.StringVarP(&optNewLineType, "-new_line_type", "n", "", "")

	if flagSet.Parse(os.Args[1:]) != nil {
		printUsage()
		return 1
	}
	if optHelp {
		printUsage()
		return 0
	}

	// check command line options
	// -- required options
	if optDefineFilePath == "" ||
		optLanguage == "" {
		printUsage()
		return 1
	}
	// -- option default value
	if optOutputDir == "" {
		optOutputDir = "."
	}
	if optNewLineType == "" {
		optNewLineType = "unix"
	}

	// -- check option define_file_path
	if UtilCheckFileExists(optDefineFilePath) == false {
		fmt.Fprintf(os.Stderr,
			"error: can not find define file `%s`\n",
			optDefineFilePath)
		return 1
	}

	// -- check option language
	if optLanguage != "cpp" &&
		optLanguage != "csharp" {
		fmt.Fprintf(os.Stderr,
			"error: language `%s` is not supported\n",
			optLanguage)
		return 1
	}

	// -- check option output_dir
	if UtilCheckDirExists(optOutputDir) == false {
		fmt.Fprintf(os.Stderr,
			"error: can not find output directory `%s`\n",
			optOutputDir)
		return 1
	}

	// -- check option new_line_type
	if optNewLineType != "dos" &&
		optNewLineType != "unix" {
		fmt.Fprintf(os.Stderr,
			"error: new_line_type `%s` is invalid\n",
			optNewLineType)
		return 1
	}

	// create parser
	parser := NewTableParser()
	if parser.Parse(optDefineFilePath) == false {
		return 1
	}
	defer parser.Close()
	if optReader != "" {
		if parser.FilterByReader(optReader) == false {
			return 1
		}
	}

	// create generator
	var generator CodeGenerator = nil
	if optLanguage == "cpp" {
		generator = NewCppCodeGenerator()
	} else if optLanguage == "csharp" {
		generator = NewCSharpCodeGenerator()
	} else {
		return 1
	}
	defer generator.Close()

	// generate code
	newLineType := NewLineType_Unix
	if optNewLineType == "dos" {
		newLineType = NewLineType_Dos
	}
	if generator.Generate(parser.Descriptor,
		optReader, optOutputDir, newLineType) == false {
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
