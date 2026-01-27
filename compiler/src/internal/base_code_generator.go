package lib

import (
	"fmt"
	"strings"
)

type BaseCodeGenerator struct {
	descriptor *TableDescriptor
	reader     string
	newLineStr string
}

func (this *BaseCodeGenerator) init(
	descriptor *TableDescriptor,
	reader string, newLineType NewLineType) {

	this.descriptor = descriptor
	this.reader = reader
	if newLineType == NewLineType_Dos {
		this.newLineStr = "\r\n"
	} else {
		this.newLineStr = "\n"
	}
}

func (this *BaseCodeGenerator) close() {
	this.descriptor = nil
}

func (this *BaseCodeGenerator) writeLine(
	sb *strings.Builder, line string) {

	sb.WriteString(line)
	sb.WriteString(this.newLineStr)
}

func (this *BaseCodeGenerator) writeLineFormat(
	sb *strings.Builder, format string, args ...any) {

	fmt.Fprintf(sb, format, args...)
	sb.WriteString(this.newLineStr)
}

func (this *BaseCodeGenerator) writeEmptyLine(
	sb *strings.Builder) {

	sb.WriteString(this.newLineStr)
}
