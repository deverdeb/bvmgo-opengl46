package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// MaxMouseButtonCode est le code de bouton maximum supporté par le moteur
const MaxMouseButtonCode = 8

// MouseButtonCode représente un bouton de la souris
type MouseButtonCode glfw.MouseButton

// ButtonStatus représente l'état d'un bouton.
type ButtonStatus uint8

const (
	// BUTTON_RELEASED indique que le bouton est relâché
	BUTTON_RELEASED ButtonStatus = iota
	// BUTTON_PRESSED indique que le bouton est pressé
	BUTTON_PRESSED
)

// MouseMoveMode indique le mode de déplacement de la souris
type MouseMoveMode uint8

const (
	// MOUSE_MOVE_CURSOR indique que le mode de déplacement est "curseur"
	MOUSE_MOVE_CURSOR MouseMoveMode = iota
	// MOUSE_MOVE_RAW indique que le mode de déplacement est "brut"
	MOUSE_MOVE_RAW
)

// Mouse représente une souris
type Mouse struct {
	// glfwWindow est la fenêtre GLFW.
	glfwWindow *glfw.Window
	// buttons est l'état des doubons de la souris
	buttons [MaxMouseButtonCode]ButtonStatus
	// pos est la position du curseur de la souris
	pos mgl64.Vec2
	// scrollWheels est la position des molettes de la souris
	scrollWheels mgl64.Vec2
	// moveMode est le mode de déplacement de la souris
	moveMode MouseMoveMode
	// cursorHidden indique si le curseur est masqué
	cursorHidden bool
	// cursorLocked indique si le curseur est "bloqué" au centre de la fenêtre
	cursorLocked bool
}

// X retourne la position X de la souris
func (mouse *Mouse) X() float64 {
	return mouse.pos.X()
}

// Y retourne la position Y de la souris
func (mouse *Mouse) Y() float64 {
	return mouse.pos.Y()
}

// Pos retourne la position de la souris
func (mouse *Mouse) Pos() (x, y float64) {
	return mouse.pos.X(), mouse.pos.Y()
}

// SetPos déplace la position de la souris
func (mouse *Mouse) SetPos(x, y float64) {
	mouse.pos = mgl64.Vec2{x, y}
	mouse.glfwWindow.SetCursorPos(x, y)
}

// Pos2d retourne la position de la souris
func (mouse *Mouse) Pos2d() mgl64.Vec2 {
	return mouse.pos
}

// SetPos2d déplace la position de la souris
func (mouse *Mouse) SetPos2d(pos mgl64.Vec2) {
	mouse.pos = pos
	mouse.glfwWindow.SetCursorPos(pos.X(), pos.Y())
}

// HScroll retourne la position de la molette horizontale de la souris
func (mouse *Mouse) HScroll() float64 {
	return mouse.scrollWheels.X()
}

// VScroll retourne la position de la molette verticale de la souris
func (mouse *Mouse) VScroll() float64 {
	return mouse.scrollWheels.Y()
}

// ScrollWheels retourne la position des molettes de la souris
func (mouse *Mouse) ScrollWheels() (x, y float64) {
	return mouse.scrollWheels.X(), mouse.scrollWheels.Y()
}

// ScrollWheelsPos2d retourne la position des molettes de la souris
func (mouse *Mouse) ScrollWheelsPos2d() mgl64.Vec2 {
	return mouse.scrollWheels
}

// ButtonStatus retourne l'état d'un bouton
func (mouse *Mouse) ButtonStatus(mouseButtonCode MouseButtonCode) ButtonStatus {
	var code = int(mouseButtonCode)
	if code < 0 || code >= MaxMouseButtonCode {
		return BUTTON_RELEASED
	} else {
		return mouse.buttons[code]
	}
}

func (mouse *Mouse) MoveMode() MouseMoveMode {
	return mouse.moveMode
}

func (mouse *Mouse) SetMoveMode(moveMode MouseMoveMode) {
	mouse.moveMode = moveMode
	if moveMode == MOUSE_MOVE_RAW {
		if glfw.RawMouseMotionSupported() {
			mouse.glfwWindow.SetInputMode(glfw.RawMouseMotion, glfw.True)
		}
	} else if moveMode == MOUSE_MOVE_CURSOR {
		mouse.glfwWindow.SetInputMode(glfw.RawMouseMotion, glfw.False)
	}
}

func (mouse *Mouse) CursorHidden() bool {
	return mouse.cursorHidden
}

func (mouse *Mouse) SetCursorHidden(hidden bool) {
	mouse.cursorHidden = hidden
	if hidden {
		mouse.glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	} else {
		mouse.glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

func (mouse *Mouse) CursorLocked() bool {
	return mouse.cursorLocked
}

func (mouse *Mouse) SetCursorLocked(locked bool) {
	mouse.cursorLocked = locked
	if locked {
		x, y := mouse.glfwWindow.GetSize()
		mouse.pos = mgl64.Vec2{float64(x) / 2., float64(y) / 2.}
		mouse.glfwWindow.SetCursorPos(mouse.pos.X(), mouse.pos.Y())
	}
}

func (mouse *Mouse) updateWithMouseButtonEvent(mouseButtonEvent *MouseButtonEvent) {
	var code = int(mouseButtonEvent.Button())
	if code >= 0 && code < MaxKeyCode {
		if mouseButtonEvent.Action() == BUTTON_PRESS {
			mouse.buttons[code] = BUTTON_PRESSED
		} else if mouseButtonEvent.Action() == BUTTON_RELEASE {
			mouse.buttons[code] = BUTTON_RELEASED
		}
	}
}

func (mouse *Mouse) updateWithMouseMoveEvent(mouseMoveEvent *MouseMoveEvent) {
	if mouse.cursorLocked {
		mouse.glfwWindow.SetCursorPos(mouse.Pos())
	} else {
		mouse.pos = mgl64.Vec2{mouseMoveEvent.X(), mouseMoveEvent.Y()}
	}
}

func (mouse *Mouse) updateWithMouseScrollEvent(mouseScrollEvent *MouseScrollEvent) {
	mouse.scrollWheels = mgl64.Vec2{mouseScrollEvent.HOffset(), mouseScrollEvent.VOffset()}
}
