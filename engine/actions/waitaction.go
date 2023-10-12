package actions

import (
	"ogl46/engine"
)

type WaitAction struct {
	waitTimeInMs int64
}

func BuildWaitAction(waitTimeInMs int64) Action {
	return &WaitAction{
		waitTimeInMs: waitTimeInMs,
	}
}

func (action *WaitAction) Execute(timer *engine.Timer) {
	if !action.IsEnd() {
		action.waitTimeInMs -= timer.ElapsedTime()
	}
}

func (action *WaitAction) IsEnd() bool {
	return action.waitTimeInMs <= 0
}
