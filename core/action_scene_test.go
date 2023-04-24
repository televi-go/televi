package core

import (
	"github.com/televi-go/televi/core/builders"
	"testing"
)

type SomeScene struct {
}

func (s SomeScene) View(builder ActionSceneBuilder) {
	builder.Head(func(headBuilder builders.Head) {
		for i := 0; i < 1000; i++ {
			headBuilder.Text("abc").Bold()
		}
	})
}

func BenchmarkHead(b *testing.B) {

	for i := 0; i < b.N; i++ {
		scene := SomeScene{}
		ctx := &SceneBuildContext{}
		scene.View(ctx)
	}
}
