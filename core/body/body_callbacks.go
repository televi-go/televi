package body

import (
	"fmt"
	"github.com/televi-go/televi/core/builders"
	"log"
)

type Callbacks struct {
	currentMessageIndex int
	triggers            map[string]func(ctx builders.ClickContext)
}

func NewCallbacks() *Callbacks {
	return &Callbacks{
		currentMessageIndex: 0,
		triggers:            map[string]func(ctx builders.ClickContext){},
	}
}

func (callbacks *Callbacks) beginNextMessage() {
	callbacks.currentMessageIndex++
}

func (callbacks *Callbacks) bind(descriptor string, callback func(ctx builders.ClickContext)) string {
	callbackData := fmt.Sprintf("_body_cb:%d:%s", callbacks.currentMessageIndex, descriptor)
	callbacks.triggers[callbackData] = callback
	return callbackData
}

func (callbacks *Callbacks) Execute(data string, ctx builders.ClickContext) {
	callback, hasCallback := callbacks.triggers[data]
	if hasCallback {
		callback(ctx)
		return
	}
	log.Printf("no callback for data %s", data)
}
