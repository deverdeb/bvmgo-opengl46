package graphic

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"image"
	"image/draw"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Texture struct {
	// Adresse OpenGL
	handle uint32
	// Unité de texture
	unit uint32 // Texture unit that is currently bound to ex: gl.TEXTURE0
	// Dimension de la texture : largeur
	width int32
	// Dimension de la texture : hauteur
	height int32
}

func LoadTextureFromFile(filename string) (*Texture, error) {
	slog.Debug("texture creation from file %s", filename)
	extension := ImageFormat(strings.ToLower(filepath.Ext(filename)))
	decoder, found := imageDecoders[extension]
	if !found {
		return nil, fmt.Errorf("cannot build texture from '%s' file "+
			"(unsupported '%s' extension - only '.png', '.jpeg' and '.gif' are supported)", filename, extension)
	}
	fileReader, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from file '%s'\n - %w", filename, err)
	}
	texture, err := buildTexture(fileReader, decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to build texture from file '%s'\n - %w", filename, err)
	}
	return texture, nil
}

func LoadTextureFromBytes(content []byte, format ImageFormat) (*Texture, error) {
	slog.Debug("texture creation from byte array")
	decoder, found := imageDecoders[format]
	if !found {
		return nil, fmt.Errorf("cannot build texture from byte array "+
			"(unsupported '%s' format - only '.png', '.jpeg' and '.gif' are supported)", format)
	}
	textureReader := bytes.NewReader(content)
	texture, err := buildTexture(textureReader, decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to build texture from byte array\n - %w", err)
	}
	return texture, nil
}

func NewTextureFromImage(img image.Image) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported image stride (stride value = %d)", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	var handle uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &handle)
	gl.BindTexture(gl.TEXTURE_2D, handle)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	slog.Debug("texture info: handler=%d / size=%dx%d", handle, rgba.Rect.Size().X, rgba.Rect.Size().Y)

	return &Texture{
		handle: handle,
		width:  int32(rgba.Rect.Size().X),
		height: int32(rgba.Rect.Size().Y),
	}, nil
}

func (texture *Texture) Bind() {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.handle)
	texture.unit = gl.TEXTURE0
}

func (texture *Texture) BindTextureUnit(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, texture.handle)
	texture.unit = textureUnit
}

func (texture *Texture) UnBind() {
	gl.ActiveTexture(texture.unit)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (texture *Texture) Release() {
	if texture.handle != 0 {
		slog.Debug("texture destruction")
		gl.DeleteTextures(1, &texture.handle)
		texture.handle = 0
	}
}

func (texture *Texture) SetUniform(uniformLocation int32) error {
	if texture.unit == 0 {
		return fmt.Errorf("failed to modify uniform location on without texture unit")
	}
	gl.Uniform1i(uniformLocation, int32(texture.unit-gl.TEXTURE0))
	return nil
}

func (texture *Texture) Width() int32 {
	return texture.width
}

func (texture *Texture) Height() int32 {
	return texture.height
}

func (texture *Texture) Rectangle() Rectangle {
	return Rectangle{0, 0, float32(texture.width), float32(texture.height)}
}

func buildTexture(imageReader io.Reader, decode func(reader io.Reader) (image.Image, error)) (*Texture, error) {
	img, err := decode(imageReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image\n - %w", err)
	}
	texture, err := NewTextureFromImage(img)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture from image\n - %w", err)
	}
	return texture, nil
}
