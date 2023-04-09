package results

import "gtihub.com/televi-go/televi/telegram/messages/keyboards"

type KeyboardResult interface {
	Kind() string
	CanBeUpdated(kbResult KeyboardResult) UpdateAction
	ToReplyMarkup() keyboards.ReplyMarkup
}

type UpdateAction int

const (
	NoAction      UpdateAction = iota
	EditAction    UpdateAction = iota
	ReplaceAction UpdateAction = iota
)

type InlineKeyboardResult struct {
	Keyboard [][]keyboards.InlineKeyboardButton
}

func (inlineKeyboardResult *InlineKeyboardResult) Kind() string {
	return "inline"
}

func (inlineKeyboardResult *InlineKeyboardResult) CanBeUpdated(kbResult KeyboardResult) UpdateAction {
	inlineKb, isInlineKb := kbResult.(*InlineKeyboardResult)
	if !isInlineKb {
		return ReplaceAction
	}
	if len(inlineKb.Keyboard) != len(inlineKeyboardResult.Keyboard) {
		return EditAction
	}

	for i := 0; i < len(inlineKb.Keyboard); i++ {
		if len(inlineKb.Keyboard[i]) != len(inlineKeyboardResult.Keyboard[i]) {
			return EditAction
		}
		for j := 0; j < len(inlineKb.Keyboard[i]); j++ {
			if inlineKb.Keyboard[i][j] != inlineKeyboardResult.Keyboard[i][j] {
				return EditAction
			}
		}
	}

	return NoAction
}

func (inlineKeyboardResult *InlineKeyboardResult) ToReplyMarkup() keyboards.ReplyMarkup {
	return keyboards.InlineKeyboardMarkup{Keyboard: inlineKeyboardResult.Keyboard}
}
