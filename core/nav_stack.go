package core

import (
	"github.com/televi-go/televi/core/body"
	"github.com/televi-go/televi/core/external"
	"github.com/televi-go/televi/core/head_management"
	"github.com/televi-go/televi/core/magic"
	"github.com/televi-go/televi/models/render/results"
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/telegram/messages/keyboards"
	"log"
)

type NavStackEntry struct {
	ActionScene       ActionScene
	WasInitialized    bool
	Previous          *NavStackEntry
	Kind              NavigationKind
	Destination       telegram.Destination
	HeadResult        *head_management.HeadResultContainer
	BodyRoot          *body.Root
	BodyResult        *body.Result
	BodyCallbacks     *body.Callbacks
	ExternalCallbacks external.CallbackDispatcher
	UpdateC           chan<- body.NeedsPaintTask
}

func (entry *NavStackEntry) dehydrate(api *bot.Api) {
	for _, resultEntry := range entry.BodyResult.Entries {
		for _, messageId := range resultEntry.BoundIds {
			api.LaunchRequest(keyboards.EditInlineKeyboardRequest{
				Destination: entry.Destination,
				MessageId:   messageId,
				NewKeyboard: nil,
			})
		}
		resultEntry.Actions = results.InlineKeyboardResult{}
	}
}

func (entry *NavStackEntry) Render(needsRepaint bool, api *bot.Api) {
	if needsRepaint {
		entry.BodyRoot.SetNeedsUpdate()
	}
	context := NewSceneContext(entry.BodyRoot)
	entry.ActionScene.View(context)
	wasReplaced := false
	if needsRepaint {
		headResult := context.HeadBuilder.Build()
		wasReplaced = entry.HeadResult.CompareAgainst(headResult, entry.Destination, api)
		entry.HeadResult.HeadCallbacks = context.HeadBuilder.Callbacks
	}

	bodyResultLine, callbacks := context.Root.ProvideResult()
	entry.BodyResult.CompareAgainst(
		bodyResultLine,
		api,
		entry.Destination,
		wasReplaced,
	)
	entry.BodyCallbacks = callbacks
}

func NewNavStackEntry(
	scene ActionScene,
	stateUpdateChannel chan<- body.NeedsPaintTask,
	destination telegram.Destination,
	previous *NavStackEntry,
) *NavStackEntry {
	entry := &NavStackEntry{
		ActionScene:       scene,
		WasInitialized:    false,
		Previous:          previous,
		Kind:              0,
		Destination:       destination,
		HeadResult:        &head_management.HeadResultContainer{},
		BodyRoot:          body.NewRoot(stateUpdateChannel),
		BodyResult:        &body.Result{},
		BodyCallbacks:     body.NewCallbacks(),
		ExternalCallbacks: external.EmptyDispatcher,
		UpdateC:           stateUpdateChannel,
	}
	entry.initScene()
	return entry
}

func (entry *NavStackEntry) initScene() {
	if entry.WasInitialized {
		return
	}
	entry.WasInitialized = true
	magic.InjectInPlace(entry.ActionScene, func() {
		entry.UpdateC <- body.NeedsPaintTask{SceneStateUpdated: true}
	})

	defer func() {
		if v := recover(); v != nil {
			log.Printf("error in rendering scene %T for %s: %v\n", entry.ActionScene, entry.Destination.ToString(), v)
		}
	}()
	builder := external.NewDispatcherBuilder()
	entry.ActionScene.Init(builder)
	entry.ExternalCallbacks = builder
}

type NavigationStack struct {
	Current     *NavStackEntry
	Destination telegram.Destination
}

func (stack *NavigationStack) pop() {
	if stack.Current.Previous != nil {
		stack.Current.ActionScene.Dispose()
		stack.Current = stack.Current.Previous
	}
}

type NavImpl struct {
	NavC chan<- NavigationTask
}

func (n NavImpl) Push(scene ActionScene) {
	n.NavC <- NavigationTask{
		Action: ExtendAction,
		Target: scene,
	}
}

func (n NavImpl) Replace(scene ActionScene) {
	n.NavC <- NavigationTask{
		Action: ReplaceAction,
		Target: scene,
	}
}

func (n NavImpl) PopScene() {
	n.NavC <- NavigationTask{Action: PopAction}
}

func (n NavImpl) PopToMain() {
	n.NavC <- NavigationTask{Action: PopToMainAction}
}

func (stack *NavigationStack) popMain() {
	for stack.Current.Previous != nil {
		stack.Current.ActionScene.Dispose()
		stack.Current = stack.Current.Previous
	}
}

func (stack *NavigationStack) push(scene ActionScene, updateC chan<- body.NeedsPaintTask, isExtend bool, api *bot.Api) {
	if scene == nil {
		log.Printf("possible misuse of navigation api, called push with nil")
		return
	}
	var prev *NavStackEntry = stack.Current
	if !isExtend {
		prev = prev.Previous
	}
	stack.Current.dehydrate(api)
	stack.Current = NewNavStackEntry(scene, updateC, stack.Destination, prev)
	if isExtend {
		stack.Current.BodyResult = prev.BodyResult
		stack.Current.HeadResult = prev.HeadResult
	}

}

type NavigationKind int

const (
	ReplaceNavigationKind NavigationKind = iota
	ExtendNavigationKind  NavigationKind = iota
)

type NavigateAction int

const (
	ExtendAction    NavigateAction = iota
	ReplaceAction   NavigateAction = iota
	PopAction       NavigateAction = iota
	PopToMainAction NavigateAction = iota
)

type NavigationTask struct {
	Action NavigateAction
	Target ActionScene
}
