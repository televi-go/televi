package gopage

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"reflect"
)

type Component[T any] func(data T, ctx Context)

func (c Component[T]) Mount(ctx Context, data T) {
	c(data, ctx)
}

type RenderAction func(ctx Context)

func (r RenderAction) Mount(ctx Context) {
	r(ctx)
}

func Bind[T any](component Component[T], data T) RenderAction {
	return func(ctx Context) {
		component(data, ctx)
	}
}

type contextApplicable struct {
	Tasks []func(ctx Context)
}

func (c *contextApplicable) AddTask(task func(ctx Context)) {
	c.Tasks = append(c.Tasks, task)
}

func stringifyData[T any](data T) map[string]string {
	result := make(map[string]string)
	rv := reflect.ValueOf(data)
	rt := reflect.TypeOf(data)
	for i := 0; i < rt.NumField(); i++ {
		isExported := rt.Field(i).PkgPath == ""
		if !isExported {
			continue
		}
		val := rv.Field(i).Interface()
		var marshal []byte
		s, isStr := val.(string)

		if isStr {
			marshal = []byte(s)
		} else {
			marshal, _ = json.Marshal(val)
		}
		if marshal != nil {
			result[rt.Field(i).Name] = string(marshal)
		}
	}
	return result
}

func componentFromApplicable[T any](applicable contextApplicable) Component[T] {
	return func(data T, ctx Context) {
		mapping := stringifyData[T](data)
		fragment := NewContext("")
		for _, task := range applicable.Tasks {
			task(&fragment)
		}
		node := fragment.rootNode
		node = node.PropagateData(mapping)
		ctx.AddChildTask(node)
	}
}

func MakeComponent[T any](source string) Component[T] {
	applicable := parseFragment(source)
	return componentFromApplicable[T](applicable)
}

func parseFragment(source string) contextApplicable {
	c := &contextApplicable{}
	z := html.NewTokenizer(bytes.NewReader([]byte(source)))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			isSelfClosing := tt == html.SelfClosingTagToken
			name, hasAttrs := z.TagName()
			var (
				key   []byte
				val   []byte
				attrs []Attr
			)
			for hasAttrs {
				key, val, hasAttrs = z.TagAttr()
				attrs = append(attrs, Attr{Key: string(key), Value: string(val)})
			}

			c.AddTask(func(ctx Context) {
				if isSelfClosing {
					ctx.OpenSelfClosing(string(name))
				} else {
					ctx.OpenTag(string(name))
				}
				ctx.Attributes(attrs...)
			})
			break
		case html.TextToken:
			text := z.Text()
			c.AddTask(func(ctx Context) {
				ctx.Content(string(text))
			})
			break
		case html.EndTagToken:
			c.AddTask(func(ctx Context) {
				ctx.CloseTag()
			})
			break
		}
	}
	return *c
}
