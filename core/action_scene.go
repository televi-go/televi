package core

import (
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/telegram/dto"
)

type ActionScene interface {
	View(builder builders.Scene)
	Init()
	Dispose()
	OnMessage(message dto.Message)
}
