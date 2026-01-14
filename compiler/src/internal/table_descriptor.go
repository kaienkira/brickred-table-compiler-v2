package lib

type TableDescriptor struct {
	FilePath string

	// reader define
	// ReaderDef.Name -> ReaderDef
	Readers map[string]*ReaderDef

	// global struct define
	// in file define order
	GlobalStructs []*StructDef
	// StructDef.Name -> StructDef
	GlobalStructNameIndex map[string]*StructDef
}

func NewTableDescriptor() *TableDescriptor {
	newObj := new(TableDescriptor)
	newObj.Readers = make(map[string]*ReaderDef)

	return newObj
}

func (this *TableDescriptor) Close() {
	if this.GlobalStructNameIndex != nil {
		clear(this.GlobalStructNameIndex)
		this.GlobalStructNameIndex = nil
	}
	if this.GlobalStructs != nil {
		for _, def := range this.GlobalStructs {
			def.Close()
		}
		clear(this.GlobalStructs)
		this.GlobalStructs = nil
	}
	if this.Readers != nil {
		for _, def := range this.Readers {
			def.Close()
		}
		clear(this.Readers)
		this.Readers = nil
	}
}

// ----------------------------------------------------------------------------
type ReaderDef struct {
	// reader name
	Name string
	// define in line number
	LineNumber int

	Namespace      string
	NamespaceParts []string
}

func NewReadDef(name string, lineNumber int) *ReaderDef {
	newObj := new(ReaderDef)
	newObj.Name = name
	newObj.LineNumber = lineNumber

	return newObj
}

func (this *ReaderDef) Close() {
	this.NamespaceParts = nil
}

// ----------------------------------------------------------------------------
type StructFieldType int

const (
	StructFieldType_None StructFieldType = iota
	StructFieldType_Int
	StructFieldType_String
)

// ----------------------------------------------------------------------------
type StructFieldDef struct {
	// link to parent define
	ParentRef *StructDef
	// field name
	Name string
	// define in line number
	LineNumber int

	Type StructFieldType
}

func NewStructFieldDef(
	parentRef *StructDef, name string, lineNumber int) *StructFieldDef {

	newObj := new(StructFieldDef)
	newObj.ParentRef = parentRef
	newObj.Name = name
	newObj.LineNumber = lineNumber

	return newObj
}

func (this *StructFieldDef) Close() {
	this.ParentRef = nil
}

// ----------------------------------------------------------------------------
type StructDef struct {
	// link to parent define, null when struct is global
	ParentRef *TableDef
	// struct name
	Name string
	// define in line number
	LineNumber int

	// in file define order
	Fields []*StructFieldDef
	// FieldDef.Name -> FieldDef
	FieldNameIndex map[string]*StructFieldDef
}

func NewStructDef(
	parentRef *TableDef, name string, lineNumber int) *StructDef {

	newObj := new(StructDef)
	newObj.ParentRef = parentRef
	newObj.Name = name
	newObj.LineNumber = lineNumber
	newObj.Fields = make([]*StructFieldDef, 0)
	newObj.FieldNameIndex = make(map[string]*StructFieldDef)

	return newObj
}

func (this *StructDef) Close() {
	if this.FieldNameIndex != nil {
		clear(this.FieldNameIndex)
		this.FieldNameIndex = nil
	}
	if this.Fields != nil {
		for _, def := range this.Fields {
			def.Close()
		}
		clear(this.Fields)
		this.Fields = nil
	}
	this.ParentRef = nil
}

// ----------------------------------------------------------------------------
type TableColumnType int

const (
	TableColumnType_None TableColumnType = iota
	TableColumnType_Int
	TableColumnType_String
	TableColumnType_Struct
	TableColumnType_List
)

// ----------------------------------------------------------------------------
type TableKeyType int

const (
	TableKeyType_None TableKeyType = iota
	TableKeyType_SingleKey
	TableKeyType_SetKey
)

// ----------------------------------------------------------------------------
type TableColumnDef struct {
	// link to parent define
	ParentRef *TableDef
	// column name
	Name string
	// define in line number
	LineNumber int

	Type         TableColumnType
	ListType     TableColumnType
	RefStructDef *StructDef
	Readers      map[string]*ReaderDef
}

func NewTableColumnDef(
	parentRef *TableDef, name string, lineNumber int) *TableColumnDef {

	newObj := new(TableColumnDef)
	newObj.ParentRef = parentRef
	newObj.Name = name
	newObj.LineNumber = lineNumber
	newObj.Readers = make(map[string]*ReaderDef)

	return newObj
}

func (this *TableColumnDef) Close() {
	if this.Readers != nil {
		clear(this.Readers)
		this.Readers = nil
	}
	this.RefStructDef = nil
	this.ParentRef = nil
}

// ----------------------------------------------------------------------------
type TableDef struct {
	// table name
	Name string
	// define in line number
	LineNumber int

	// table key
	TableKey *TableColumnDef
	// table key type
	TableKeyType TableKeyType
	// table key column index
	TableKeyColumnIndex int
	// file name
	FileName string
	// read by
	Readers map[string]*ReaderDef
	// in file define order
	LocalStructs []*StructDef
	// StructDef.Name -> StructDef
	LocalStructNameIndex map[string]*StructDef
	// in file define order
	Columns []*TableColumnDef
	// ColumnDef.Name -> ColumnDef
	ColumnNameIndex map[string]*TableColumnDef
}

func NewTableDef(name string, lineNumber int) *TableDef {
	newObj := new(TableDef)
	newObj.Name = name
	newObj.LineNumber = lineNumber
	newObj.Readers = make(map[string]*ReaderDef)
	newObj.LocalStructs = make([]*StructDef, 0)
	newObj.LocalStructNameIndex = make(map[string]*StructDef)
	newObj.Columns = make([]*TableColumnDef, 0)
	newObj.ColumnNameIndex = make(map[string]*TableColumnDef)

	return newObj
}

func (this *TableDef) Close() {
	if this.ColumnNameIndex != nil {
		clear(this.ColumnNameIndex)
		this.ColumnNameIndex = nil
	}
	if this.Columns != nil {
		for _, def := range this.Columns {
			def.Close()
		}
		clear(this.Columns)
		this.Columns = nil
	}
	if this.LocalStructNameIndex != nil {
		clear(this.LocalStructNameIndex)
		this.LocalStructNameIndex = nil
	}
	if this.LocalStructs != nil {
		for _, def := range this.LocalStructs {
			def.Close()
		}
		clear(this.LocalStructs)
		this.LocalStructs = nil
	}
	if this.Readers != nil {
		clear(this.Readers)
		this.Readers = nil
	}
	this.TableKey = nil
}
