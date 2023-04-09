package abstractions

import (
	"fmt"
	"gtihub.com/televi-go/televi/models/pages"
	"strings"
)

type TextHtmlBuilder struct {
	buildTasks []*FormatOptions
}

func (builder *TextHtmlBuilder) TextF(value string, args ...any) pages.IFormatter {
	formatter := &FormatOptions{source: fmt.Sprintf(value, args...)}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

type FormatOptions struct {
	source       string
	envelopOpen  string
	envelopClose string
}

func (builder *TextHtmlBuilder) length() int {
	acc := 0
	for _, task := range builder.buildTasks {
		acc += task.length()
	}
	return acc
}

func (formatOptions *FormatOptions) length() int {
	return len(formatOptions.source) + len(formatOptions.envelopClose) + len(formatOptions.envelopClose)
}

func (formatOptions *FormatOptions) writeTo(builder *strings.Builder) {
	builder.WriteString(formatOptions.envelopOpen)
	builder.WriteString(formatOptions.source)
	builder.WriteString(formatOptions.envelopClose)
}

func (formatOptions *FormatOptions) envelopIn(open, close string) {
	formatOptions.envelopOpen = open + formatOptions.envelopOpen
	formatOptions.envelopClose += close
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
	formatter := &FormatOptions{source: value}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

func (builder *TextHtmlBuilder) TextLine(value string) pages.IFormatter {
	formatter := &FormatOptions{source: value + "\n"}
	builder.buildTasks = append(builder.buildTasks, formatter)
	return formatter
}

func (builder *TextHtmlBuilder) ToString() string {
	tBuilder := strings.Builder{}
	tBuilder.Grow(builder.length())
	for _, task := range builder.buildTasks {
		task.writeTo(&tBuilder)
	}
	return tBuilder.String()
}
