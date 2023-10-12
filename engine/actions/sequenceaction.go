package actions

import (
	"ogl46/engine"
)

type SequenceAction struct {
	actions      []Action
	currentIndex int
}

func BuildSequenceAction() *SequenceAction {
	return &SequenceAction{
		actions:      make([]Action, 0),
		currentIndex: 0,
	}
}

func (action *SequenceAction) Add(newAction Action) {
	if newAction != nil && !newAction.IsEnd() {
		action.actions = append(action.actions, newAction)
	}
}

func (action *SequenceAction) Clear() {
	action.actions = make([]Action, 0)
}

func (action *SequenceAction) Execute(timer *engine.Timer) {
	if action.IsEnd() {
		return
	}
	currentAction := action.actions[action.currentIndex]
	if currentAction == nil || currentAction.IsEnd() {
		action.currentIndex++
		action.Execute(timer)
		return
	} else {
		currentAction.Execute(timer)
	}
}

func (action *SequenceAction) IsEnd() bool {
	return action.currentIndex >= len(action.actions)
}
