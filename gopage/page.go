package gopage

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Page = func(context Context)

type HtmlWriter struct {
	ctx     *fiber.Ctx
	context ContextImpl
}

func (writer *HtmlWriter) WritePage(p Page) error {
	p(&writer.context)
	builder := strings.Builder{}
	writer.context.rootNode.Render(&builder, 0)
	writer.ctx.Set("content-type", "text/html; charset='utf-8'")
	return writer.ctx.SendString(builder.String())
}

func NewContext(tag string) ContextImpl {
	rootNode := &NodeBuildTask{
		Tag: tag,
	}
	return ContextImpl{rootNode: rootNode, current: rootNode}
}

func NewHtmlWriter(ctx *fiber.Ctx) HtmlWriter {

	return HtmlWriter{
		ctx:     ctx,
		context: NewContext("html"),
	}
}
