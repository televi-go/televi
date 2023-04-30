package gopage

import (
	"fmt"
	"strings"
	"testing"
)

const fragment = `
<div>
	<inner src="some-attr-val" href="a">{A}</inner>
	<br/>
</div>
`

func TestParseComponent(t *testing.T) {
	applicable := parseFragment(fragment)
	ctx := NewContext("")
	for _, task := range applicable.Tasks {
		task(&ctx)
	}
	builder := strings.Builder{}
	ctx.rootNode.Render(&builder, 0)
	println(builder.String())
}

type SampleViewData struct {
	A int64
}

func TestMakeComponent(t *testing.T) {
	component := MakeComponent[SampleViewData](fragment)
	ctx := NewContext("")
	component(SampleViewData{A: 2}, &ctx)
	builder := strings.Builder{}
	ctx.rootNode.Render(&builder, 0)
	println(builder.String())
}

func TestTokenize(t *testing.T) {
	str := "{A} {B}"
	tokens := tokenize(str)
	fmt.Printf("%#v", tokens)
}
