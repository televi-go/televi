package abstractions

type TwoDimensionBuilder[T any] struct {
	Elements   [][]T
	CurrentRow []T
}

func (dimensionBuilder *TwoDimensionBuilder[T]) Add(element T) {
	dimensionBuilder.CurrentRow = append(dimensionBuilder.CurrentRow, element)
}

func (dimensionBuilder *TwoDimensionBuilder[T]) CommitRow() {
	if len(dimensionBuilder.CurrentRow) != 0 {
		dimensionBuilder.Elements = append(dimensionBuilder.Elements, dimensionBuilder.CurrentRow)
	}
	dimensionBuilder.CurrentRow = nil
}
