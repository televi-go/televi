package core

import "github.com/televi-go/televi/telegram/dto"

type NavProvider interface {
	Push(scene ActionScene)
	Replace(scene ActionScene)
	PopScene()
	PopToMain()
}

type Platform interface {
	NavProvider
	GetUser() dto.User
}
