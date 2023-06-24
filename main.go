package main

/*
Adapted from this tutorial: http://www.learnopengl.com/#!Getting-started/Camera

Shows how to create a basic controllable FPS camera. This has been refactored into
classes to allow better reuse going forward.
*/

import (
	"log"
	"log/slog"
	"ogl46/assets"
	"ogl46/cam"
	"ogl46/win"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// vertices to draw 6 faces of a cube
var cubeVertices = []float32{
	// position        // texture position
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5, 0.5, -0.5, 0.0, 1.0,
	0.5, 0.5, -0.5, 1.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0,
}

var cubePositions = [][]float32{
	{0.0, 0.0, -3.0},
	{2.0, 5.0, -15.0},
	{-1.5, -2.2, -2.5},
	{-3.8, -2.0, -12.3},
	{2.4, -0.4, -3.5},
	{-1.7, 3.0, -7.5},
	{1.3, -2.0, -2.5},
	{1.5, 2.0, -2.5},
	{1.5, 0.2, -1.5},
	{-1.3, 1.0, -1.5},
}

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window := win.NewWindow(1280, 720, "basic camera")

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	vendor := gl.GoStr(gl.GetString(gl.VENDOR))
	slog.Info("OpenGL", "version", version, "vendor", vendor)

	err := programLoop(window)
	if err != nil {
		log.Fatalln(err)
	}
}

func programLoop(window *win.Window) error {

	// Début 2D
	renderer2d := Renderer2d{}
	err := renderer2d.Initialize()
	defer renderer2d.Release()
	if err != nil {
		return err
	}
	// Fin 2D	// the linked shader program determines how the data will be rendered
	vertShader, err := assets.LoadShaderFromFile("./resources/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragShader, err := assets.LoadShaderFromFile("./resources/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	program, err := assets.NewShaderProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer program.Delete()

	VAO := assets.NewVertexArrayWithTexture(cubeVertices, gl.TRIANGLES)
	texture0, err := assets.LoadTextureFromFile("./resources/images/RTS_Crate.png")
	if err != nil {
		panic(err.Error())
	}

	texture1, err := assets.LoadTextureFromFile("./resources/images/trollface-transparent.png")
	if err != nil {
		panic(err.Error())
	}

	// ensure that triangles that are "behind" others do not draw over top of them
	gl.Enable(gl.DEPTH_TEST)

	camera := cam.NewFpsCamera(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 1, 0}, -90, 0, window.InputManager())

	for !window.ShouldClose() {

		// swaps in last buffer, polls for window events, and generally sets up for a new render frame
		window.StartFrame()

		// update camera position and direction from input evevnts
		camera.Update(window.SinceLastFrame())

		// background color
		gl.ClearColor(0., 0., 0., 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // depth buffer needed for DEPTH_TEST
		gl.Enable(gl.DEPTH_TEST)

		// Culling (masquer les faces non visibles)
		gl.Disable(gl.CULL_FACE)
		//gl.Enable(gl.CULL_FACE)
		//gl.CullFace(gl.BACK)

		program.Use()

		// bind textures
		texture0.BindTextureUnit(gl.TEXTURE0)
		texture0.SetUniform(program.GetUniformLocation("texture0"))

		texture1.BindTextureUnit(gl.TEXTURE1)
		texture1.SetUniform(program.GetUniformLocation("texture1"))

		// cube rotation matrices
		rotateX := (mgl32.Rotate3DX(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateY := (mgl32.Rotate3DY(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))
		rotateZ := (mgl32.Rotate3DZ(mgl32.DegToRad(-60 * float32(glfw.GetTime()))))

		// creates perspective
		fov := float32(60.0)
		projectTransform := mgl32.Perspective(mgl32.DegToRad(fov),
			float32(window.Width())/float32(window.Height()),
			0.1,
			100.0)

		camTransform := camera.GetTransform()
		program.UniformMatrix4fv("camera", &camTransform)
		program.UniformMatrix4fv("projection", &projectTransform)

		VAO.Bind()

		// draw each cube after all coordinate system transforms are bound
		for _, pos := range cubePositions {
			worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			worldTransform := (worldTranslate.Mul4(rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4()))

			program.UniformMatrix4fv("world", &worldTransform)

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		gl.BindVertexArray(0)

		texture0.UnBind()
		texture1.UnBind()

		// Et ensuite la 2D
		renderer2d.Begin()
		renderer2d.DrawSprite(texture1 /*, ...*/)
		renderer2d.End()

		// end of draw loop
	}

	return nil
}
