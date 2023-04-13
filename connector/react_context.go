package connector

import (
	"fmt"
	"gtihub.com/televi-go/televi/connector/abstractions"
	"gtihub.com/televi-go/televi/models/pages"
	"gtihub.com/televi-go/televi/telegram/dto"
	"gtihub.com/televi-go/televi/telegram/messages"
)

type reactContextImpl struct {
	controller             *Controller
	message                *dto.Message
	WasTransitRequested    bool
	AlertRequest           messages.AnswerCallbackRequest
	includeTransitToAnchor bool
	origin                 pages.ViewSequenceOrigin
}

func (reactionContext *reactContextImpl) Contact() *dto.Contact {
	return reactionContext.message.Contact
}

func (reactionContext *reactContextImpl) transitToAnchorModel() {
	if reactionContext.controller.ActiveCallbacksModel == nil {
		return
	}
	backTransitionsCount := 0
	backTransitModel := reactionContext.controller.CurrentModel
	for backTransitModel != nil && backTransitModel != reactionContext.controller.ActiveCallbacksModel {
		backTransitionsCount++
		backTransitModel = backTransitModel.Previous
	}
	if backTransitModel == nil {
		fmt.Println("Cannot transit to anchor model as it was possibly lost in stack")
		return
	}
	for i := 0; i < backTransitionsCount; i++ {
		reactionContext.controller._transitBack()
	}
}

func (reactionContext *reactContextImpl) TransitToMain() bool {
	reactionContext.controller.transitToMain()
	reactionContext.WasTransitRequested = true
	return true
}

func (reactionContext *reactContextImpl) ShowAlert(alertBuilder func(builder pages.TextPartBuilder)) {
	textBuilder := &abstractions.TextHtmlBuilder{}
	alertBuilder(textBuilder)
	reactionContext.AlertRequest = messages.AnswerCallbackRequest{
		Id:        reactionContext.AlertRequest.Id,
		Text:      textBuilder.ToString(),
		ShowAlert: true,
	}
}

func (reactionContext *reactContextImpl) Message() *dto.Message {
	return reactionContext.message
}

func (reactionContext *reactContextImpl) TransitTo(page pages.Scene, policy pages.TransitPolicy) {
	if reactionContext.includeTransitToAnchor {
		reactionContext.transitToAnchorModel()
	}
	reactionContext.controller.transitTo(page, policy, reactionContext.origin)
	reactionContext.WasTransitRequested = true
}

func (reactionContext *reactContextImpl) TransitBack() bool {
	canTransit := reactionContext.controller.CurrentModel.Previous != nil
	if !canTransit {
		return false
	}
	reactionContext.controller.transitBack()
	reactionContext.WasTransitRequested = true
	return true
}
