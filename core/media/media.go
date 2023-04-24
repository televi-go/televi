package media

type Media struct {
	Kind
	Content    []byte
	FileId     string
	Key        string
	HasSpoiler bool
}

type Insertable interface {
	Media(media Media)
}
