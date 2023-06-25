package engine

import (
	"ogl46/engine/input"
)

type Stage interface {
	Initialize(app *Application) error

	Release(app *Application)

	Display(app *Application, timer *Timer)

	Execute(app *Application, timer *Timer)

	ProcessEvent(app *Application, event input.Event, timer *Timer)
}
