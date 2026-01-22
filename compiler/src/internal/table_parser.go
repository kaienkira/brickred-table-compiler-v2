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

	// parse global structs
	{
		nodes := xmlquery.Find(rootNode, "/struct")
		for _, node := range nodes {
			if this.addStructDef(nil, node) == false {
				return false
			}
		}
	}

	// parse tables
	{
		nodes := xmlquery.Find(rootNode, "/table")
		for _, node := range nodes {
			if this.addTableDef(node) == false {
				return false
			}
		}
	}

	return true
}

func (this *TableParser) FilterByReader(reader string) bool {
	if this.Descriptor == nil {
		return false
	}

	if _, ok := this.Descriptor.Readers[reader]; ok == false {
		fmt.Fprintf(os.Stderr,
			"error: reader `%s` is not defined\n",
			reader)
		return false
	}

	// remove unread tables
	filteredTables := make([]*TableDef, 0)
	for _, tableDef := range this.Descriptor.Tables {
		used := false
		if len(tableDef.Readers) == 0 {
			used = true
		} else {
			if _, ok := tableDef.Readers[reader]; ok {
				used = true
			}
		}
		if used {
			filteredTables = append(filteredTables, tableDef)
		} else {
			delete(this.Descriptor.TableNameIndex, tableDef.Name)
			tableDef.Close()
		}
	}
	this.Descriptor.Tables = filteredTables

	// remove unread columns
	for _, tableDef := range this.Descriptor.Tables {
		filteredColumns := make([]*TableColumnDef, 0)
		for _, columnDef := range tableDef.Columns {
			used := false
			if columnDef == tableDef.TableKey {
				used = true
			} else if len(columnDef.Readers) == 0 {
				used = true
			} else {
				if _, ok := columnDef.Readers[reader]; ok {
					used = true
				}
			}
			if used {
				filteredColumns = append(filteredColumns, columnDef)
			} else {
				delete(tableDef.ColumnNameIndex, columnDef.Name)
				columnDef.Close()
			}
		}
		tableDef.Columns = filteredColumns
		this.calculateTableKeyColumnIndex(tableDef)
	}

	// collect used structs
	usedStructs := make(map[*StructDef]bool)
	for _, tableDef := range this.Descriptor.Tables {
		for _, columnDef := range tableDef.Columns {
			if columnDef.RefStructDef != nil {
				usedStructs[columnDef.RefStructDef] = true
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

func (this *TableParser) addStructDef(
	tableDef *TableDef, node *xmlquery.Node) bool {

	// check name attr
	var name string
	{
		attr := this.getNodeAttr(node, "name")
		if attr == nil {
			this.printNodeError(node,
				"`struct` node must contain a `name` attribute")
			return false
		}
		name = attr.Value
	}
	if this.isStrValidVarName(name) == false {
		this.printNodeError(node,
			"`struct` node `name` attribute is invalid")
		return false
	}
	if tableDef == nil {
		ok := false
		if _, ok = this.Descriptor.GlobalStructNameIndex[name]; ok == false {
			_, ok = this.Descriptor.TableNameIndex[name]
		}
		if ok {
			this.printNodeError(node,
				"`struct` node `name` attribute duplicated")
			return false
		}
	} else {
		if _, ok := tableDef.LocalStructNameIndex[name]; ok {
			this.printNodeError(node,
				"`struct` node `name` attribute duplicated")
			return false
		}
		if name == "Row" ||
			name == "Rows" ||
			name == "RowSet" ||
			name == "RowSets" {
			this.printNodeError(node, ""+
				"local struct can not be named as "+
				"`Row`, `Rows`, `RowSet` or `RowSets`")
			return false
		}
	}

	def := NewStructDef(tableDef, name, node.LineNumber)

	// parse fields
	for _, childNode := range node.ChildNodes() {
		if childNode.Type != xmlquery.ElementNode {
			continue
		}
		if childNode.Data != "field" {
			this.printNodeError(childNode,
				"expect a `field` node")
			return false
		}

		if this.addStructFieldDef(def, childNode) == false {
			return false
		}
	}

	if tableDef == nil {
		this.Descriptor.GlobalStructs =
			append(this.Descriptor.GlobalStructs, def)
		this.Descriptor.GlobalStructNameIndex[def.Name] = def
	} else {
		tableDef.LocalStructs = append(tableDef.LocalStructs, def)
		tableDef.LocalStructNameIndex[def.Name] = def
	}

	return true
}

func (this *TableParser) addStructFieldDef(
	structDef *StructDef, node *xmlquery.Node) bool {

	// check name attr
	var name string
	{
		attr := this.getNodeAttr(node, "name")
		if attr == nil {
			this.printNodeError(node,
				"`field` node must contain a `name` attribute")
			return false
		}
		name = attr.Value
	}
	if this.isStrValidVarName(name) == false {
		this.printNodeError(node,
			"`field` node `name` attribute is invalid")
		return false
	}
	if _, ok := structDef.FieldNameIndex[name]; ok {
		this.printNodeError(node,
			"`field` node `name` attribute duplicated")
		return false
	}

	// check type attr
	var typ string
	{
		attr := this.getNodeAttr(node, "type")
		if attr == nil {
			this.printNodeError(node,
				"`field` node must contain a `type` attribute")
			return false
		}
		typ = attr.Value
	}

	def := NewStructFieldDef(structDef, name, node.LineNumber)

	if typ == "int" {
		def.Type = StructFieldType_Int
	} else if typ == "string" {
		def.Type = StructFieldType_String
	} else {
		this.printNodeError(node,
			"type `%s` is invalid", typ)
		return false
	}

	structDef.Fields = append(structDef.Fields, def)
	structDef.FieldNameIndex[def.Name] = def

	return true
}

func (this *TableParser) addTableDef(node *xmlquery.Node) bool {
	// check name attr
	var name string
	{
		attr := this.getNodeAttr(node, "name")
		if attr == nil {
			this.printNodeError(node,
				"`table` node must contain a `name` attribute")
			return false
		}
		name = attr.Value
	}
	if this.isStrValidVarName(name) == false {
		this.printNodeError(node,
			"`table` node `name` attribute is invalid")
		return false
	}

	ok := false
	if _, ok = this.Descriptor.TableNameIndex[name]; ok == false {
		_, ok = this.Descriptor.GlobalStructNameIndex[name]
	}
	if ok {
		this.printNodeError(node,
			"`table` node `name` attribute duplicated")
		return false
	}

	def := NewTableDef(name, node.LineNumber)

	for _, childNode := range node.ChildNodes() {
		if childNode.Type != xmlquery.ElementNode {
			continue
		}

		if childNode.Data == "struct" {
			// parse local struct
			if this.addStructDef(def, childNode) == false {
				return false
			}
		} else if childNode.Data == "col" {
			// parse column
			if this.addTableColumnDef(def, childNode) == false {
				return false
			}
		} else {
			this.printNodeError(childNode,
				"expect a `struct` or `col` node")
		}
	}

	// check key/setkey attr
	{
		var key string
		attr := this.getNodeAttr(node, "key")
		if attr != nil {
			key = attr.Value
			def.TableKeyType = TableKeyType_SingleKey
		} else {
			attr = this.getNodeAttr(node, "setkey")
			if attr != nil {
				key = attr.Value
				def.TableKeyType = TableKeyType_SetKey
			} else {
				this.printNodeError(node,
					"`table` node must contain a `key` or `setkey` attribute")
				return false
			}
		}

		tableKey, ok := def.ColumnNameIndex[key]
		if ok == false {
			this.printNodeError(node,
				"table key `%s` is not defined", key)
			return false
		}
		if tableKey.Type != TableColumnType_Int &&
			tableKey.Type != TableColumnType_String {
			this.printNodeError(node,
				"table key can only be `int` or `string` type")
			return false
		}
		def.TableKey = tableKey
	}

	// check file attr
	{
		attr := this.getNodeAttr(node, "file")
		if attr == nil {
			this.printNodeError(node,
				"`table` node must contain a `file` attribute")
			return false
		}
		def.FileName = attr.Value
	}

	// check readby attr
	{
		attr := this.getNodeAttr(node, "readby")
		if attr != nil {
			for reader := range strings.SplitSeq(attr.Value, "|") {
				readerDef, ok := this.Descriptor.Readers[reader]
				if ok == false {
					this.printNodeError(node,
						"reader `%s` is not defined", reader)
					return false
				}
				def.Readers[reader] = readerDef
			}
		}
	}

	this.calculateTableKeyColumnIndex(def)
	this.Descriptor.Tables = append(this.Descriptor.Tables, def)
	this.Descriptor.TableNameIndex[def.Name] = def

	return true
}

func (this *TableParser) addTableColumnDef(
	tableDef *TableDef, node *xmlquery.Node) bool {

	// check name attr
	var name string
	{
		attr := this.getNodeAttr(node, "name")
		if attr == nil {
			this.printNodeError(node,
				"`col` node must contain a `name` attribute")
			return false
		}
		name = attr.Value
	}
	if this.isStrValidVarName(name) == false {
		this.printNodeError(node,
			"`col` node `name` attribute is invalid")
		return false
	}
	if _, ok := tableDef.ColumnNameIndex[name]; ok {
		this.printNodeError(node,
			"`col` node `name` attribute duplicated")
		return false
	}

	// check type attr
	var typ string
	{
		attr := this.getNodeAttr(node, "type")
		if attr == nil {
			this.printNodeError(node,
				"`col` node must contain a `type` attribute")
			return false
		}
		typ = attr.Value
	}

	def := NewTableColumnDef(tableDef, name, node.LineNumber)

	// get type info
	columnTypeStr := typ
	{
		m := g_fetchListTypeRegexp.FindStringSubmatch(columnTypeStr)
		if m != nil {
			columnTypeStr = m[1]
			def.Type = TableColumnType_List
		}
	}

	columnType := TableColumnType_None
	if columnTypeStr == "int" {
		columnType = TableColumnType_Int
	} else if columnTypeStr == "string" {
		columnType = TableColumnType_String
	} else {
		if refStructDef, ok :=
			tableDef.LocalStructNameIndex[columnTypeStr]; ok {
			// check is local struct
			columnType = TableColumnType_Struct
			def.RefStructDef = refStructDef
		} else if refStructDef, ok :=
			this.Descriptor.GlobalStructNameIndex[columnTypeStr]; ok {
			// check is global struct
			columnType = TableColumnType_Struct
			def.RefStructDef = refStructDef
		} else {
			this.printNodeError(node,
				"type `%s` is invalid", typ)
			return false
		}
	}

	if def.Type == TableColumnType_List {
		def.ListType = columnType
	} else {
		def.Type = columnType
	}

	// check readby attr
	{
		attr := this.getNodeAttr(node, "readby")
		if attr != nil {
			for reader := range strings.SplitSeq(attr.Value, "|") {
				readerDef, ok := this.Descriptor.Readers[reader]
				if ok == false {
					this.printNodeError(node,
						"reader `%s` is not defined", reader)
					return false
				}
				def.Readers[reader] = readerDef
			}
		}
	}

	tableDef.Columns = append(tableDef.Columns, def)
	tableDef.ColumnNameIndex[def.Name] = def

	return true
}

func (this *TableParser) calculateTableKeyColumnIndex(tableDef *TableDef) {
	for i, columnDef := range tableDef.Columns {
		if columnDef == tableDef.TableKey {
			tableDef.TableKeyColumnIndex = i
			return
		}
	}
}
