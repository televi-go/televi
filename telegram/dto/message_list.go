package dto

import "encoding/json"

type MessageList []Message

func (ml *MessageList) UnmarshalJSON(data []byte) error {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err == nil {
		*ml = []Message{msg}
		return nil
	}
	if string(data) == "true" {
		*ml = nil
		return nil
	}
	var msgs []Message
	err = json.Unmarshal(data, &msgs)
	if err == nil {
		*ml = msgs
	}
	return nil
}

func (ml *MessageList) CollectIds() []int {
	ids := make([]int, len(*ml))
	for i, msg := range *ml {
		ids[i] = msg.MessageID
	}
	return ids
}
