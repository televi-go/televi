package head_management

import (
	"github.com/televi-go/televi/abstractions"
	"github.com/televi-go/televi/core/builders"
	"github.com/televi-go/televi/core/callbacks"
	"github.com/televi-go/televi/core/media"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram/dto"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"github.com/televi-go/televi/util"
)

type HeadBuildContext struct {
	abstractions.TextHtmlBuilder
	abstractions.TwoDimensionBuilder[keyboards.ReplyKeyboardButton]
	FMedia         *media.Media
	ProtectContent bool
	Callbacks      callbacks.MenuCallbacks
}

func (headBuildContext *HeadBuildContext) SetProtection() {
	headBuildContext.ProtectContent = true
}

func (headBuildContext *HeadBuildContext) Button(caption string, onclick func()) {
	headBuildContext.Add(keyboards.ReplyKeyboardButton{
		Text:            caption,
		RequestContact:  false,
		RequestLocation: false,
	})
	headBuildContext.Callbacks.BindButton(caption, onclick)
}

func (headBuildContext *HeadBuildContext) Contact(caption string, onclick func(contact dto.Contact)) {
	headBuildContext.Add(keyboards.ReplyKeyboardButton{
		Text:            caption,
		RequestContact:  true,
		RequestLocation: false,
	})
	headBuildContext.Callbacks.OnContact = util.OptionValue(onclick)
}

func (headBuildContext *HeadBuildContext) Location(caption string, onclick func(location dto.Location)) {
	//TODO implement me
	panic("implement me")
}

func (headBuildContext *HeadBuildContext) Build() HeadResult {
	return HeadResult{
		Media:          headBuildContext.FMedia,
		Text:           headBuildContext.ToString(),
		ProtectContent: headBuildContext.ProtectContent,
		Keyboard:       results.ReplyKeyboardResult{Buttons: headBuildContext.Elements},
	}
}

func (headBuildContext *HeadBuildContext) Media(media2 media.Media) {
	headBuildContext.FMedia = &media2
}

func (headBuildContext *HeadBuildContext) Row(f func(builder builders.MenuRow)) {
	headBuildContext.CommitRow()
	f(headBuildContext)
	headBuildContext.CommitRow()
}
