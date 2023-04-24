package body

import (
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
)

type RealComponentBuilder interface {
	builders.ComponentBuilder
	SetIsModified(value bool)
}

type FragmentRunner interface {
	RunWith(builder builders.ComponentBuilder) (isCached bool)
}

type FragmentProxyBuilder struct {
	BuildTasks []FragmentRunner
	Callbacks  *Callbacks
	UpdateC    chan<- NeedsPaintTask
}

type StatefulFragment struct {
	view            builders.View
	StateHasChanged bool
	Descendants     []FragmentRunner
	Callbacks       *Callbacks
	UpdateC         chan<- NeedsPaintTask
}

func (statefulFragment *StatefulFragment) RunWith(builder builders.ComponentBuilder) (isCached bool) {
	isCached = !statefulFragment.StateHasChanged

	if statefulFragment.StateHasChanged {

		disposableView, isDisposableView := statefulFragment.view.(builders.DisposableView)
		if isDisposableView {

			disposableView.Dispose()
		}

		statefulFragment.Descendants = nil
		proxy := &FragmentProxyBuilder{
			UpdateC:   statefulFragment.UpdateC,
			Callbacks: statefulFragment.Callbacks,
		}
		statefulFragment.view.View(proxy)
		statefulFragment.Descendants = proxy.BuildTasks
		statefulFragment.StateHasChanged = false
	}
	for _, buildTask := range statefulFragment.Descendants {
		buildTask.RunWith(builder)
	}
	return
}

func NewStatefulFragment(v builders.View, c chan<- NeedsPaintTask) *StatefulFragment {
	fragment := &StatefulFragment{
		view:            v,
		StateHasChanged: true,
		UpdateC:         c,
	}
	if v != nil {
		magic.InjectInPlace(v, func() {
			fragment.StateHasChanged = true
			c <- NeedsPaintTask{SceneStateUpdated: false}
		})
		go v.Init()
	}

	return fragment
}

type StatelessFragment struct {
	Setup   func(builder builders.ComponentBuilder)
	UpdateC chan<- NeedsPaintTask
}

func (statelessFragment StatelessFragment) RunWith(builder builders.ComponentBuilder) (isCached bool) {
	statelessFragment.Setup(builder)
	return true
}

func (fragmentProxy *FragmentProxyBuilder) Component(view builders.View) {
	fragmentProxy.BuildTasks = append(fragmentProxy.BuildTasks, NewStatefulFragment(
		view,
		fragmentProxy.UpdateC,
	))
}

type proxyStub struct {
	builders.Message
}

func (p proxyStub) SetIsModified(value bool) {}

func (fragmentProxy *FragmentProxyBuilder) Message(builder func(viewBuilder builders.Message)) {

	proxyMessageBuilder := &proxyMessageBuilderImpl{StateUpdateChan: fragmentProxy.UpdateC}
	// collect descendant nodes
	builder(proxyMessageBuilder)

	fragmentProxy.BuildTasks = append(fragmentProxy.BuildTasks, StatelessFragment{
		Setup: func(component builders.ComponentBuilder) {
			component.Message(func(viewBuilder builders.Message) {
				for _, descendant := range proxyMessageBuilder.Descendants {
					descendant.RunWith(proxyStub{viewBuilder})
				}
			})
		},
		UpdateC: fragmentProxy.UpdateC,
	})
}
