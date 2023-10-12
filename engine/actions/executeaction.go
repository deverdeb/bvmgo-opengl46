package actions

import (
	"ogl46/engine"
)

type ExecuteAction struct {
	function func()
	end      bool
}

func BuildExecuteAction(function func()) Action {
	return &ExecuteAction{
		function: function,
		end:      false,
	}
}

func (action *ExecuteAction) Execute(timer *engine.Timer) {
	if action.function != nil {
		action.function()
	}
	action.end = true
}

func (action *ExecuteAction) IsEnd() bool {
	return action.end
}
