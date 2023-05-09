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
	RegisterAction(domain string, action string)
	GetUser() dto.User
}
