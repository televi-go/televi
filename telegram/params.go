package telegram

import (
	"encoding/json"
	"strconv"
)

type Params map[string]string

func (params Params) WriteInt(name string, value int) {
	params[name] = strconv.Itoa(value)
}

func (params Params) WriteBool(name string, value bool) {
	params[name] = strconv.FormatBool(value)
}

func (params Params) WriteNonZero(name string, value int) {
	if value != 0 {
		params.WriteInt(name, value)
	}
}

func (params Params) WriteJson(key string, value any) error {
	if value != nil {
		bytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		params[key] = string(bytes)
	}
	return nil
}

func (params Params) WriteString(name string, value string) {
	params[name] = value
}
