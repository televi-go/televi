package telegram

import "io"

type Request interface {
	Method() string
	Params() (Params, error)
}

type File struct {
	FieldName string
	Reader    io.Reader
	FileId    string
}

type RequestWithFiles interface {
	Request
	Files() []File
}
