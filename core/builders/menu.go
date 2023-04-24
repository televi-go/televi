package builders

import "github.com/televi-go/televi/telegram/dto"

type MenuRow interface {
	Button(caption string, onclick func())
	Contact(caption string, onclick func(contact dto.Contact))
	Location(caption string, onclick func(location dto.Location))
}

type Menu interface {
	Row(func(builder MenuRow))
}
