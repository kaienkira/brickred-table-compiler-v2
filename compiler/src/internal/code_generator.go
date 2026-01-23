package lib

type NewLineType int

const (
	NewLineType_None NewLineType = iota
	NewLineType_Unix
	NewLineType_Dos
)

type CodeGenerator interface {
	Close()
	Generate(descriptor *TableDescriptor,
		outputDir string, newLineType NewLineType) bool
}
