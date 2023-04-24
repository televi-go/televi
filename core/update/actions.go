package update

import (
	"github.com/televi-go/televi/telegram"
)

type MediaAction struct {
}

func (m *MediaAction) GetRequest(destination telegram.Destination) telegram.Request {
	panic("")
}

func (m *MediaAction) InflateResult(response telegram.Response) {

}
