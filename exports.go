package televi

import "televi/models/pages"

type Scene interface {
	View(ctx BuildContext)
}
type BuildContext = pages.PageBuildContext
type TransitPolicy = pages.TransitPolicy

func ForEach[T any](data []T, runner func(element T)) {
	for _, datum := range data {
		runner(datum)
	}
}
