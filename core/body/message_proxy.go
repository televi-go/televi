package body

import (
	"fmt"
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/models/pages"
)

type ProxyMessageBuilder interface {
	RunWith(builder RealMessageBuilder)
}

type StatelessMessageBuilder func(builder builders.Message)

func (statelessMessageBuilder StatelessMessageBuilder) RunWith(builder RealMessageBuilder) {
	statelessMessageBuilder(builder)
}

type StatefulMessageBuilder struct {
	Component         builders.InMessageView
	StateHasChanged   bool
	StateUpdateChan   chan<- NeedsPaintTask
	CachedDescendants []ProxyMessageBuilder
}

func NewStatefulMessageBuilder(
	component builders.InMessageView,
	stateUpdateChannel chan<- NeedsPaintTask,
) *StatefulMessageBuilder {

	builder := &StatefulMessageBuilder{
		Component:       component,
		StateHasChanged: true,
		StateUpdateChan: stateUpdateChannel,
	}
	magic.InjectInPlace(component, func() {
		builder.StateHasChanged = true
		stateUpdateChannel <- NeedsPaintTask{SceneStateUpdated: false}
	})
	go component.Init()
	return builder
}

func (statefulMessageBuilder *StatefulMessageBuilder) RunWith(builder RealMessageBuilder) {
	builder.SetIsModified(statefulMessageBuilder.StateHasChanged)
	if statefulMessageBuilder.StateHasChanged {

		statefulMessageBuilder.CachedDescendants = nil
		proxy := &proxyMessageBuilderImpl{
			StateUpdateChan: statefulMessageBuilder.StateUpdateChan,
		}

		statefulMessageBuilder.Component.View(proxy)

		// cache all actions
		statefulMessageBuilder.CachedDescendants = proxy.Descendants
		statefulMessageBuilder.StateHasChanged = false
	}
	for _, descendant := range statefulMessageBuilder.CachedDescendants {
		descendant.RunWith(builder)
	}
}

type proxyMessageBuilderImpl struct {
	Descendants     []ProxyMessageBuilder
	StateUpdateChan chan<- NeedsPaintTask
}

func (proxy *proxyMessageBuilderImpl) AddBuildable(buildable abstractions.Buildable) {
	proxy.Descendants = append(proxy.Descendants, StatelessMessageBuilder(func(builder builders.Message) {
		builder.AddBuildable(buildable)
	}))
}

func (proxy *proxyMessageBuilderImpl) Media(media media.Media) {
	proxy.Descendants = append(proxy.Descendants, StatelessMessageBuilder(func(builder builders.Message) {
		builder.Media(media)
	}))
}

func (proxy *proxyMessageBuilderImpl) Text(value string) pages.IFormatter {
	formatter := NewFormatNode(value)
	proxy.Descendants = append(proxy.Descendants, StatelessMessageBuilder(func(builder builders.Message) {
		builder.AddBuildable(formatter)
	}))
	return formatter
}
func (proxy *proxyMessageBuilderImpl) SetProtection() {
	proxy.Descendants = append(proxy.Descendants, StatelessMessageBuilder(func(builder builders.Message) {
		builder.SetProtection()
	}))
}

func (proxy *proxyMessageBuilderImpl) TextLine(value string) pages.IFormatter {
	return proxy.Text(value + "\n")
}

func (proxy *proxyMessageBuilderImpl) TextF(format string, values ...any) pages.IFormatter {
	text := fmt.Sprintf(format, values...)
	return proxy.Text(text)
}

func (proxy *proxyMessageBuilderImpl) Row(f func(builder builders.ActionRowBuilder)) {
	proxy.Descendants = append(proxy.Descendants, StatelessMessageBuilder(func(builder builders.Message) {
		builder.Row(f)
	}))
}

func (proxy *proxyMessageBuilderImpl) Component(view builders.InMessageView) {
	proxy.Descendants = append(proxy.Descendants, NewStatefulMessageBuilder(view, proxy.StateUpdateChan))
}
