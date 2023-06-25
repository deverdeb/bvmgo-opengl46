package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// MouseButtonAction représente une action sur un bouton de la souris.
type MouseButtonAction uint8

const (
	// BUTTON_PRESS indique que le bouton a été pressé
	BUTTON_PRESS MouseButtonAction = iota
	// BUTTON_RELEASE indique que le bouton a été relâché
	BUTTON_RELEASE
	// BUTTON_REPEAT indique que le bouton est toujours pressé
	BUTTON_REPEAT
)

// MouseEvent représente un évènement provenant de la souris.
type mouseEvent struct {
	// mouse est le périphérique concerné par l'évènement
	mouse *Mouse
}

// Mouse retourne le périphérique concerné par l'évènement
func (event *mouseEvent) Mouse() *Mouse {
	return event.mouse
}

// MouseMoveEvent représente un évènement provenant d'un déplacement de la souris.
type MouseMoveEvent struct {
	mouseEvent
	// pos est la position du curseur de la souris
	pos mgl64.Vec2
	// shift est le déplacement du curseur de la souris
	shift mgl64.Vec2
}

// Mouse retourne le périphérique concerné par l'évènement
func (event *MouseMoveEvent) Mouse() *Mouse {
	return event.mouse
}

// X retourne la position X de la souris
func (event *MouseMoveEvent) X() float64 {
	return event.pos.X()
}

// Y retourne la position Y de la souris
func (event *MouseMoveEvent) Y() float64 {
	return event.pos.Y()
}

// Pos retourne la position de la souris
func (event *MouseMoveEvent) Pos() (x, y float64) {
	return event.pos.X(), event.pos.Y()
}

// Pos2d retourne la position de la souris
func (event *MouseMoveEvent) Pos2d() mgl64.Vec2 {
	return event.pos
}

// HShift retourne le déplacement horizontal du curseur de la souris
func (event *MouseMoveEvent) HShift() float64 {
	return event.shift.X()
}

// VShift retourne le déplacement vertical du curseur de la souris
func (event *MouseMoveEvent) VShift() float64 {
	return event.shift.Y()
}

// Shift retourne le déplacement du curseur de la souris
func (event *MouseMoveEvent) Shift() (x, y float64) {
	return event.shift.X(), event.shift.Y()
}

// ShiftPos2d retourne le déplacement du curseur de la souris
func (event *MouseMoveEvent) ShiftPos2d() mgl64.Vec2 {
	return event.shift
}

func (event *MouseMoveEvent) Source() EventSource {
	return MOUSE_MOVE
}

// MouseButtonEvent représente un évènement provenant d'un bouton de la souris.
type MouseButtonEvent struct {
	mouseEvent
	// Button indique le bouton concerné par l'évènement
	button glfw.MouseButton
	// Action sur le bouton
	action MouseButtonAction
}

func (event *MouseButtonEvent) Source() EventSource {
	return MOUSE_BUTTON
}

func (event *MouseButtonEvent) Button() glfw.MouseButton {
	return event.button
}

func (event *MouseButtonEvent) Action() MouseButtonAction {
	return event.action
}

func convertGlfwActionToMouseButtonAction(action glfw.Action) MouseButtonAction {
	if action == glfw.Press {
		return BUTTON_PRESS
	} else if action == glfw.Release {
		return BUTTON_RELEASE
	} else {
		return BUTTON_REPEAT
	}
}

// MouseScrollEvent représente un évènement provenant de la molette de la souris.
type MouseScrollEvent struct {
	mouseEvent
	// hOffset représente la position de la molette de la souris
	offset mgl64.Vec2
	// shift est le déplacement de la molette de la souris
	shift mgl64.Vec2
}

func (event *MouseScrollEvent) Source() EventSource {
	return MOUSE_SCROLL
}

// HOffset retourne la position de la molette horizontale
func (event *MouseScrollEvent) HOffset() float64 {
	return event.offset.X()
}

// VOffset retourne la position de la molette verticale
func (event *MouseScrollEvent) VOffset() float64 {
	return event.offset.Y()
}

// Offset retourne la position de la molette de la souris
func (event *MouseScrollEvent) Offset() (x, y float64) {
	return event.offset.X(), event.offset.Y()
}

// OffsetPos2d retourne la position de la molette de la souris
func (event *MouseScrollEvent) OffsetPos2d() mgl64.Vec2 {
	return event.offset
}

// HShift retourne le déplacement de la molette horizontale
func (event *MouseScrollEvent) HShift() float64 {
	return event.shift.X()
}

// VShift retourne le déplacement de la molette verticale
func (event *MouseScrollEvent) VShift() float64 {
	return event.shift.Y()
}

// Shift retourne le déplacement de la molette de la souris
func (event *MouseScrollEvent) Shift() (x, y float64) {
	return event.shift.X(), event.shift.Y()
}

// ShiftPos2d retourne le déplacement de la molette de la souris
func (event *MouseScrollEvent) ShiftPos2d() mgl64.Vec2 {
	return event.shift
}
