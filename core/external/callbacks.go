package external

import (
	"log"
	"runtime"
)

type Callback = func(data any)

type CallbackDispatcher interface {
	Dispatch(name string, data any)
}

type emptyDispatcher struct{}

func (emptyDispatcher emptyDispatcher) Dispatch(name string, data any) {}

var EmptyDispatcher CallbackDispatcher = emptyDispatcher{}

type DispatcherBuilder struct {
	collection map[string][]Callback
}

func NewDispatcherBuilder() *DispatcherBuilder {
	return &DispatcherBuilder{collection: map[string][]Callback{}}
}

func (builder *DispatcherBuilder) OnExternal(kind string, callback Callback) {
	builder.collection[kind] = append(builder.collection[kind], callback)
}

func wrapCall(call Callback, data any) {
	defer func() {
		if v := recover(); v != nil {
			var buf = make([]byte, 2048)
			runtime.Stack(buf, false)
			log.Printf("Error in calling %v\n%s\n", v, buf)
		}
	}()
	call(data)
}

func (builder *DispatcherBuilder) Dispatch(kind string, data any) {
	for _, callback := range builder.collection[kind] {
		wrapCall(callback, data)
	}
}
