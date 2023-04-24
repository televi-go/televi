package views

import (
	"github.com/televi-go/televi/core/builders"
)

type listView[T any] struct {
	sourceData []T
	builder    func(data T) builders.View
}

func (listView listView[T]) Init() {}

func (listView listView[T]) View(builder builders.ComponentBuilder) {
	for _, datum := range listView.sourceData {
		datum := datum
		builder.Component(listView.builder(datum))
	}
}

func ListView[T any](context builders.ComponentBuilder, data []T, builder func(data T) builders.View) {
	context.Component(listView[T]{
		sourceData: data,
		builder:    builder,
	})
}

func ListFuncView[T any](context builders.ComponentBuilder, data []T, funcBuilder func(item T, componentBuilder builders.ComponentBuilder)) {
	ListView(context, data, func(data T) builders.View {
		return FuncView(func(builder builders.ComponentBuilder) {
			funcBuilder(data, builder)
		})
	})
}
