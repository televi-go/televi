package bot

import (
	"github.com/televi-go/televi/telegram"
)

type GetMeRequest struct{}

func (getMe GetMeRequest) Method() string {
	return "getMe"
}

func (getMe GetMeRequest) Params() (telegram.Params, error) {
	return map[string]string{}, nil
}
