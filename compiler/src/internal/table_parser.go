package lib

import (
	"fmt"
	"os"
	"strings"

	"github.com/antchfx/xmlquery"
)

type TableParser struct {
	Descriptor *TableDescriptor
}

func NewTableParser() *TableParser {
	newObj := new(TableParser)

	return newObj
}

func (this *TableParser) Close() {
	if this.Descriptor != nil {
		this.Descriptor.Close()
		this.Descriptor = nil
	}
}

func (this *TableParser) Parse(defineFilePath string) bool {
	// get file full path
	defineFileFullPath := UtilGetFullPath(defineFilePath)
	if defineFileFullPath == "" {
		fmt.Fprintf(os.Stderr,
			"error: can not find define file `%s`\n",
			defineFilePath)
		return false
	}

	// load xml xmlDoc
	xmlDoc := this.loadDefineFile(defineFileFullPath)
	if xmlDoc == nil {
		return false
	}

	this.Descriptor = NewTableDescriptor(defineFileFullPath)

	// check root node name
	var rootNode *xmlquery.Node = nil
	for _, child := range xmlDoc.ChildNodes() {
		if child.Type == xmlquery.ElementNode {
			rootNode = child
			break
		}
	}
	if rootNode == nil ||
		rootNode.Type != xmlquery.ElementNode ||
		rootNode.Data != "define" {
		this.printNodeError(rootNode,
			"root node must be `define` node")
		return false
	}

	// parse readers
	{
		nodes := xmlquery.Find(rootNode, "/reader")
		for _, node := range nodes {
			if this.addReaderDef(node) == false {
				return false
			}
		}
	}

	return true
}

func (this *TableParser) isStrValidVarName(str string) bool {
	return g_isVarNameRegexp.MatchString(str)
}

func (this *TableParser) printLineError(
	fileName string, lineNumber int, format string, args ...any) {

	fmt.Fprintf(os.Stderr,
		"error:%s:%d: %s\n",
		fileName, lineNumber,
		fmt.Sprintf(format, args...))
}

func (this *TableParser) printNodeError(
	node *xmlquery.Node, format string, args ...any) {

	this.printLineError(
		this.Descriptor.FilePath, node.LineNumber, format, args...)
}

func (this *TableParser) getNodeAttr(
	node *xmlquery.Node, attrName string) *xmlquery.Attr {

	for _, attr := range node.Attr {
		if attr.Name.Local == attrName {
			return &attr
		}
	}

	return nil
}

func (this *TableParser) loadDefineFile(filePath string) *xmlquery.Node {
	fileBin, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"error: can not read define file `%s`: %s\n",
			filePath, err.Error())
		return nil
	}
	fileText := string(fileBin)

	xmlDoc, err := xmlquery.ParseWithOptions(strings.NewReader(fileText),
		xmlquery.ParserOptions{
			WithLineNumbers: true,
		})
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"error: can not parse define file `%s`: %s\n",
			filePath, err.Error())
		return nil
	}

	return xmlDoc
}

func (this *TableParser) addReaderDef(node *xmlquery.Node) bool {
	// check name attr
	var name string
	{
		attr := this.getNodeAttr(node, "name")
		if attr == nil {
			this.printNodeError(node,
				"`reader` node must contain a `name` attribute")
			return false
		}
		name = attr.Value
	}
	if this.isStrValidVarName(name) == false {
		this.printNodeError(node,
			"`reader` node `name` attribute is invalid")
		return false
	}
	if _, ok := this.Descriptor.Readers[name]; ok {
		this.printNodeError(node,
			"`reader` node `name` attribute duplicated")
		return false
	}

	// check namespace attr
	var namespaceStr string
	{
		attr := this.getNodeAttr(node, "namespace")
		if attr == nil {
			this.printNodeError(node,
				"`reader` node must contain a `namespace` attribute")
			return false
		}
		namespaceStr = attr.Value
	}

	// check namespace parts
	namespaceParts := strings.Split(namespaceStr, ".")
	for _, part := range namespaceParts {
		if this.isStrValidVarName(part) == false {
			this.printNodeError(node,
				"`reader` node `namespace` attribute is invalid")
			return false
		}
	}

	def := NewReadDef(name, node.LineNumber)
	def.Namespace = namespaceStr
	def.NamespaceParts = namespaceParts

	this.Descriptor.Readers[def.Name] = def

	return true
}
