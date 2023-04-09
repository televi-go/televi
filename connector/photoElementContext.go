package connector

import (
	"fmt"
	"io"
	"televi/models/pages"
	"televi/models/render"
	"televi/models/render/results"
	"televi/telegram/messages/keyboards"
)

type photoElementContext struct {
	textElementContext
	PhotoProvider
}

type PhotoProvider struct {
	photoReader io.Reader
	photoFileId string
	key         string
	hasSpoiler  bool
}

func (elementContext *PhotoProvider) Spoiler() pages.ImageOptionsSetter {
	elementContext.hasSpoiler = true
	return elementContext
}

func (elementContext *PhotoProvider) Image(key string, source io.Reader) pages.ImageOptionsSetter {
	elementContext.key = key
	elementContext.photoReader = source
	return elementContext
}

func (elementContext *photoElementContext) BuildResult() (render.IResult, error) {

	var buttons [][]keyboards.InlineKeyboardButton
	if elementContext.keyboardBuilder != nil && elementContext.keyboardBuilder.Elements != nil {
		buttons = elementContext.keyboardBuilder.Elements
	}

	var err error

	if elementContext.photoReader == nil && elementContext.photoFileId == "" {
		err = fmt.Errorf("no image specified")
	}

	return &results.SingleMediaResult{
		Text:           elementContext.ToString(),
		ProtectContent: elementContext.protectContent,
		Silent:         elementContext.produceSilent,
		Key:            elementContext.key,
		FileId:         elementContext.photoFileId,
		FileReader:     elementContext.photoReader,
		ReplyMarkup:    &results.InlineKeyboardResult{Keyboard: buttons},
		Type:           "photo",
		HasSpoiler:     elementContext.hasSpoiler,
	}, err
}

type activePhotoContext struct {
	activeElementContext
	PhotoProvider
}

func (context *activePhotoContext) BuildResult() (render.IResult, error) {
	var buttons [][]keyboards.ReplyKeyboardButton
	if context.keyboardBuilder != nil && context.keyboardBuilder.Elements != nil {
		buttons = context.keyboardBuilder.Elements
	}

	var err error

	if context.photoReader == nil && context.photoFileId == "" {
		err = fmt.Errorf("no image specified")
	}

	return &results.SingleMediaResult{
		Text:           context.ToString(),
		ProtectContent: context.protectContent,
		Silent:         context.produceSilent,
		Key:            context.key,
		FileId:         context.photoFileId,
		FileReader:     context.photoReader,
		ReplyMarkup:    &results.ReplyKeyboardResult{Buttons: buttons},
		Type:           "photo",
		HasSpoiler:     context.hasSpoiler,
	}, err
}
