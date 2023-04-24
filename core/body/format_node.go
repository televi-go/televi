package body

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/models/pages"
	"strings"
)

type FormatNode struct {
	preamble string
	content  abstractions.Buildable
	epilogue string
}

type StrBuildable string

func (s StrBuildable) WriteTo(builder *strings.Builder) {
	builder.WriteString(string(s))
}

func NewFormatNode(source string) *FormatNode {
	return &FormatNode{
		preamble: "",
		content:  StrBuildable(source),
		epilogue: "",
	}
}

func (formatNode *FormatNode) envelopIn(open, close string) *FormatNode {
	formatNodeCopy := *formatNode
	formatNode.preamble = open
	formatNode.content = &formatNodeCopy
	formatNode.epilogue = close
	return formatNode
}

func (formatNode *FormatNode) Bold() pages.IFormatter {
	return formatNode.envelopIn("<b>", "</b>")
}

func (formatNode *FormatNode) Mono() pages.IFormatter {
	return formatNode.envelopIn("<code>", "</code>")
}

func (formatNode *FormatNode) Spoiler() pages.IFormatter {
	return formatNode.envelopIn("<tg-spoiler>", "</tg-spoiler>")
}

func (formatNode *FormatNode) WriteTo(builder *strings.Builder) {
	builder.WriteString(formatNode.preamble)
	formatNode.content.WriteTo(builder)
	builder.WriteString(formatNode.epilogue)
}
