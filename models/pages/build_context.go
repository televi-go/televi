package pages

import (
	"github.com/televi-go/televi/telegram/dto"
	"io"
)

type TextPartBuilder interface {
	Text(value string) IFormatter
	TextF(value string, args ...any) IFormatter
	TextLine(value string) IFormatter
}

type PageBuildContext interface {
	TextElement(buildAction func(component TextContext))
	ActiveElement(buildAction func(component ActiveTextContext))
	PhotoElement(buildAction func(component PhotoContext))
	AnimationElement(buildAction func(component AnimationContext))
	ActiveAnimationElement(buildAction func(component ActiveAnimationContext))
	ActivePhoto(buildAction func(component ActivePhotoContext))
	GetUserId() int
	GetUserInfo() *dto.User
	GetNavigator() Navigator
}

type PhotoConsumer interface {
	Image(key string, source io.Reader) ImageOptionsSetter
}

type AnimationConsumer interface {
	Animation(key string, source io.Reader, filename string)
}

type AnimationContext interface {
	TextContext
	AnimationConsumer
}

type ActiveAnimationContext interface {
	ActiveTextContext
	AnimationConsumer
}

type PhotoContext interface {
	TextContext
	PhotoConsumer
}

type ActiveTextContext interface {
	TextPartBuilder
	ReplyKeyboard(buildAction func(builder ReplyKeyboardBuilder))
}

type ActivePhotoContext interface {
	ActiveTextContext
	PhotoConsumer
}

type ImageOptionsSetter interface {
	Spoiler() ImageOptionsSetter
}

type ReplyKeyboardBuilder ReplyBuilder[ReplyRowBuilder]
type InlineKeyboardBuilder ReplyBuilder[InlineRowBuilder]

type ReplyBuilder[T any] interface {
	ButtonsRow(builder func(rowBuilder T))
}

type RowBuilder interface {
	ActionButton(caption string, callback ClickCallback)
}

type MessageOriginatedRowBuilder interface {
	ActionButton(caption string, callback ClickCallback)
}

type ReplyRowBuilder interface {
	MessageOriginatedRowBuilder
	ContactButton(caption string, callback ContactCallback)
}

type InlineRowBuilder interface {
	RowBuilder
	UrlButton(caption string, url string)
}

type IFormatter interface {
	Bold() IFormatter
	Mono() IFormatter
	Spoiler() IFormatter
	//TODO: to be added
}

type TextContext interface {
	TextPartBuilder
	ButtonsRow(builder func(rowBuilder InlineRowBuilder))
	InlineKeyboard(buildAction func(builder InlineKeyboardBuilder))
}
