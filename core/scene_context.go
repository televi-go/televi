package core

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/body"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/callbacks"
	"github.com/televi-go/televi/core/head_management"
	"github.com/televi-go/televi/telegram/messages/keyboards"
)

type SceneBuildContext struct {
	HeadBuilder *head_management.HeadBuildContext
	*body.Root
}

func (s *SceneBuildContext) Head(builder func(headBuilder builders.Head)) {
	builder(s.HeadBuilder)
}

func (s *SceneBuildContext) Navigator() {
	//TODO implement me
	panic("implement me")
}

func NewSceneContext(root *body.Root) *SceneBuildContext {
	return &SceneBuildContext{
		Root: root,
		HeadBuilder: &head_management.HeadBuildContext{
			TextHtmlBuilder:     abstractions.TextHtmlBuilder{},
			TwoDimensionBuilder: abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]{},
			FMedia:              nil,
			Callbacks:           *callbacks.NewCallbacks(),
		},
	}
}
