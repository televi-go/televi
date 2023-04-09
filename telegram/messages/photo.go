package messages

import (
	"gtihub.com/televi-go/televi/telegram"
	"gtihub.com/televi-go/televi/telegram/messages/keyboards"
	"io"
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

type SendPhotoRequest struct {
	Base        MediaMessageBase
	Photo       io.Reader
	PhotoFileId string
	HasSpoiler  bool
}

func (sendPhotoRequest SendPhotoRequest) Method() string {
	return "sendPhoto"
}

func (sendPhotoRequest SendPhotoRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)
	params.WriteBool("has_spoiler", sendPhotoRequest.HasSpoiler)
	err := sendPhotoRequest.Base.WriteParameter(params)
	return params, err
}

func (sendPhotoRequest SendPhotoRequest) Files() []telegram.File {
	return []telegram.File{
		{
			FieldName: "photo",
			Reader:    sendPhotoRequest.Photo,
			FileId:    sendPhotoRequest.PhotoFileId,
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
