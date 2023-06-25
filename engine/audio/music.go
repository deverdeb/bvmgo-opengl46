package assets

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"io"
	"log/slog"
	"os"
)

// Music permet de jouer de la musique.
type Music struct {
	streamer   beep.StreamSeekCloser
	controller *beep.Ctrl
	volume     *effects.Volume
}

func LoadMusicFromFile(filename string) (*Music, error) {
	format, found := findAudioFormatFromFilename(filename)
	if !found {
		return nil, fmt.Errorf("cannot build sound from '%s' file "+
			"(unsupported '%s' format - only '.ogg' and '.mp3' are supported)", filename, format)
	}
	fileReader, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s' music file\n - %w", filename, err)
	}
	music, err := loadMusicFromReader(fileReader, format)
	if err != nil {
		return nil, fmt.Errorf("failed to load music from '%s' file\n - %w", filename, err)
	}
	return music, nil
}

func LoadMusicFromOggFile(filename string) (*Music, error) {
	slog.Debug("music creation from OGG file", "filename", filename)
	oggFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s' OGG file\n - %w", filename, err)
	}
	music, err := loadMusicFromReader(oggFile, OGG)
	if err != nil {
		return nil, fmt.Errorf("failed to load music from '%s' OGG file\n - %w", filename, err)
	}
	return music, nil
}

func LoadMusicFromMp3File(filename string) (*Music, error) {
	slog.Debug("music creation from MP3 file", "filename", filename)
	mp3File, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s' MP3 file\n - %w", filename, err)
	}
	music, err := loadMusicFromReader(mp3File, MP3)
	if err != nil {
		return nil, fmt.Errorf("failed to load music from '%s' MP3 file\n - %w", filename, err)
	}
	return music, nil
}

func LoadMusicFromBytes(content []byte, format AudioFormat) (*Music, error) {
	bytesReader := io.NopCloser(bytes.NewReader(content))
	music, err := loadMusicFromReader(bytesReader, format)
	if err != nil {
		return nil, fmt.Errorf("failed to load '%s' music from byte array\n - %w", format, err)
	}
	return music, nil
}

func loadMusicFromReader(audioReader io.ReadCloser, format AudioFormat) (*Music, error) {
	slog.Debug("music creation from data reader", "format", format)
	decoder, found := findAudioDecoderFromAudioFormat(format)
	if !found {
		return nil, fmt.Errorf("cannot build music from byte array "+
			"(unsupported '%s' format - only '.ogg' and '.mp3' are supported)", format)
	}
	streamer, _, err := decoder(audioReader)
	if err != nil {
		return nil, fmt.Errorf("failed to load '%s' music\n - %w", format, err)
	}
	return buildMusicFromStreamer(streamer), nil
}

func buildMusicFromStreamer(streamer beep.StreamSeekCloser) *Music {
	controller := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	controller.Paused = true
	volume := &effects.Volume{Streamer: controller, Base: 1}
	beepInit()
	speaker.Play(volume)
	return &Music{
		streamer:   streamer,
		controller: controller,
		volume:     volume,
	}
}

func (music *Music) Release() {
	if music.streamer != nil {
		slog.Debug("music destruction")
		err := music.streamer.Close()
		if err != nil {
			slog.Error("music destruction failed", "error", err)
		}
	}
}

func (music *Music) Play() {
	speaker.Lock()
	music.controller.Paused = false
	speaker.Unlock()
}

func (music *Music) Pause() {
	speaker.Lock()
	music.controller.Paused = true
	speaker.Unlock()
}

func (music *Music) Stop() {
	speaker.Lock()
	music.controller.Paused = true
	_ = music.streamer.Seek(0)
	speaker.Unlock()
}

func (music *Music) SetVolume(volume int) {
	speaker.Lock()
	music.volume.Volume = float64(volume) / 100.
	speaker.Unlock()
}
