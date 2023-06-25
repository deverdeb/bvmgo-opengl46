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

// Sound permet de jouer des sons.
type Sound struct {
	source  string
	content []byte
	format  AudioFormat
	volume  float64
}

func LoadSoundFromFile(filename string) (*Sound, error) {
	format, found := findAudioFormatFromFilename(filename)
	if !found {
		return nil, fmt.Errorf("cannot build sound from '%s' file "+
			"(unsupported '%s' format - only '.ogg' and '.mp3' are supported)", filename, format)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound file '%s'\n - %w", filename, err)
	}
	return &Sound{
		source:  "'" + filename + "' file",
		content: content,
		format:  format,
		volume:  1,
	}, nil
}

func LoadSoundFromOggFile(filename string) (*Sound, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound file '%s'\n - %w", filename, err)
	}
	return &Sound{
		source:  "'" + filename + "' file",
		content: content,
		format:  OGG,
		volume:  1,
	}, nil
}

func LoadSoundFromMp3File(filename string) (*Sound, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound file '%s'\n - %w", filename, err)
	}
	return &Sound{
		source:  "'" + filename + "' file",
		content: content,
		format:  MP3,
		volume:  1,
	}, nil
}

func LoadSoundFromBytes(content []byte, format AudioFormat) *Sound {
	var source string
	switch format {
	case OGG:
		source = "ogg []bytes"
	case MP3:
		source = "mp3 []bytes"
	}
	return &Sound{
		source:  source,
		content: content,
		format:  format,
		volume:  1,
	}
}

func (sound *Sound) buildStreamer() (beep.StreamSeekCloser, error) {
	// Build byte reader
	reader := io.NopCloser(bytes.NewReader(sound.content))
	// Build audio streamer
	decoder, _ := findAudioDecoderFromAudioFormat(sound.format)
	streamer, _, err := decoder(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode '%s' sound\n - %w", sound.source, err)
	}
	return streamer, nil
}

func (sound *Sound) Release() {
}

func (sound *Sound) Play() {
	streamer, err := sound.buildStreamer()
	if err != nil {
		slog.Error("play sound failed", "source", sound.source, "error", err)
		return
	}
	volume := &effects.Volume{Streamer: streamer, Base: sound.volume}
	beepInit()
	speaker.Play(beep.Seq(volume, beep.Callback(func() {
		err := streamer.Close()
		if err != nil {
			slog.Error("close sound failed", "source", sound.source, "error", err)
		}
	})))
}

func (sound *Sound) SetVolume(volume int) {
	sound.volume = float64(volume) / 100.
}
