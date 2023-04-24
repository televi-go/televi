package connector

import (
	"github.com/televi-go/televi/models/pages"
	"github.com/televi-go/televi/models/render"
	"github.com/televi-go/televi/telegram/dto"
	"strconv"
)

type BuildContext struct {
	elements        []render.IResultProvider
	everySilent     bool
	everyProtected  bool
	Callbacks       *pages.Callbacks
	ActiveCallbacks *pages.Callbacks
	UserId          int
	UserInfo        *dto.User
	controller      *Controller
	stackPoint      *pages.Model
}

func (p *BuildContext) AnimationElement(buildAction func(component pages.AnimationContext)) {
	context := &singleMediaContext{
		textElementContext: textElementContext{
			componentPrefix: strconv.Itoa(len(p.elements)),
			produceSilent:   p.everySilent,
			protectContent:  p.everyProtected,
			callbacks:       p.Callbacks,
		},
		SingleMediaProvider: SingleMediaProvider{MediaType: "Video"},
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) ActiveAnimationElement(buildAction func(component pages.ActiveAnimationContext)) {
	context := &activePhotoContext{
		activeElementContext: activeElementContext{
			Callbacks:      p.ActiveCallbacks,
			produceSilent:  p.everySilent,
			protectContent: p.everyProtected,
		},
		SingleMediaProvider: SingleMediaProvider{MediaType: "Video"},
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) PhotoElement(buildAction func(component pages.PhotoContext)) {
	context := &singleMediaContext{
		textElementContext: textElementContext{
			componentPrefix: strconv.Itoa(len(p.elements)),
			produceSilent:   p.everySilent,
			protectContent:  p.everyProtected,
			callbacks:       p.Callbacks,
		},
		SingleMediaProvider: SingleMediaProvider{MediaType: "Content"},
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
	return int(p.UserInfo.ID)
}

func (p *BuildContext) GetUserInfo() *dto.User {
	return p.UserInfo
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
		SingleMediaProvider: SingleMediaProvider{MediaType: "Content"},
	}
	buildAction(context)
	p.elements = append(p.elements, context)
}

func (p *BuildContext) GetNavigator() pages.Navigator {
	return navigator{controller: p.controller, stackPoint: p.stackPoint}
}
