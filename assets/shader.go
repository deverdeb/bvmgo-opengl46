package assets

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Shader struct {
	handle uint32
}

func NewShaderFromSource(src string, shaderType uint32) (*Shader, error) {
	handle := gl.CreateShader(shaderType)
	glSrc, freeFn := gl.Strs(src + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err := getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog)
	if err != nil {
		return nil, fmt.Errorf("failed to compile shader from sources\n - %w", err)
	}
	return &Shader{handle: handle}, nil
}

func LoadShaderFromFile(filename string, shaderType uint32) (*Shader, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load '%s' shader file\n - %w", filename, err)
	}
	shader, err := NewShaderFromSource(string(fileContent), shaderType)
	if err != nil {
		return nil, fmt.Errorf("failed to build shader from '%s' file\n - %w", filename, err)
	}
	return shader, nil
}

func LoadShaderFromBytes(content []byte, shaderType uint32) (*Shader, error) {
	shader, err := NewShaderFromSource(string(content), shaderType)
	if err != nil {
		return nil, fmt.Errorf("failed to build shader byte array\n - %w", err)
	}
	return shader, nil
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.handle)
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv, getObjInfoLogFn getObjInfoLog) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, log)

		return fmt.Errorf(gl.GoStr(log))
	}

	return nil
}
