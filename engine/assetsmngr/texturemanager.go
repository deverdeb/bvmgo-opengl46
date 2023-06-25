package assetsmngr

import (
	"fmt"
	"ogl46/engine/graphic"
)

// TextureManager permet de gérer le chargement et la libération des textures
type TextureManager struct {
	// Manager est une instance d'asset manager pour gérer les ressources
	Manager[graphic.Texture]
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		Manager: CreateManager[graphic.Texture](),
	}
}

func (textureManager *TextureManager) RegisterTextureFromFile(name string, filename string) {
	textureManager.Manager.Register(name,
		func() (*graphic.Texture, error) {
			texture, err := graphic.LoadTextureFromFile(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' texture from file '%s'.\n - %w", name, filename, err)
			}
			return texture, nil
		},
		func(font *graphic.Texture) {
			if font != nil {
				font.Release()
			}
		})
}

func (textureManager *TextureManager) RegisterTextureFromBytes(name string, content []byte, format graphic.ImageFormat) {
	textureManager.Manager.Register(name,
		func() (*graphic.Texture, error) {
			texture, err := graphic.LoadTextureFromBytes(content, format)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' texture from byte array.\n - %w", name, err)
			}
			return texture, nil
		},
		func(font *graphic.Texture) {
			if font != nil {
				font.Release()
			}
		})
}
