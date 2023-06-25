package engine

import (
	"ogl46/engine/input"
)

type EmptyStage struct {
}

func (stage *EmptyStage) Initialize(app *Application) error {
	// FIXME à implémenter
	return nil
}

func (stage *EmptyStage) Release(app *Application) {
	// FIXME à implémenter
}

func (stage *EmptyStage) Display(app *Application, timer *Timer) {
	// FIXME à implémenter
}

func (stage *EmptyStage) Execute(app *Application, timer *Timer) {
	// FIXME à implémenter
}

func (stage *EmptyStage) ProcessEvent(app *Application, event input.Event, timer *Timer) {
	// FIXME à implémenter
	switch event.Source() {
	case input.WINDOW:
		// Traitement minimum : sur "bouton fermeture" de la fenêtre, arrêter l'application
		if event.(*input.WindowEvent).Action() == input.WINDOW_CLOSED {
			app.Stop()
		}
	}
}
