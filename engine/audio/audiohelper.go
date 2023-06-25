package assets

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
)

const beepSampleRate = beep.SampleRate(44100)

// AudioFormat is an audio file supported format
type AudioFormat string

// List of supported image formats
const (
	OGG AudioFormat = ".ogg"
	MP3             = ".mp3"
)

// decoders contains image decoder by supported format
var audioDecoders = map[AudioFormat]func(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error){
	OGG: vorbis.Decode,
	MP3: mp3.Decode,
}

var beepInitialized = false

func beepInit() {
	if !beepInitialized {
		err := speaker.Init(beepSampleRate, beepSampleRate.N(time.Second/10))
		if err != nil {
			log.Fatalf("beep speaker initialization failed\n - %v", err)
		}
		beepInitialized = true
	}
}

func findAudioDecoderFromAudioFormat(audioFormat AudioFormat) (func(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error), bool) {
	decoder, found := audioDecoders[audioFormat]
	return decoder, found
}

func findAudioFormatFromFilename(filename string) (AudioFormat, bool) {
	extension := strings.ToLower(strings.TrimSpace(filepath.Ext(filename)))
	switch extension {
	case ".ogg":
		return OGG, true
	case ".mp3":
		return MP3, true
	default:
		return "", false
	}
}
