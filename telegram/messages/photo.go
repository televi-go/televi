package messages

import (
	"bytes"
	"fmt"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"io"
	"strings"
)

type MediaMessageBase struct {
	Destination    telegram.Destination
	Caption        string
	ProtectContent bool
	Silent         bool
	ReplyTo        int
	ReplyMarkup    keyboards.ReplyMarkup
}

func (mediaMessageBase MediaMessageBase) WriteParameter(params telegram.Params) error {
	var err error
	mediaMessageBase.Destination.WriteParameter(params)
	params.WriteString("caption", mediaMessageBase.Caption)
	params.WriteString("parse_mode", "HTML")
	params.WriteBool("protect_content", mediaMessageBase.ProtectContent)
	params.WriteBool("disable_notification", mediaMessageBase.Silent)
	params.WriteNonZero("reply_to_message_id", mediaMessageBase.ReplyTo)
	if mediaMessageBase.ReplyMarkup != nil {
		err = mediaMessageBase.ReplyMarkup.WriteParameter(params)
	}
	return err
}

type SingleMediaRequest struct {
	Base        MediaMessageBase
	Content     []byte
	FileName    string
	PhotoFileId string
	HasSpoiler  bool
	MediaType   string
}

func (sendPhotoRequest SingleMediaRequest) Method() string {
	return fmt.Sprintf("send%s", sendPhotoRequest.MediaType)
}

func (sendPhotoRequest SingleMediaRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)
	params.WriteBool("has_spoiler", sendPhotoRequest.HasSpoiler)
	err := sendPhotoRequest.Base.WriteParameter(params)
	return params, err
}

func ToSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

func (sendPhotoRequest SingleMediaRequest) Files() []telegram.File {
	var reader io.Reader
	if sendPhotoRequest.Content != nil {
		reader = bytes.NewReader(sendPhotoRequest.Content)
	}
	return []telegram.File{
		{
			FieldName: ToSnake(sendPhotoRequest.MediaType),
			Reader:    reader,
			FileId:    sendPhotoRequest.PhotoFileId,
			Name:      sendPhotoRequest.FileName,
		},
	}
}

type UpdateMediaRequest struct {
	EditMessageCaptionRequest
	Media InputMedia
}

func (updateMedia UpdateMediaRequest) Method() string {
	return "editMessageMedia"
}

func (updateMedia UpdateMediaRequest) Params() (telegram.Params, error) {
	params, err := updateMedia.EditInlineKeyboardRequest.Params()
	if err != nil {
		return nil, err
	}
	updateMedia.Caption = updateMedia.EditMessageCaptionRequest.Caption
	err = params.WriteJson("media", updateMedia.Media)
	return params, err
}

func (updateMedia UpdateMediaRequest) Files() []telegram.File {
	return []telegram.File{
		{
			FieldName: "media",
			Reader:    updateMedia.Media.Raw,
			FileId:    updateMedia.Media.FileId,
		},
	}
}

type InputMedia struct {
	Raw    io.Reader
	FileId string
	Type   string `json:"type"`
	// Caption is auto-supplied
	Caption string `json:"caption,omitempty"`
}

type EditMessageCaptionRequest struct {
	keyboards.EditInlineKeyboardRequest
	Caption string
}

func (editMessageCaption EditMessageCaptionRequest) Method() string {
	return "editMessageCaption"
}

func (editMessageCaption EditMessageCaptionRequest) Params() (telegram.Params, error) {
	params, err := editMessageCaption.EditInlineKeyboardRequest.Params()
	if err != nil {
		return nil, err
	}
	params.WriteString("caption", editMessageCaption.Caption)
	return params, nil
}
