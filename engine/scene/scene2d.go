package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
	"ogl46/engine/graphic"
	"ogl46/engine/input"
	"slices"
)

type Scene2d struct {
	// size est la dimension interne de la scène
	internalSize mgl32.Vec2
	// position est la position de la scène à l'écran (coordonnées "écran")
	position mgl32.Vec2
	// ratio et le ratio dimension écran / dimension scène
	ratio mgl32.Vec2

	// Liste des composants de la scène
	components []Component

	componentUnderMouse Component

	// Sur dessin
	OnDraw OnDrawFunc
}

// BuildScene2d permet d'obtenir une scène 2D
func BuildScene2d(width, height float32) *Scene2d {
	return &Scene2d{
		internalSize: mgl32.Vec2{width, height},
		position:     mgl32.Vec2{0, 0},
		ratio:        mgl32.Vec2{1, 1},
	}
}

// DrawerBuilder permet d'obtenir une structure de mise en place d'un "scene2d" pour la scène
func (scene2d *Scene2d) DrawerBuilder(application *engine.Application) *Scene2dDrawerBuilder {
	screenWidth, screenHeight := application.Size()
	return &Scene2dDrawerBuilder{
		application:      application,
		viewSize:         mgl32.Vec2{float32(screenWidth), float32(screenHeight)},
		screenTargetZone: graphic.Rectangle{0, 0, -1, -1},
		scene:            scene2d,
		keepRatio:        true,
	}
}

// Traitement
func (scene2d *Scene2d) Execute(application *engine.Application, timer *engine.Timer) {
	for _, component := range scene2d.components {
		component.Execute(application, scene2d, timer)
	}
}

// Dessin
func (scene2d *Scene2d) Draw(application *engine.Application, timer *engine.Timer) {
	drawer := scene2d.DrawerBuilder(application).Build()
	drawer.Begin()
	defer drawer.End()
	for _, component := range scene2d.components {
		if component.Visible() {
			component.Draw(application, scene2d, drawer, timer)
		}
	}
	if scene2d.OnDraw != nil {
		scene2d.OnDraw(application, scene2d, drawer, timer)
	}
}

// Traiter les évènements
func (scene2d *Scene2d) ProcessEvent(application *engine.Application, event input.Event, timer *engine.Timer) {
	switch event.Source() {
	case input.MOUSE_BUTTON:
		mouseEvent := event.(*input.MouseButtonEvent)
		mousePosition := scene2d.PosScreenToScene(mgl32.Vec2{float32(mouseEvent.Mouse().X()), float32(mouseEvent.Mouse().Y())})
		component := scene2d.GetComponentUnder(mousePosition)
		if component != nil && component.OnMouseButton() != nil {
			component.OnMouseButton()(application, scene2d, component, mouseEvent, mousePosition, timer)
		}
	case input.MOUSE_MOVE:
		mouseEvent := event.(*input.MouseMoveEvent)
		mousePosition := scene2d.PosScreenToScene(mgl32.Vec2{float32(mouseEvent.Mouse().X()), float32(mouseEvent.Mouse().Y())})
		component := scene2d.GetComponentUnder(mousePosition)
		if scene2d.componentUnderMouse != component {
			if scene2d.componentUnderMouse != nil {
				scene2d.componentUnderMouse.SetMouseOver(false)
			}
			scene2d.componentUnderMouse = component
			if scene2d.componentUnderMouse != nil {
				scene2d.componentUnderMouse.SetMouseOver(true)
			}
		}
		if component != nil && component.OnMouseMove() != nil {
			component.OnMouseMove()(application, scene2d, component, mouseEvent, mousePosition, timer)
		}
	case input.MOUSE_SCROLL:
		mouseEvent := event.(*input.MouseScrollEvent)
		mousePosition := scene2d.PosScreenToScene(mgl32.Vec2{float32(mouseEvent.Mouse().X()), float32(mouseEvent.Mouse().Y())})
		component := scene2d.GetComponentUnder(mousePosition)
		if component != nil && component.OnMouseScroll() != nil {
			component.OnMouseScroll()(application, scene2d, component, mouseEvent, mousePosition, timer)
		}
	}
}

// Taille de la scène
func (scene2d *Scene2d) Dimension() mgl32.Vec2 {
	return scene2d.internalSize
}

// Trouver le composant sous la position indiqué (retourne null si rien)
func (scene2d *Scene2d) GetComponentUnder(position mgl32.Vec2) Component {
	for idx := len(scene2d.components) - 1; idx >= 0; idx-- {
		component := scene2d.components[idx]
		if componentUnder := component.GetComponentUnder(position); componentUnder != nil && componentUnder.Visible() {
			return componentUnder
		}
	}
	return nil
}

// RemoveComponent supprime un composant
func (scene2d *Scene2d) RemoveComponent(component Component) {
	// Enlever le composant s'il est déjà présent
	slices.DeleteFunc(scene2d.components, func(comp Component) bool { return comp == component })
}

// AddComponent Ajoute un composant
func (scene2d *Scene2d) AddComponent(component Component) {
	scene2d.RemoveComponent(component)
	// Ajouter le composant à la fin des composants
	scene2d.components = append(scene2d.components, component)
}

// RectSceneToScreen convertit un rectangle de la scène en rectangle à l'écran
func (scene2d *Scene2d) RectSceneToScreen(rectangleScene graphic.Rectangle) graphic.Rectangle {
	return graphic.Rectangle{
		rectangleScene.X()*scene2d.ratio.X() + scene2d.position.X(), rectangleScene.Y()*scene2d.ratio.Y() + scene2d.position.Y(),
		rectangleScene.Width() * scene2d.ratio.X(), rectangleScene.Height() * scene2d.ratio.Y(),
	}
}

// PosSceneToScreen convertit une position de la scène en position à l'écran
func (scene2d *Scene2d) PosSceneToScreen(posScene mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		posScene.X()*scene2d.ratio.X() + scene2d.position.X(), posScene.Y()*scene2d.ratio.Y() + scene2d.position.Y(),
	}
}

// DimSceneToScreen convertit une dimension de la scène en dimension à l'écran
func (scene2d *Scene2d) DimSceneToScreen(dimScene mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		dimScene.X() * scene2d.ratio.X(), dimScene.Y() * scene2d.ratio.Y(),
	}
}

// RectScreenToScene convertit un rectangle à l'écran en rectangle de la scène
func (scene2d *Scene2d) RectScreenToScene(rectangleScene graphic.Rectangle) graphic.Rectangle {
	return graphic.Rectangle{
		(rectangleScene.X() - scene2d.position.X()) / scene2d.ratio.X(), (rectangleScene.Y() - scene2d.position.Y()) / scene2d.ratio.Y(),
		rectangleScene.Width() / scene2d.ratio.X(), rectangleScene.Height() / scene2d.ratio.Y(),
	}
}

// PosScreenToScene convertit une position à l'écran en position de la scène
func (scene2d *Scene2d) PosScreenToScene(posScene mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		(posScene.X() - scene2d.position.X()) / scene2d.ratio.X(), (posScene.Y() - scene2d.position.Y()) / scene2d.ratio.Y(),
	}
}

// DimScreenToScene convertit une dimension à l'écran en dimension de la scène
func (scene2d *Scene2d) DimScreenToScene(dimScene mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		dimScene.Y() / scene2d.ratio.X(), dimScene.Y() / scene2d.ratio.Y(),
	}
}
