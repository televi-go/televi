package media

import "fmt"

type Kind struct {
	value string
}

var (
	// NoMedia is to be used only in comparisons
	NoMedia   = Kind{value: ""}
	ImageKind = Kind{value: "Photo"}
	VideoKind = Kind{value: "Video"}
)

func (kind Kind) Method() string {
	return fmt.Sprintf("send%s", kind.value)
}

func (kind Kind) FieldName() string {
	// TODO: add more logic
	return kind.value
}
