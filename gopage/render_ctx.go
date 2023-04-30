package gopage

import "strings"

type Attr struct {
	Key   string
	Value string
}

type Context interface {
	OpenTag(tag string)
	OpenSelfClosing(tag string)
	Attributes(attrs ...Attr)
	Content(content string)
	CloseTag()
	AddChildTask(task *NodeBuildTask)
}

func WriteAttribute(context Context, key string, value string) {
	context.Attributes(Attr{Key: key, Value: value})
}

type CanRender interface {
	Context
	Render(builder *strings.Builder)
}

type NodeBuildTask struct {
	WasDataPropagated bool
	Tag               string
	IsSelfClosing     bool
	Attributes        []Attr
	ChildContent      []*NodeBuildTask
	Raw               string
	Parent            *NodeBuildTask
}

type strTokKind int

const (
	rawKind   strTokKind = iota
	openKind  strTokKind = iota
	tokKind   strTokKind = iota
	closeKind strTokKind = iota
)

type strToken struct {
	TokName string
	Raw     string
	kind    strTokKind
}

func tokenize(in string) []strToken {
	result := make([]strToken, 0, 3)
	currentFragBegin := 0
	currentFragEnd := 0
	for i := 0; i < len(in); i++ {
		if in[i] == '{' {
			if currentFragBegin != currentFragEnd {
				result = append(result, strToken{kind: rawKind, Raw: in[currentFragBegin:currentFragEnd]})
			}
			result = append(result, strToken{kind: openKind})
			currentFragBegin = i + 1
			currentFragEnd = i + 1
			continue
		}

		if in[i] == '}' {
			result = append(result, strToken{kind: tokKind, TokName: in[currentFragBegin:currentFragEnd]})
			result = append(result, strToken{kind: closeKind})
			currentFragBegin = i + 1
			currentFragEnd = i + 1
			continue
		}

		currentFragEnd++
	}

	result = append(result, strToken{kind: rawKind, Raw: in[currentFragBegin:currentFragEnd]})

	return result
}

func formatTokens(tokens []strToken, data map[string]string) string {
	builder := strings.Builder{}
	for _, token := range tokens {
		switch token.kind {
		case rawKind:
			builder.WriteString(token.Raw)
			break
		case tokKind:
			datum, hasData := data[token.TokName]
			if hasData {
				builder.WriteString(datum)
			}
			break
		}
	}
	return builder.String()
}

func (node NodeBuildTask) PropagateData(data map[string]string) *NodeBuildTask {

	if node.WasDataPropagated {
		return &node
	}

	newTask := &NodeBuildTask{
		WasDataPropagated: true,
		Tag:               node.Tag,
		IsSelfClosing:     node.IsSelfClosing,
		Attributes:        nil,
		ChildContent:      nil,
		Raw:               "",
		Parent:            nil,
	}
	for _, attribute := range node.Attributes {
		attrTokens := tokenize(attribute.Value)
		newTask.Attributes = append(newTask.Attributes, Attr{Key: attribute.Key, Value: formatTokens(attrTokens, data)})
	}
	rawTokens := tokenize(node.Raw)
	newTask.Raw = formatTokens(rawTokens, data)
	for _, task := range node.ChildContent {
		child := task.PropagateData(data)
		child.Parent = newTask
		newTask.ChildContent = append(newTask.ChildContent, child)
	}
	return newTask
}

func writePadding(builder *strings.Builder, padding int) {
	for i := 0; i < padding; i++ {
		builder.WriteByte(' ')
	}
}

func (node NodeBuildTask) Render(builder *strings.Builder, depth int) {
	if node.Tag != "" {
		writePadding(builder, depth)
		builder.WriteByte('<')
		builder.WriteString(node.Tag)
		for _, attribute := range node.Attributes {
			builder.WriteByte(' ')
			builder.WriteString(attribute.Key)
			builder.WriteString("=\"")
			builder.WriteString(attribute.Value)
			builder.WriteByte('"')
		}
		builder.WriteByte('>')
	}
	lines := strings.Split(node.Raw, "\n")
	for _, line := range lines {
		writePadding(builder, depth+4)
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	for _, task := range node.ChildContent {
		task.Render(builder, depth+4)
	}
	if node.Tag != "" && !node.IsSelfClosing {
		writePadding(builder, depth)
		builder.WriteString("</")
		builder.WriteString(node.Tag)
		builder.WriteByte('>')
		builder.WriteByte('\n')
	}
}

type ContextImpl struct {
	rootNode *NodeBuildTask
	current  *NodeBuildTask
}

func (contextImpl *ContextImpl) AddChildTask(task *NodeBuildTask) {
	contextImpl.current.ChildContent = append(contextImpl.current.ChildContent, task)
	task.Parent = contextImpl.current
}

func (contextImpl *ContextImpl) OpenSelfClosing(tag string) {
	childNode := &NodeBuildTask{Tag: tag, Parent: contextImpl.current, IsSelfClosing: true}
	contextImpl.current.ChildContent = append(contextImpl.current.ChildContent, childNode)
	contextImpl.current = childNode
}

func (contextImpl *ContextImpl) OpenTag(tag string) {
	childNode := &NodeBuildTask{Tag: tag, Parent: contextImpl.current}
	contextImpl.current.ChildContent = append(contextImpl.current.ChildContent, childNode)
	contextImpl.current = childNode
}

func (contextImpl *ContextImpl) Attributes(attrs ...Attr) {
	contextImpl.current.Attributes = append(contextImpl.current.Attributes, attrs...)
}

func (contextImpl *ContextImpl) Content(content string) {
	contextImpl.current.ChildContent = append(contextImpl.current.ChildContent, &NodeBuildTask{
		Tag:          "",
		Attributes:   nil,
		ChildContent: nil,
		Raw:          content,
		Parent:       contextImpl.current,
	})
}

func (contextImpl *ContextImpl) CloseTag() {
	if contextImpl.current.Parent != nil {
		contextImpl.current = contextImpl.current.Parent
	}
}
