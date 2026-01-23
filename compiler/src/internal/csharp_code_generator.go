package lib

type CSharpCodeGenerator struct {
	BaseCodeGenerator
}

func NewCSharpCodeGenerator() *CSharpCodeGenerator {
	newObj := new(CSharpCodeGenerator)

	return newObj
}

func (this *CSharpCodeGenerator) Close() {
	this.close()
}

func (this *CSharpCodeGenerator) Generate(
	descriptor *TableDescriptor,
	outputDir string, newLineType NewLineType) bool {

	this.init(descriptor, newLineType)

	return true
}
