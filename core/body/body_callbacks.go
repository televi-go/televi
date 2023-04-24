package body

import (
	"fmt"
	"log"
)

type Callbacks struct {
	currentMessageIndex int
	triggers            map[string]func()
}

func NewCallbacks() *Callbacks {
	return &Callbacks{
		currentMessageIndex: 0,
		triggers:            map[string]func(){},
	}
}

func (callbacks *Callbacks) beginNextMessage() {
	callbacks.currentMessageIndex++
}

func (callbacks *Callbacks) bind(descriptor string, callback func()) string {
	callbackData := fmt.Sprintf("_body_cb:%d:%s", callbacks.currentMessageIndex, descriptor)
	callbacks.triggers[callbackData] = callback
	return callbackData
}

func (callbacks *Callbacks) Execute(data string) {
	callback, hasCallback := callbacks.triggers[data]
	if hasCallback {
		callback()
		return
	}
	log.Printf("no callback for data %s", data)
}
