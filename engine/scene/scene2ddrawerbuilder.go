package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
	"ogl46/engine/graphic"
)

// Scene2dDrawerBuilder est une classe de construction de l'outil de rendu de la scène à l'écran
type Scene2dDrawerBuilder struct {
	// application est l'application
	application *engine.Application
	// scene correspond à la scène à rendre
	scene *Scene2d
	// 	viewSize est la dimension de la vue (pour initialiser le rendu 2d côté OpenGL) - par défaut, la taille de l'écran
	viewSize mgl32.Vec2
	// screenTargetZone est la zone de l'écran sur laquelle il faut rendre la scène - tout l'écran par défaut
	screenTargetZone graphic.Rectangle
	// keepRatio indique s'il faut conserver le ratio Largeur / hauteur - vrai par défaut
	keepRatio bool
}

func (builder *Scene2dDrawerBuilder) KeepRatio(keep bool) *Scene2dDrawerBuilder {
	builder.keepRatio = keep
	return builder
}

func (builder *Scene2dDrawerBuilder) ViewSize(size mgl32.Vec2) *Scene2dDrawerBuilder {
	builder.viewSize = size
	return builder
}

func (builder *Scene2dDrawerBuilder) ScreenZone(zone graphic.Rectangle) *Scene2dDrawerBuilder {
	builder.screenTargetZone = zone
	return builder
}

func (builder *Scene2dDrawerBuilder) Build() *Scene2dDrawer {
	if builder.screenTargetZone.Width() == -1 {
		builder.screenTargetZone.SetWidth(builder.viewSize.X())
	}
	if builder.screenTargetZone.Height() == -1 {
		builder.screenTargetZone.SetHeight(builder.viewSize.Y())
	}
	ratioX := builder.screenTargetZone.Width() / builder.scene.internalSize.X()
	ratioY := builder.screenTargetZone.Height() / builder.scene.internalSize.Y()
	if builder.keepRatio {
		ratioX = min(ratioX, ratioY)
	}
	renderWidth := builder.scene.internalSize.X() * ratioX
	renderHeight := builder.scene.internalSize.Y() * ratioY
	renderPosX := (builder.screenTargetZone.Width()-renderWidth)/2 + builder.screenTargetZone.X()
	renderPosY := (builder.screenTargetZone.Height()-renderHeight)/2 + builder.screenTargetZone.Y()

	builder.scene.position = mgl32.Vec2{renderPosX, renderPosY}
	builder.scene.ratio = mgl32.Vec2{ratioX, ratioY}

	return &Scene2dDrawer{
		application: builder.application,
		scene:       builder.scene,
		viewSize:    builder.viewSize,
	}
}
