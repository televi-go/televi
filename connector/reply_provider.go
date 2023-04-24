package connector

import "github.com/televi-go/televi/models/render/results"

type replyProvider interface {
	BuildKeyboard() (result results.KeyboardResult, err error)
}
