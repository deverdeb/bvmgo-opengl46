package assetsmngr

import (
	"fmt"
	"ogl46/engine/audio"
)

// MusicManager permet de gérer le chargement et la libération des musiques
type MusicManager struct {
	// Manager est une instance d'asset manager pour gérer les ressources
	Manager[assets.Music]
}

func NewMusicManager() *MusicManager {
	return &MusicManager{
		Manager: CreateManager[assets.Music](),
	}
}

func (musicManager *MusicManager) RegisterMusicFromFile(name string, filename string) {
	musicManager.Manager.Register(name,
		func() (*assets.Music, error) {
			music, err := assets.LoadMusicFromFile(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' music from file '%s'.\n - %w", name, filename, err)
			}
			return music, nil
		},
		func(music *assets.Music) {
			if music != nil {
				music.Release()
			}
		})
}

func (musicManager *MusicManager) RegisterMusicFromBytes(name string, data []byte, format assets.AudioFormat) {
	musicManager.Manager.Register(name,
		func() (*assets.Music, error) {
			music, err := assets.LoadMusicFromBytes(data, format)
			if err != nil {
				return nil, fmt.Errorf("failed to load '%s' music from byte array. - %w", name, err)
			}
			return music, nil
		},
		func(music *assets.Music) {
			if music != nil {
				music.Release()
			}
		})
}
