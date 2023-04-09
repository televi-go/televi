package render

import (
	"gtihub.com/televi-go/televi/telegram"
	"gtihub.com/televi-go/televi/telegram/messages"
	"time"
)

// CompletedResult
// LatestResult is nullable
type CompletedResult struct {
	MessageIds   []int
	LatestResult IResult
	MountedAt    time.Time
}

// CompareTo
//
// canBeChanged declares if there is a way of editing message
// changes are the changes can be made
// if returns false, there is a need to delete and resend
func (completedResult *CompletedResult) CompareTo(result IResult, destination telegram.Destination) (canBeChanged bool, changes []telegram.Request) {
	if completedResult.LatestResult == nil {
		return false, nil
	}

	if completedResult.LatestResult.Kind() != result.Kind() {
		return false, nil
	}

	return completedResult.LatestResult.CompareTo(result, destination, completedResult.MessageIds)
}

func (completedResult *CompletedResult) Cleanup(destination telegram.Destination) (result []telegram.Request) {
	for _, messageId := range completedResult.MessageIds {

		result = append(result, messages.DeleteMessageRequest{
			MessageId:   messageId,
			Destination: destination,
		})
	}
	return
}
