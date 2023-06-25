package assetsmngr

import (
	"fmt"
	"ogl46/engine/audio"
)

// SoundManager permet de gérer le chargement et la libération des sons
type SoundManager struct {
	// Manager est une instance d'asset manager pour gérer les ressources
	Manager[assets.Sound]
}

func NewSoundManager() *SoundManager {
	return &SoundManager{
		Manager: CreateManager[assets.Sound](),
	}
}

func (soundManager *SoundManager) RegisterSoundFromFile(name string, filename string) {
	soundManager.Manager.Register(name,
		func() (*assets.Sound, error) {
			sound, err := assets.LoadSoundFromFile(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' sound from file '%s'.\n - %w", name, filename, err)
			}
			return sound, nil
		},
		func(sound *assets.Sound) {
			if sound != nil {
				sound.Release()
			}
		})
}

func (soundManager *SoundManager) RegisterSoundFromBytes(name string, data []byte, format assets.AudioFormat) {
	soundManager.Manager.Register(name,
		func() (*assets.Sound, error) {
			sound := assets.LoadSoundFromBytes(data, format)
			return sound, nil
		},
		func(sound *assets.Sound) {
			if sound != nil {
				sound.Release()
			}
		})
}
