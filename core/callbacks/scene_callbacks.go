package callbacks

import (
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/util"
)

type MenuCallbacks struct {
	OnRegularButton map[string]func()
	OnContact       util.Option[func(contact dto.Contact)]
}

func NewCallbacks() *MenuCallbacks {
	return &MenuCallbacks{
		OnContact:       util.OptionEmpty[func(contact dto.Contact)](),
		OnRegularButton: map[string]func(){},
	}
}

func (callbacks *MenuCallbacks) BindButton(text string, callback func()) {
	callbacks.OnRegularButton[text] = callback
}
