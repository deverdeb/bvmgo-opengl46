package graphic

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"log/slog"
)

type Font struct {
	// Texture
	texture *Texture
	// Dimension de la texture : largeur
	textureWidth int32
	// Dimension de la texture : hauteur
	textureHeight int32
	// Dimension des caractères : largeur
	charWidth int32
	// Dimension des caractères : hauteur
	charHeight int32
}

func (font *Font) Release() {
	if font.texture != nil {
		slog.Debug("font texture destruction")
		font.texture.Release()
		font.texture = nil
	}
}

func LoadFontFromBitmapFile(filename string, charWidth int32, charHeight int32) (*Font, error) {
	slog.Debug("font texture creation from file", "filename", filename)
	// Charger l'image
	texture, err := LoadTextureFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create bitmap font from texture file '%s'\n - %w", filename, err)
	}

	// On crée notre objet texture qui sera retournée
	asset := &Font{
		texture:    texture,
		charWidth:  charWidth,
		charHeight: charHeight,
	}
	return asset, nil
}

func LoadFontFromBytes(content []byte, format ImageFormat, charWidth int32, charHeight int32) (*Font, error) {
	slog.Debug("font texture creation from byte array")
	// Charger l'image
	texture, err := LoadTextureFromBytes(content, format)
	if err != nil {
		return nil, fmt.Errorf("failed to create bitmap font from byte array\n - %w", err)
	}

	// On crée notre objet texture qui sera retournée
	asset := &Font{
		texture:    texture,
		charWidth:  charWidth,
		charHeight: charHeight,
	}
	return asset, nil
}

func (font *Font) Texture() *Texture {
	return font.texture
}

func (font *Font) CharacterWidth() int32 {
	return font.charWidth
}

func (font *Font) CharacterHeight() int32 {
	return font.charHeight
}

func (font *Font) CharacterPosition(character int32) mgl32.Vec2 {
	idxChar := character - 32
	return mgl32.Vec2{
		float32((idxChar % font.charWidth) * font.charWidth),
		float32((idxChar / font.charWidth) * font.charHeight),
	}
}

func (font *Font) CharacterRectangle(character int32) Rectangle {
	idxChar := character - 32
	return BuildRectangle(
		float32((idxChar%font.charWidth)*font.charWidth),
		float32((idxChar/font.charWidth)*font.charHeight),
		float32(font.charWidth),
		float32(font.charHeight),
	)
}
