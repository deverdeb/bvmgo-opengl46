package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
	"ogl46/engine/graphic"
)

// Scene2dDrawer est un outil de dessin de la scène sur l'écran
type Scene2dDrawer struct {
	// application est l'application (permet de récupérer l'outil de rendu 2d pour le dessin à l'écran)
	application *engine.Application
	// scene correspond à la scène à rendre
	scene *Scene2d
	// 	viewSize est la dimension de la vue (pour initialiser le rendu 2d côté OpenGL) - par défaut, la taille de l'écran
	viewSize mgl32.Vec2
}

// Begin initialise le rendu de la vue 2d
func (drawer *Scene2dDrawer) Begin() {
	drawer.application.Renderer2d().Begin(drawer.viewSize.X(), drawer.viewSize.Y())
}

// End finalise le rendu de la vue 2d
func (drawer *Scene2dDrawer) End() {
	drawer.application.Renderer2d().End()
}

// DrawSprite draws a picture
func (drawer *Scene2dDrawer) DrawSprite(texture *graphic.Texture, targetPosition mgl32.Vec2) {
	// Texture zone
	sourcePosition := mgl32.Vec2{0., 0.}
	sourceDimension := mgl32.Vec2{float32(texture.Width()), float32(texture.Height())}
	// Screen / target zone
	screenTargetPosition := drawer.scene.PosSceneToScreen(targetPosition)
	screenTargetDimension := drawer.scene.DimSceneToScreen(sourceDimension)

	drawer.application.Renderer2d().DrawSpriteExWithRotateAndColor(
		texture, sourcePosition, sourceDimension, screenTargetPosition, screenTargetDimension,
		mgl32.Vec2{}, 0., graphic.White)
}

// DrawSpriteFromRect dessine une image à l'écran
func (drawer *Scene2dDrawer) DrawSpriteFromRect(texture *graphic.Texture, sourceRectangle graphic.Rectangle, targetRectangle graphic.Rectangle) {
	screenTarget := drawer.scene.RectSceneToScreen(targetRectangle)
	drawer.application.Renderer2d().DrawSpriteFromRect(texture, sourceRectangle, screenTarget)
}

// DrawSpriteFromRectWithColor teinte et dessine une image à l'écran
func (drawer *Scene2dDrawer) DrawSpriteFromRectWithColor(texture *graphic.Texture, sourceRectangle graphic.Rectangle, targetRectangle graphic.Rectangle, color graphic.Color) {
	screenTarget := drawer.scene.RectSceneToScreen(targetRectangle)
	drawer.application.Renderer2d().DrawSpriteFromRectWithColor(texture, sourceRectangle, screenTarget, color)
}

// DrawSpriteExWithRotateAndColor draws part of texture at position and applies a color and a rotation.
//
//	sourcePosition et sourceDimension définissent la zone à rendre au niveau de la texture.
//	targetPosition et targetDimension définissent la zone où effectuer le rendu à l'écran.
//	rotationCenter indique le centre de rotation à utiliser (ratio relatif à la taille de la zone à dessiner : 0. = à gauche ou en haut / 1. = à droite ou en bas)
//	rotationInDegree est la rotation à appliquer (angle en degrés)
//	color est la couleur à appliquer
func (drawer *Scene2dDrawer) DrawSpriteExWithRotateAndColor(texture *graphic.Texture, sourcePosition mgl32.Vec2, sourceDimension mgl32.Vec2,
	targetPosition mgl32.Vec2, targetDimension mgl32.Vec2, rotationCenter mgl32.Vec2, rotation graphic.Angle, color graphic.Color) {
	drawer.application.Renderer2d().DrawSpriteExWithRotateAndColor(texture, sourcePosition, sourceDimension,
		drawer.scene.PosSceneToScreen(targetPosition), drawer.scene.DimSceneToScreen(targetDimension),
		drawer.scene.DimSceneToScreen(rotationCenter), rotation, color)
}

// DrawText dessine un texte à la position indiquée
func (drawer *Scene2dDrawer) DrawText(font *graphic.Font, text string, position mgl32.Vec2, color graphic.Color) {
	screenPosition := drawer.scene.PosSceneToScreen(position)
	drawer.application.Renderer2d().DrawTextEx(font, text, screenPosition, drawer.scene.ratio, color)
}

// DrawTextInRect dessine un texte dans le rectangle indiqué en essayant d'utiliser toute la place
func (drawer *Scene2dDrawer) DrawTextInRect(font *graphic.Font, text string, targetRectangle graphic.Rectangle, color graphic.Color) {
	screenTarget := drawer.scene.RectSceneToScreen(targetRectangle)
	drawer.application.Renderer2d().DrawTextInRect(font, text, screenTarget, color)
}

// Taille de la scène
func (drawer *Scene2dDrawer) SceneDimension() mgl32.Vec2 {
	return drawer.scene.Dimension()
}

// RectSceneToScreen convertit un rectangle de la scène en rectangle à l'écran
func (drawer *Scene2dDrawer) RectSceneToScreen(rectangleScene graphic.Rectangle) graphic.Rectangle {
	return drawer.scene.RectSceneToScreen(rectangleScene)
}

// PosSceneToScreen convertit une position de la scène en position à l'écran
func (drawer *Scene2dDrawer) PosSceneToScreen(posScene mgl32.Vec2) mgl32.Vec2 {
	return drawer.scene.PosSceneToScreen(posScene)
}

// DimSceneToScreen convertit une dimension de la scène en dimension à l'écran
func (drawer *Scene2dDrawer) DimSceneToScreen(dimScene mgl32.Vec2) mgl32.Vec2 {
	return drawer.scene.DimSceneToScreen(dimScene)
}

// RectScreenToScene convertit un rectangle à l'écran en rectangle de la scène
func (drawer *Scene2dDrawer) RectScreenToScene(rectangleScene graphic.Rectangle) graphic.Rectangle {
	return drawer.scene.RectScreenToScene(rectangleScene)
}

// PosScreenToScene convertit une position à l'écran en position de la scène
func (drawer *Scene2dDrawer) PosScreenToScene(posScene mgl32.Vec2) mgl32.Vec2 {
	return drawer.scene.PosScreenToScene(posScene)
}

// DimScreenToScene convertit une dimension à l'écran en dimension de la scène
func (drawer *Scene2dDrawer) DimScreenToScene(dimScene mgl32.Vec2) mgl32.Vec2 {
	return drawer.scene.DimScreenToScene(dimScene)
}
