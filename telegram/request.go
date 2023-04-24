package telegram

import "io"

type Request interface {
	Method() string
	Params() (Params, error)
}

type File struct {
	FieldName   string
	ContentType string
	Reader      io.Reader
	FileId      string
	Name        string
}

type RequestWithFiles interface {
	Request
	Files() []File
}
