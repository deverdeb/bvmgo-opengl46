package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// KeyAction représente une action sur une touche.
type KeyAction uint8

const (
	// KEY_PRESS indique que la touche a été pressée
	KEY_PRESS KeyAction = iota
	// KEY_RELEASE indique que la touche a été relâchée
	KEY_RELEASE
	// KEY_REPEAT indique que la touche est toujours pressée
	KEY_REPEAT
)

// KeyboardEvent représente un évènement provenant du clavier.
type KeyboardEvent struct {
	// Keyboard est le périphérique concerné par l'évènement
	keyboard *Keyboard
	// Key indique la touche concernée par l'évènement
	key glfw.Key
	// ScanCode indique le code de la touche concernée par l'évènement
	scanCode int
	// Action sur la touche
	action KeyAction
}

func (event *KeyboardEvent) Keyboard() *Keyboard {
	return event.keyboard
}

func (event *KeyboardEvent) Source() EventSource {
	return KEYBOARD
}

func (event *KeyboardEvent) Key() glfw.Key {
	return event.key
}

func (event *KeyboardEvent) ScanCode() int {
	return event.scanCode
}

func (event *KeyboardEvent) Action() KeyAction {
	return event.action
}

func convertGlfwActionToKeyAction(action glfw.Action) KeyAction {
	if action == glfw.Press {
		return KEY_PRESS
	} else if action == glfw.Release {
		return KEY_RELEASE
	} else {
		return KEY_REPEAT
	}
}
