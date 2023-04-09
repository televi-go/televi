package telegram

import "encoding/json"

type Response struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

func ParseAs[T any](response Response) (T, error) {
	var t T
	err := json.Unmarshal(response.Result, &t)
	return t, err
}
