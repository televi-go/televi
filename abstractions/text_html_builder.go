package abstractions

import (
	"fmt"
	"github.com/televi-go/televi/models/pages"
	"strings"
)

type Buildable interface {
	WriteTo(builder *strings.Builder)
}

type wrapper struct {
	Buildable
}

func (w wrapper) writeTo(builder *strings.Builder) {
	w.WriteTo(builder)
}

type TextHtmlBuilder struct {
	buildTasks []*FormatOptions
}

func (builder *TextHtmlBuilder) AddBuildable(buildable Buildable) {
	builder.buildTasks = append(builder.buildTasks, &FormatOptions{
		source:       wrapper{buildable},
		envelopOpen:  "",
		envelopClose: "",
	})
}

func (builder *TextHtmlBuilder) TextF(value string, args ...any) pages.IFormatter {
	formatter := &FormatOptions{source: StrFormatWriteable(fmt.Sprintf(value, args...))}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

type FormatWritable interface {
	writeTo(builder *strings.Builder)
}

type StrFormatWriteable string

func (s StrFormatWriteable) writeTo(builder *strings.Builder) {
	builder.WriteString(string(s))
}

type FormatOptions struct {
	source       FormatWritable
	envelopOpen  string
	envelopClose string
}

func (formatOptions *FormatOptions) writeTo(builder *strings.Builder) {
	builder.WriteString(formatOptions.envelopOpen)
	formatOptions.source.writeTo(builder)
	builder.WriteString(formatOptions.envelopClose)
}

func (formatOptions *FormatOptions) envelopIn(open, close string) {
	nodeCopy := *formatOptions
	formatOptions.envelopOpen = open
	formatOptions.envelopClose = close
	formatOptions.source = &nodeCopy
}

func (formatOptions *FormatOptions) Bold() pages.IFormatter {
	formatOptions.envelopIn("<b>", "</b>")
	return formatOptions
}

func (formatOptions *FormatOptions) Mono() pages.IFormatter {
	formatOptions.envelopIn("<code>", "</code>")
	return formatOptions
}

func (formatOptions *FormatOptions) Spoiler() pages.IFormatter {
	formatOptions.envelopIn("<tg-spoiler>", "</tg-spoiler>")
	return formatOptions
}

func (builder *TextHtmlBuilder) Text(value string) pages.IFormatter {
	formatter := &FormatOptions{source: StrFormatWriteable(value)}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

func (builder *TextHtmlBuilder) TextLine(value string) pages.IFormatter {
	formatter := &FormatOptions{source: StrFormatWriteable(value + "\n")}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

func (builder *TextHtmlBuilder) ToString() string {
	tBuilder := strings.Builder{}
	for _, task := range builder.buildTasks {
		task.writeTo(&tBuilder)
	}
	return tBuilder.String()
}

type TextBuilder interface {
	Text(value string) IFormatter
	TextF(value string, args ...any) IFormatter
	TextLine(value string) IFormatter
}

type IFormatter interface {
	Bold() IFormatter
	Mono() IFormatter
	Spoiler() IFormatter
	//TODO: to be added
}
