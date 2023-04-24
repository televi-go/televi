package pages

import "github.com/televi-go/televi/telegram/dto"

type EventData struct {
	Kind    string
	Payload string
}

type ClickCallback = func()
type MessageCallback = func(message *dto.Message)
type ExternalCallback = func(data any)
type ContactCallback = func(contact *dto.Contact)

type Callbacks struct {
	Buttons   map[EventData][]ClickCallback
	OnMessage []MessageCallback
	External  map[string][]ExternalCallback
}

func NewCallbacks() *Callbacks {
	return &Callbacks{
		Buttons:  map[EventData][]ClickCallback{},
		External: map[string][]ExternalCallback{},
	}
}

func (callbacks *Callbacks) IsEmpty() bool {
	return len(callbacks.OnMessage) == 0 && len(callbacks.External) == 0 && len(callbacks.Buttons) == 0
}

func (callbacks *Callbacks) AddExternalListener(name string, callback ExternalCallback) {
	callbacks.External[name] = append(callbacks.External[name], callback)
}

func (callbacks *Callbacks) InvokeOnMessage(message *dto.Message) bool {
	if len(callbacks.OnMessage) == 0 {
		return false
	}
	for _, callback := range callbacks.OnMessage {
		callback(message)
	}
	return true
}

func (callbacks *Callbacks) InvokeExternal(event string, ctx ExternalReactionContext) bool {
	listeners := callbacks.External[event]
	if len(listeners) == 0 {
		return false
	}

	for _, listener := range listeners {
		listener(ctx)
	}

	return true
}

func (callbacks *Callbacks) AddButtonListener(data EventData, callback ClickCallback) {
	listener := callbacks.Buttons[data]
	listener = append(listener, callback)
	callbacks.Buttons[data] = listener
}

func (callbacks *Callbacks) AddMessageListener(callback MessageCallback) {
	callbacks.OnMessage = append(callbacks.OnMessage, callback)
}

// InvokeButton returns if any callbacks were triggered
func (callbacks *Callbacks) InvokeButton(data EventData) bool {
	listeners := callbacks.Buttons[data]
	if len(listeners) == 0 {
		return false
	}

	for _, listener := range listeners {
		listener()
	}

	return true
}

type CallbacksCall struct {
	Callbacks Callbacks
	Next      *CallbacksCall
}

func (c CallbacksCall) Call(call func(callbacks Callbacks) bool) bool {
	current := call(c.Callbacks)
	if current {
		return true
	}
	if c.Next == nil {
		return false
	}
	return c.Next.Call(call)
}
