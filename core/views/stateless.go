package views

import "github.com/televi-go/televi/core/builders"

type FuncView func(builder builders.ComponentBuilder)

func (f FuncView) Init() {

}

func (f FuncView) View(builder builders.ComponentBuilder) {
	f(builder)
}
