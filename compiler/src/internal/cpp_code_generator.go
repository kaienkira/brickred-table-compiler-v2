package lib

type CppCodeGenerator struct {
	BaseCodeGenerator
}

func NewCppCodeGenerator() *CppCodeGenerator {
	newObj := new(CppCodeGenerator)

	return newObj
}

func (this *CppCodeGenerator) Close() {
	this.close()
}

func (this *CppCodeGenerator) Generate(
	descriptor *TableDescriptor,
	outputDir string, newLineType NewLineType) bool {

	this.init(descriptor, newLineType)

	return true
}
