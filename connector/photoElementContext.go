package connector

import (
	"fmt"
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/models/render"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"io"
)

type singleMediaContext struct {
	textElementContext
	SingleMediaProvider
}

type SingleMediaProvider struct {
	reader     io.Reader
	fileId     string
	MediaType  string
	key        string
	hasSpoiler bool
	filename   string
}

func (provider *SingleMediaProvider) singleMedia(key string, source io.Reader, filename string) {
	provider.key = key
	provider.reader = source
	provider.filename = filename
}

func (provider *SingleMediaProvider) Animation(key string, source io.Reader, filename string) {
	provider.singleMedia(key, source, filename)
}

func (provider *SingleMediaProvider) Spoiler() pages.ImageOptionsSetter {
	provider.hasSpoiler = true
	return provider
}

func (provider *SingleMediaProvider) Image(key string, source io.Reader) pages.ImageOptionsSetter {
	provider.singleMedia(key, source, "")
	return provider
}

func (elementContext *singleMediaContext) buildSingleMediaInternal() (*results.SingleMediaResult, error) {
	var err error

	kbResult, _ := elementContext.keyboardBuilder.BuildKeyboard()

	if elementContext.reader == nil && elementContext.fileId == "" {
		return nil, err
	}
	return &results.SingleMediaResult{
		Text:           elementContext.ToString(),
		ProtectContent: elementContext.protectContent,
		Silent:         elementContext.produceSilent,
		Key:            elementContext.key,
		FileId:         elementContext.fileId,
		FileReader:     elementContext.reader,
		ReplyMarkup:    kbResult,
		Type:           elementContext.MediaType,
		HasSpoiler:     elementContext.hasSpoiler,
		FileName:       elementContext.filename,
	}, nil
}

func (elementContext *singleMediaContext) BuildResult() (render.IResult, error) {
	return elementContext.buildSingleMediaInternal()
}

type activePhotoContext struct {
	activeElementContext
	SingleMediaProvider
}

func (context *activePhotoContext) BuildResult() (render.IResult, error) {
	var buttons [][]keyboards.ReplyKeyboardButton
	if context.keyboardBuilder != nil && context.keyboardBuilder.Elements != nil {
		buttons = context.keyboardBuilder.Elements
	}

	var err error

	if context.reader == nil && context.fileId == "" {
		err = fmt.Errorf("no image specified")
	}

	return &results.SingleMediaResult{
		Text:           context.ToString(),
		ProtectContent: context.protectContent,
		Silent:         context.produceSilent,
		Key:            context.key,
		FileId:         context.fileId,
		FileReader:     context.reader,
		ReplyMarkup:    &results.ReplyKeyboardResult{Buttons: buttons},
		Type:           context.MediaType,
		HasSpoiler:     context.hasSpoiler,
		FileName:       context.filename,
	}, err
}
