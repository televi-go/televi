package builders

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/models/pages"
)

type ContentBuilder interface {
	Text(value string) pages.IFormatter
	AddBuildable(buildable abstractions.Buildable)
	TextLine(value string) pages.IFormatter
	TextF(format string, values ...any) pages.IFormatter
}

type Head interface {
	ContentBuilder
	media.Insertable
	SetProtection()
	Menu
}

type InMessageView interface {
	View(builder Message)
	Init()
}

type Message interface {
	media.Insertable
	ContentBuilder
	ActionsBuilder
	SetProtection()
	Component(view InMessageView)
}
