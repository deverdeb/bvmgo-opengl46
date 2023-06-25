package ogl

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderProgram struct {
	handle  uint32
	shaders []*Shader
}

func NewShaderProgram(shaders ...*Shader) (*ShaderProgram, error) {
	program := &ShaderProgram{
		handle:  gl.CreateProgram(),
		shaders: make([]*Shader, 0),
	}
	program.attach(shaders...)
	if err := program.link(); err != nil {
		return nil, err
	}
	return program, nil
}

func (shaderProgram *ShaderProgram) Delete() {
	if shaderProgram.handle != 0 {
		for _, shader := range shaderProgram.shaders {
			shader.Delete()
		}
		gl.DeleteProgram(shaderProgram.handle)
	}
}

func (shaderProgram *ShaderProgram) link() error {
	gl.LinkProgram(shaderProgram.handle)
	return getGlError(shaderProgram.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog)
}

func (shaderProgram *ShaderProgram) attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(shaderProgram.handle, shader.handle)
		shaderProgram.shaders = append(shaderProgram.shaders, shader)
	}
}

func (shaderProgram *ShaderProgram) Use() {
	gl.UseProgram(shaderProgram.handle)
}

func (shaderProgram *ShaderProgram) Unuse() {
	gl.UseProgram(0)
}

func (shaderProgram *ShaderProgram) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(shaderProgram.handle, gl.Str(name+"\x00"))
}

func (shaderProgram *ShaderProgram) Uniform1i(name string, value int32) {
	gl.Uniform1i(shaderProgram.GetUniformLocation(name), value)
}

func (shaderProgram *ShaderProgram) Uniform1f(name string, value float32) {
	gl.Uniform1f(shaderProgram.GetUniformLocation(name), value)
}

func (shaderProgram *ShaderProgram) Uniform2f(name string, x, y float32) {
	gl.Uniform2f(shaderProgram.GetUniformLocation(name), x, y)
}

func (shaderProgram *ShaderProgram) UniformVector2f(name string, vec2 mgl32.Vec2) {
	gl.Uniform2f(shaderProgram.GetUniformLocation(name), vec2.X(), vec2.Y())
}

func (shaderProgram *ShaderProgram) Uniform3f(name string, x, y, z float32) {
	gl.Uniform3f(shaderProgram.GetUniformLocation(name), x, y, z)
}

func (shaderProgram *ShaderProgram) UniformVector3f(name string, vec3 mgl32.Vec3) {
	gl.Uniform3f(shaderProgram.GetUniformLocation(name), vec3.X(), vec3.Y(), vec3.Z())
}

func (shaderProgram *ShaderProgram) Uniform4f(name string, x, y, z, w float32) {
	gl.Uniform4f(shaderProgram.GetUniformLocation(name), x, y, z, w)
}

func (shaderProgram *ShaderProgram) UniformVector4f(name string, vec4 mgl32.Vec4) {
	gl.Uniform4f(shaderProgram.GetUniformLocation(name), vec4.X(), vec4.Y(), vec4.Z(), vec4.W())
}

func (shaderProgram *ShaderProgram) UniformMatrix4fv(name string, mat4 *mgl32.Mat4) {
	gl.UniformMatrix4fv(shaderProgram.GetUniformLocation(name), 1, false, &mat4[0])
}
