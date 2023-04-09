package connector

import (
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/models/render"
	"strconv"
)

type BuildContext struct {
	elements        []render.IResultProvider
	everySilent     bool
	everyProtected  bool
	Callbacks       *pages.Callbacks
	ActiveCallbacks *pages.Callbacks
	UserId          int
}

func (p *BuildContext) PhotoElement(buildAction func(component pages.PhotoContext)) {
	context := &photoElementContext{
		textElementContext: textElementContext{
			componentPrefix: strconv.Itoa(len(p.elements)),
			produceSilent:   p.everySilent,
			protectContent:  p.everyProtected,
			callbacks:       p.Callbacks,
		},
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) buildLine() []render.IResult {
	var result []render.IResult
	for _, element := range p.elements {
		r, err := element.BuildResult()
		if err != nil {
			panic(err)
		}
		result = append(result, r)
	}
	return result
}

func (p *BuildContext) GetUserId() int {
	return p.UserId
}

func (p *BuildContext) TextElement(buildAction func(ctx pages.TextContext)) {
	context := &textElementContext{
		componentPrefix: strconv.Itoa(len(p.elements)),
		produceSilent:   p.everySilent,
		protectContent:  p.everyProtected,
		callbacks:       p.Callbacks,
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) ActiveElement(buildAction func(ctx pages.ActiveTextContext)) {
	context := &activeElementContext{
		Callbacks:      p.ActiveCallbacks,
		produceSilent:  p.everySilent,
		protectContent: p.everyProtected,
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) ActivePhoto(buildAction func(ctx pages.ActivePhotoContext)) {
	context := &activePhotoContext{
		activeElementContext: activeElementContext{
			Callbacks:      p.ActiveCallbacks,
			produceSilent:  p.everySilent,
			protectContent: p.everyProtected,
		},
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}
