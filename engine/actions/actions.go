package actions

import (
	"ogl46/engine"
	"slices"
)

type Action interface {
	Execute(timer *engine.Timer)
	IsEnd() bool
}

type ActionEngine struct {
	actions []Action
}

func (engine *ActionEngine) Actions() []Action {
	if engine.actions == nil {
		engine.actions = make([]Action, 0)
	}
	return engine.actions
}

func (engine *ActionEngine) Add(action Action) {
	engine.actions = append(engine.Actions(), action)
}

func (engine *ActionEngine) Clear() {
	engine.actions = make([]Action, 0)
}

func (engine *ActionEngine) Execute(timer *engine.Timer) {
	// Traiter les actions
	toRemove := make([]Action, 0)
	for _, action := range engine.Actions() {
		action.Execute(timer)
		if action.IsEnd() { // Stocker les actions finies
			toRemove = append(toRemove, action)
		}
	}
	// Enlever les actions finies
	engine.actions = slices.DeleteFunc(engine.actions, func(elt Action) bool {
		return slices.Contains(toRemove, elt)
	})
}
