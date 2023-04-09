package messages

import (
	"televi/telegram"
)

type AnswerCallbackRequest struct {
	Id        string
	Text      string
	ShowAlert bool
}

func (answerCallbackRequest AnswerCallbackRequest) Method() string {
	return "answerCallbackQuery"
}

func (answerCallbackRequest AnswerCallbackRequest) Params() (telegram.Params, error) {
	params := make(telegram.Params)
	params.WriteString("text", answerCallbackRequest.Text)
	params.WriteString("callback_query_id", answerCallbackRequest.Id)
	params.WriteBool("show_alert", answerCallbackRequest.ShowAlert)
	return params, nil
}
