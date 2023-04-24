package media

import (
	"os"
)

func File(filename string, kind Kind) Media {
	bytes, _ := os.ReadFile(filename)
	return Bytes(bytes, filename, kind)
}

func ImageFile(builder Insertable, filename string) {
	builder.Media(File(filename, ImageKind))
}

func Bytes(source []byte, key string, kind Kind) Media {
	return Media{
		Kind:       kind,
		Content:    source,
		FileId:     "",
		Key:        key,
		HasSpoiler: false,
	}
}
