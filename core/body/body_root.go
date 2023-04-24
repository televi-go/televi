package body

import (
	"github.com/televi-go/televi/core/builders"
)

type NeedsPaintTask struct {
	SceneStateUpdated bool
}

type Root struct {
	//UpdateC chan<- struct{}
	Main *StatefulFragment
}

func (root *Root) SetNeedsUpdate() {
	root.Main.view = nil
}

type Wrapper struct {
	BodySetupFn func(builder builders.ComponentBuilder)
}

func NewRoot(signalChannel chan<- NeedsPaintTask) *Root {
	root := &Root{
		Main: NewStatefulFragment(
			nil,
			signalChannel,
		),
	}
	return root
}

// ReplaceView is called when root state is updated
func (root *Root) ReplaceView(body func(builder builders.ComponentBuilder)) {
	root.Main.view = Wrapper{BodySetupFn: body}
	root.Main.StateHasChanged = true
}

func (root *Root) Body(builder func(builder builders.ComponentBuilder)) {
	if root.Main.view == nil {
		root.Main.view = Wrapper{
			BodySetupFn: builder,
		}
		root.Main.StateHasChanged = true
	}
}

func (w Wrapper) Init() {
}

func (w Wrapper) View(builder builders.ComponentBuilder) {
	w.BodySetupFn(builder)
}

func (root *Root) ProvideResult() ([]Message, *Callbacks) {
	callbacks := NewCallbacks()
	componentBuilder := &FragmentBuilderImpl{
		Callbacks: callbacks,
	}
	root.Main.RunWith(componentBuilder)
	return componentBuilder.Messages, callbacks
}
