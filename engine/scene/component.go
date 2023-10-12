package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
	"ogl46/engine/graphic"
	"ogl46/engine/input"
)

type OnMouseMoveFunc func(application *engine.Application, scene *Scene2d, component Component, event *input.MouseMoveEvent, position mgl32.Vec2, timer *engine.Timer)
type OnMouseButtonFunc func(application *engine.Application, scene *Scene2d, component Component, event *input.MouseButtonEvent, position mgl32.Vec2, timer *engine.Timer)
type OnMouseScrollFunc func(application *engine.Application, scene *Scene2d, component Component, event *input.MouseScrollEvent, position mgl32.Vec2, timer *engine.Timer)
type OnDrawFunc func(application *engine.Application, scene *Scene2d, drawer *Scene2dDrawer, timer *engine.Timer)

type Component interface {
	// Position et dimension du composant
	Box() *graphic.Rectangle
	SetBox(graphic.Rectangle)

	// Le composant est-il visible ?
	Visible() bool
	SetVisible(bool)

	// Souris au dessus du composant ?
	MouseOver() bool
	SetMouseOver(bool)

	// Traitement
	Execute(application *engine.Application, scene *Scene2d, timer *engine.Timer)

	// Dessin
	Draw(application *engine.Application, scene *Scene2d, drawer *Scene2dDrawer, timer *engine.Timer)

	// Trouver le composant sous la position indiqué (retourne null si rien)
	GetComponentUnder(position mgl32.Vec2) Component

	// Traitement des évènements sur la souris
	OnMouseMove() OnMouseMoveFunc
	OnMouseButton() OnMouseButtonFunc
	OnMouseScroll() OnMouseScrollFunc
}

type GenericComponent struct {
	// Position et dimension du composant
	box graphic.Rectangle
	// Le composant est-il visible
	visible bool
	// Souris au-dessus du composant ?
	mouseOver bool

	OnMouseMoveFunc   OnMouseMoveFunc
	OnMouseButtonFunc OnMouseButtonFunc
	OnMouseScrollFunc OnMouseScrollFunc
}

func BuildGenericComponent() GenericComponent {
	return GenericComponent{
		box:               graphic.Rectangle{0, 0, 100, 100},
		visible:           true,
		mouseOver:         false,
		OnMouseMoveFunc:   nil,
		OnMouseButtonFunc: nil,
		OnMouseScrollFunc: nil,
	}
}

// Box retourne la position et dimension du composant
func (comp *GenericComponent) Box() *graphic.Rectangle {
	return &comp.box
}

// SetBox modifie la position et dimension du composant
func (comp *GenericComponent) SetBox(box graphic.Rectangle) {
	comp.box = box
}

func (comp *GenericComponent) Visible() bool {
	return comp.visible
}

func (comp *GenericComponent) SetVisible(visible bool) {
	comp.visible = visible
}

// MouseOver indique si la souris au-dessus du composant
func (comp *GenericComponent) MouseOver() bool {
	return comp.mouseOver
}

// SetMouseOver permet d'indiquer si la souris au-dessus du composant
func (comp *GenericComponent) SetMouseOver(over bool) {
	comp.mouseOver = over
}

// Est-ce le composant sous la position indiquée ?
func (comp *GenericComponent) IsComponentUnder(position mgl32.Vec2) bool {
	return comp.box.In(position)
}

// Traitement des évènements sur la souris
func (comp *GenericComponent) OnMouseMove() OnMouseMoveFunc {
	return comp.OnMouseMoveFunc
}
func (comp *GenericComponent) OnMouseButton() OnMouseButtonFunc {
	return comp.OnMouseButtonFunc
}
func (comp *GenericComponent) OnMouseScroll() OnMouseScrollFunc {
	return comp.OnMouseScrollFunc
}
