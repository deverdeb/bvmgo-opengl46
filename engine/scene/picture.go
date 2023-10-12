package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
	"ogl46/engine/graphic"
)

type PictureComponent struct {
	GenericComponent

	texture *graphic.Texture
}

func (picture *PictureComponent) GetComponentUnder(position mgl32.Vec2) Component {
	if picture.IsComponentUnder(position) {
		return picture
	}
	return nil
}

func BuildPictureComponent(texture *graphic.Texture) *PictureComponent {
	return &PictureComponent{
		GenericComponent: GenericComponent{
			visible: true,
			box:     graphic.Rectangle{0, 0, float32(texture.Width()), float32(texture.Height())},
		},
		texture: texture,
	}
}

// Traitement
func (picture *PictureComponent) Execute(application *engine.Application, scene *Scene2d, timer *engine.Timer) {
	// Rien ici
}

// Dessin
func (picture *PictureComponent) Draw(application *engine.Application, scene *Scene2d, drawer *Scene2dDrawer, timer *engine.Timer) {
	source := picture.texture.Rectangle()
	target := picture.box
	drawer.DrawSpriteFromRect(picture.texture, source, target)
}
