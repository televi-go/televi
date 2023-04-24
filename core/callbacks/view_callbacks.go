package callbacks

import "fmt"

type ViewCallbacksProvider interface {
	Bind(descriptor string, callback func())
}

type ViewCallbacks struct {
	Parent ViewCallbacksProvider
	Prefix string
}

func (viewCallbacks *ViewCallbacks) Bind(descriptor string, callback func()) {
	viewCallbacks.Parent.Bind(fmt.Sprintf("%s.%s", viewCallbacks.Prefix, descriptor), callback)
}
