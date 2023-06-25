package assetsmngr

import (
	"fmt"
	"ogl46/engine/graphic"
)

// FontManager permet de gérer le chargement et la libération des polices de caractères
type FontManager struct {
	// Manager est une instance d'asset manager pour gérer les ressources
	Manager[graphic.Font]
}

func NewFontManager() *FontManager {
	return &FontManager{
		Manager: CreateManager[graphic.Font](),
	}
}

func (fontManager *FontManager) RegisterFontFromImageFile(name string, filename string, width int32, height int32) {
	fontManager.Manager.Register(name,
		func() (*graphic.Font, error) {
			font, err := graphic.LoadFontFromBitmapFile(filename, width, height)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' font from file '%s'.\n - %w", name, filename, err)
			}
			return font, nil
		},
		func(font *graphic.Font) {
			if font != nil {
				font.Release()
			}
		})
}

func (fontManager *FontManager) RegisterFontFromBytes(name string, content []byte, format graphic.ImageFormat, width int32, height int32) {
	fontManager.Manager.Register(name,
		func() (*graphic.Font, error) {
			font, err := graphic.LoadFontFromBytes(content, format, width, height)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' font from byte array.\n - %w", name, err)
			}
			return font, nil
		},
		func(font *graphic.Font) {
			if font != nil {
				font.Release()
			}
		})
}
