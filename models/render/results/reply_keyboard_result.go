package results

import (
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"github.com/televi-go/televi/util"
)

type ReplyKeyboardResult struct {
	Buttons [][]keyboards.ReplyKeyboardButton
}

func (replyKeyboardResult ReplyKeyboardResult) Kind() string {
	return "reply"
}

func (replyKeyboardResult ReplyKeyboardResult) CanBeUpdated(kbResult KeyboardResult) UpdateAction {
	formerResult, notPointerErr := util.PointerOr[ReplyKeyboardResult](kbResult)
	if notPointerErr != nil {
		return ReplaceAction
	}

	if len(formerResult.Buttons) != len(replyKeyboardResult.Buttons) {
		return ReplaceAction
	}

	for i := 0; i < len(replyKeyboardResult.Buttons); i++ {
		if len(replyKeyboardResult.Buttons[i]) != len(formerResult.Buttons[i]) {
			return ReplaceAction
		}
		for j := 0; j < len(replyKeyboardResult.Buttons[i]); j++ {
			if replyKeyboardResult.Buttons[i][j] != formerResult.Buttons[i][j] {
				return ReplaceAction
			}
		}
	}

	return NoAction
}

func (replyKeyboardResult ReplyKeyboardResult) ToReplyMarkup() keyboards.ReplyMarkup {
	if len(replyKeyboardResult.Buttons) == 0 {
		return keyboards.ReplyKeyboardRemove{}
	}
	return keyboards.ReplyKeyboardMarkup{Keyboard: replyKeyboardResult.Buttons, ResizeKeyboard: true}
}
