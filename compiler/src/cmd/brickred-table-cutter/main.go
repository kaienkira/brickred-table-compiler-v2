package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		"-o <output_dir>\n",
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

	// -- check option input_dir
	if UtilCheckDirExists(optInputDir) == false {
		fmt.Fprintf(os.Stderr,
			"error: can not find input directory `%s`\n",
			optInputDir)
		os.Exit(1)
	}

	// -- check option output_dir
	if UtilCheckDirExists(optOutputDir) == false {
		fmt.Fprintf(os.Stderr,
			"error: can not find output directory `%s`\n",
			optOutputDir)
		os.Exit(1)
	}
	if UtilGetFullPath(optInputDir) ==
		UtilGetFullPath(optOutputDir) {
		fmt.Fprintf(os.Stderr,
			"error: output directory can not be same as input directory\n")
		os.Exit(1)
	}

	// create parser
	parser := NewTableParser()
	if parser.Parse(optDefineFilePath) == false {
		os.Exit(1)
	}
	defer parser.Close()

	if cutTables(parser.Descriptor,
		optReader, optInputDir, optOutputDir) == false {
		os.Exit(1)
	}

	os.Exit(0)
}

func cutTables(descriptor *TableDescriptor,
	reader string, inputDir string, outputDir string) bool {

	// check reader
	if _, ok := descriptor.Readers[reader]; ok == false {
		fmt.Fprintf(os.Stderr,
			"error: reader `%s` is not defined\n",
			reader)
		return false
	}

	for _, def := range descriptor.Tables {
		needCut := false
		if len(def.Readers) <= 0 {
			needCut = true
		} else if _, ok := def.Readers[reader]; ok {
			needCut = true
		}
		if needCut == false {
			continue
		}
		if cutTable(def, reader, inputDir, outputDir) == false {
			return false
		}
	}

	return true
}

func cutTable(tableDef *TableDef,
	reader string, inputDir string, outputDir string) bool {

	// calucate deleted columns
	deletedColumns := make(map[int]bool)
	for i, def := range tableDef.Columns {
		if def == tableDef.TableKey {
			continue
		}
		if len(def.Readers) <= 0 {
			continue
		}
		if _, ok := def.Readers[reader]; ok {
			continue
		}
		deletedColumns[i] = true
	}

	// read input file
	inputFilePath := filepath.Join(inputDir, tableDef.FileName)
	inputFileContent, ret := UtilReadAllTextShared(inputFilePath)
	if ret == false {
		return false
	}

	// split lines
	lines := strings.Split(inputFileContent, "\r\n")
	if lines[len(lines)-1] != "" {
		fmt.Fprintf(os.Stderr,
			"error: input file `%s` file line ending is required\n",
			tableDef.FileName)
		return false
	}

	lineCount := len(lines) - 1
	if lineCount < 2 {
		fmt.Fprintf(os.Stderr,
			"error: input file `%s` comment line and name line is required\n",
			tableDef.FileName)
		return false
	}

	// split columns
	lineCols := make([][]string, 0)
	for i := range lineCount {
		cols := strings.Split(lines[i], "\t")
		if len(cols) != len(tableDef.Columns) {
			fmt.Fprintf(os.Stderr, ""+
				"error: input file `%s` line %d "+
				"column count %d is invalid, should be %d\n",
				tableDef.FileName, i+1,
				len(cols), len(tableDef.Columns))
			return false
		}
		lineCols = append(lineCols, cols)
	}

	// check name line
	for i, def := range tableDef.Columns {
		if lineCols[1][i] != def.Name {
			fmt.Fprintf(os.Stderr,
				"error: input file `%s` column %d should be named as `%s`\n",
				tableDef.FileName, i+1, def.Name)
			return false
		}
	}

	// cut columns
	var sb strings.Builder
	outputCols := make([]string, 0)
	for i := range lineCount {
		inputCols := lineCols[i]
		outputCols := outputCols[:0]
		for j, col := range inputCols {
			if _, ok := deletedColumns[j]; ok {
				continue
			}
			outputCols = append(outputCols, col)
		}
		sb.WriteString(strings.Join(outputCols, "\t"))
		sb.WriteString("\r\n")
	}
	outputFileContent := sb.String()

	// write output file
	outputFilePath := filepath.Join(outputDir, tableDef.FileName)
	if UtilWriteAllText(outputFilePath, outputFileContent) == false {
		return false
	}

	return true
}
