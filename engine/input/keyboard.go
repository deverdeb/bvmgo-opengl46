package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// MaxKeyCode est le code de touche maximum supporté par le moteur
const MaxKeyCode = 350

// KeyCode représente une touche du clavier
type KeyCode glfw.Key

// KeyStatus représente l'état d'une touche.
type KeyStatus uint8

const (
	// KEY_RELEASED indique que la touche est relâchée
	KEY_RELEASED KeyStatus = iota
	// KEY_PRESSED indique que la touche est pressée
	KEY_PRESSED
)

// Keyboard représente un clavier
type Keyboard struct {
	keys [MaxKeyCode]KeyStatus
}

// KeyStatus retourne l'état d'une touche
func (keyboard *Keyboard) KeyStatus(keyCode KeyCode) KeyStatus {
	var code = int(keyCode)
	if code < 0 || code >= MaxKeyCode {
		return KEY_RELEASED
	} else {
		return keyboard.keys[code]
	}
}

// updateKeyboardWithKeyEvent permet de mettre à jour le clavier à partir d'un évènement sur une touche
func (keyboard *Keyboard) updateKeyboardWithKeyEvent(keyEvent *KeyboardEvent) {
	var code = int(keyEvent.Key())
	if code >= 0 && code < MaxKeyCode {
		if keyEvent.Action() == KEY_PRESS {
			keyboard.keys[code] = KEY_PRESSED
		} else if keyEvent.Action() == KEY_RELEASE {
			keyboard.keys[code] = KEY_RELEASED
		}
	}
}
