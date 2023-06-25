package input

// WindowAction représente une action sur une fenêtre.
type WindowAction uint8

const (
	// WINDOW_RESIZED indique que la fenêtre est redimensionnée (ou minimized ou maximized)
	WINDOW_RESIZED WindowAction = iota
	// WINDOW_MAXIMIZED indique que la fenêtre est agrandie
	WINDOW_MAXIMIZED
	// WINDOW_ICONIFIED indique que la fenêtre est iconifiée
	WINDOW_ICONIFIED
	// WINDOW_RESTORED indique que la fenêtre est restaurée (suite à un agrandissement ou à une iconification)
	WINDOW_RESTORED
	// WINDOW_CLOSED indique que la fenêtre est fermée
	WINDOW_CLOSED
	// WINDOW_FOCUS indique que la fenêtre gagne le focus
	WINDOW_FOCUS
	// WINDOW_UNFOCUS indique que la fenêtre perd le focus
	WINDOW_UNFOCUS
)

// WindowEvent représente un évènement provenant de la fenêtre.
type WindowEvent struct {
	// window est la fenêtre concernée par l'évènement
	window *Window
	// width indique la largeur de la fenêtre
	width int
	// heigth indique la hauteur de la fenêtre
	heigth int
	// Action sur la fenêtre
	action WindowAction
}

func (event *WindowEvent) Source() EventSource {
	return WINDOW
}

func (event *WindowEvent) Window() *Window {
	return event.window
}

func (event *WindowEvent) Width() int {
	return event.width
}

func (event *WindowEvent) Heigth() int {
	return event.heigth
}

func (event *WindowEvent) Action() WindowAction {
	return event.action
}
