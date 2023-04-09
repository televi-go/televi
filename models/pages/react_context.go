package pages

import (
	"gtihub.com/televi-go/televi/telegram/dto"
)

type TransitPolicy struct {
	KeepPrevious bool
}

func (policy TransitPolicy) GetKind() TransitionKind {
	if policy.KeepPrevious {
		return ReplacingTransition
	}
	return SeparativeTransition
}

type ReactionContext interface {
	TransitTo(page Scene, policy TransitPolicy)
	TransitBack() bool
	TransitToMain() bool
	ShowAlert(alertBuilder func(builder TextPartBuilder))
}

type MessageReactionContext interface {
	ReactionContext
	Message() *dto.Message
}

type ContactReactionContext interface {
	ReactionContext
	Contact() *dto.Contact
}

type ExternalReactionContext interface {
	ReactionContext
	EventData() any
}
