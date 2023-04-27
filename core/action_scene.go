package core

import (
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/external"
	"github.com/televi-go/televi/telegram/dto"
)

type ActionScene interface {
	View(builder builders.Scene)
	Init(ctx InitContext)
	Dispose()
	OnMessage(message dto.Message)
}

type InitContext interface {
	OnExternal(kind string, callback external.Callback)
}

type Navigator interface {
	Push(scene ActionScene)
	Pop()
}
