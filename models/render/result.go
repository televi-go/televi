package render

import "github.com/televi-go/televi/telegram"

type IResult interface {
	Kind() string
	InitAction(destination telegram.Destination) telegram.Request
	// CompareTo takes result as newer
	CompareTo(result IResult, destination telegram.Destination, messageIds []int) (bool, []telegram.Request)
}

type IResultProvider interface {
	BuildResult() (IResult, error)
}
