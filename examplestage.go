package main

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"ogl46/engine"
	"ogl46/engine/graphic"
	"ogl46/engine/input"
	"ogl46/engine/ogl"
)

//go:embed resources/fonts/bitmap-test.png
var bitmapTestFontData []byte

//go:embed resources/images/brique.png
var briqueTextureData []byte

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

var ExampleStage engine.Stage = &exampleStage{}

type exampleStage struct {
	// bitmapTestFont est une police de caractères de test
	bitmapTestFont *graphic.Font
	// briqueTexture est une texture
	briqueTexture    *graphic.Texture
	trollfaceTexture *graphic.Texture
	crateTexture     *graphic.Texture
	// Shader 3D
	program *ogl.ShaderProgram
	// cube VAO
	cubeVao *ogl.VertexArray
	// fpsCounter permet de compter les FPS
	fpsCounter engine.FpsCounter
	// Caméra de la vue 3D
	camera graphic.Camera
}

func (stage *exampleStage) Initialize(app *engine.Application) error {
	var err error

	// Initialiser la police de caractères
	app.FontManager().RegisterFontFromBytes("bitmap-test", bitmapTestFontData, graphic.PNG, 16, 28)
	stage.bitmapTestFont, err = app.FontManager().Get("bitmap-test")
	if err != nil {
		return fmt.Errorf("Failed to load 'bitmap-test' font\n - %v", err)
	}

	// Initialiser la texture
	app.TextureManager().RegisterTextureFromBytes("brique", briqueTextureData, graphic.PNG)
	stage.briqueTexture, err = app.TextureManager().Get("brique")
	if err != nil {
		return fmt.Errorf("Failed to load 'brique' texture\n - %v", err)
	}
	app.TextureManager().RegisterTextureFromFile("crate", "./resources/images/RTS_Crate.png")
	stage.crateTexture, err = app.TextureManager().Get("crate")
	if err != nil {
		return fmt.Errorf("Failed to load 'crate' texture\n - %v", err)
	}
	app.TextureManager().RegisterTextureFromFile("trollface", "./resources/images/trollface-transparent.png")
	stage.trollfaceTexture, err = app.TextureManager().Get("trollface")
	if err != nil {
		return fmt.Errorf("Failed to load 'trollface' texture\n - %v", err)
	}

	vertShader, err := ogl.LoadShaderFromFile("./resources/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragShader, err := ogl.LoadShaderFromFile("./resources/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	stage.program, err = ogl.NewShaderProgram(vertShader, fragShader)
	if err != nil {
		return err
	}

	stage.cubeVao = ogl.NewVertexArrayWithTexture(cubeVertices, gl.TRIANGLES)

	// Initialiser le compteur de FPS
	stage.fpsCounter = engine.BuildFpsCounter()

	// Initialiser la caméra
	stage.camera = graphic.CreateCamera()
	stage.camera.Position = mgl32.Vec3{0., 3.5, -5.}
	stage.camera.Yaw = graphic.AngleInDegree(180.)

	return nil
}

func (stage *exampleStage) Release(app *engine.Application) {
	stage.cubeVao.Release()
	stage.program.Delete()
	app.TextureManager().Release("brique")
	app.FontManager().Release("bitmap-test")
}

func (stage *exampleStage) Display(app *engine.Application, _ *engine.Timer) {
	dimx, dimy := app.Size()
	fov := float32(60.0)

	// Affichage scene 3D
	stage.program.Use()

	// Matrice de projection
	projectTransform := mgl32.Perspective(mgl32.DegToRad(fov), float32(dimx)/float32(dimy), 0.1, 100.0)
	stage.program.UniformMatrix4fv("projection", &projectTransform)

	// Matrice de positionnement de la caméra
	camTransform := stage.camera.GetTransform() // FIXME Camera ici
	stage.program.UniformMatrix4fv("camera", &camTransform)

	// bind textures
	stage.crateTexture.BindTextureUnit(gl.TEXTURE0)
	stage.crateTexture.SetUniform(stage.program.GetUniformLocation("texture0"))

	stage.trollfaceTexture.BindTextureUnit(gl.TEXTURE1)
	stage.trollfaceTexture.SetUniform(stage.program.GetUniformLocation("texture1"))

	// cube rotation matrices
	rotateX := mgl32.Rotate3DX(mgl32.DegToRad(-60 * float32(glfw.GetTime())))
	rotateY := mgl32.Rotate3DY(mgl32.DegToRad(-60 * float32(glfw.GetTime())))
	rotateZ := mgl32.Rotate3DZ(mgl32.DegToRad(-60 * float32(glfw.GetTime())))

	stage.cubeVao.Bind()

	// draw each cube after all coordinate system transforms are bound
	for _, pos := range cubePositions {
		// Matrice de positionnement de l'élément dans le "monde"
		worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
		worldTransform := worldTranslate.Mul4(rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4())
		stage.program.UniformMatrix4fv("world", &worldTransform)

		// Dessin du cube
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}

	// bind textures
	stage.briqueTexture.BindTextureUnit(gl.TEXTURE0)
	stage.briqueTexture.SetUniform(stage.program.GetUniformLocation("texture0"))

	// draw each cube after all coordinate system transforms are bound
	for x := -10; x < 10; x++ {
		for z := -10; z < 10; z++ {
			// Matrice de positionnement de l'élément dans le "monde"
			worldTransform := mgl32.Translate3D(float32(x), -5., float32(z))
			stage.program.UniformMatrix4fv("world", &worldTransform)

			// Dessin du cube
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
	}

	gl.BindVertexArray(0)

	stage.crateTexture.UnBind()
	stage.trollfaceTexture.UnBind()

	// Et ensuite la 2D
	app.Renderer2d().Begin()
	for x := 0; x < 800; x = x + 20 {
		source := graphic.Rectangle{0., 0., float32(stage.trollfaceTexture.Width()), float32(stage.trollfaceTexture.Height())}
		target := graphic.Rectangle{float32(x), 0., 20., 20.}
		app.Renderer2d().DrawSpriteFromRect(stage.trollfaceTexture, source, target)
	}

	// Et ensuite la 2D
	app.Renderer2d().Begin()
	totalSize := 200
	squareSize := 20
	ratioX := float32(stage.trollfaceTexture.Width()) / float32(totalSize)
	ratioY := float32(stage.trollfaceTexture.Height()) / float32(totalSize)
	for x := 0; x < totalSize; x += squareSize {
		for y := 0; y < totalSize; y += squareSize {
			source := graphic.Rectangle{
				float32(x) * ratioX, float32(y) * ratioY,
				float32(squareSize)*ratioX - 1, float32(squareSize)*ratioY - 1,
			}
			target := graphic.Rectangle{
				float32(x), float32(y + 100),
				float32(squareSize - 1), float32(squareSize - 1),
			}
			app.Renderer2d().DrawSpriteFromRect(stage.trollfaceTexture, source, target)
		}
	}

	app.Renderer2d().DrawSprite(stage.trollfaceTexture, mgl32.Vec2{600., 50.})

	app.Renderer2d().DrawText(stage.bitmapTestFont, fmt.Sprintf("FPS: %.2f", stage.fpsCounter.ComputeFps()), mgl32.Vec2{20, 30}, graphic.White)
	app.Renderer2d().DrawText(stage.bitmapTestFont, fmt.Sprintf("Mouse: (%.0f, %.0f)", app.Mouse().X(), app.Mouse().Y()), mgl32.Vec2{20, 50}, graphic.White)

	app.Renderer2d().End()

}

func (stage *exampleStage) Execute(app *engine.Application, timer *engine.Timer) {
	// Move camera
	var x, y float32 = 0., 0.
	if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyD)) == input.KEY_PRESSED { // Right
		y = -1.
	}
	if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyA)) == input.KEY_PRESSED { // Left
		y = 1.
	}
	if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyW)) == input.KEY_PRESSED { // Up
		x = 1.
	}
	if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyS)) == input.KEY_PRESSED { // Down
		x = -1.
	}
	if x != 0 || y != 0 {
		move := mgl32.Vec2{x, y}.Normalize().Mul(float32(timer.ElapsedTime()) / 50)
		stage.camera.MoveForward(move.X())
		stage.camera.MoveLeft(move.Y())
	}
	if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyQ)) == input.KEY_PRESSED { // Pencher à gauche
		stage.camera.Roll += graphic.AngleInDegree(float64(timer.ElapsedTime()) / 5.)
		if stage.camera.Roll > graphic.AngleInDegree(45.) {
			stage.camera.Roll = graphic.AngleInDegree(45.)
		}
	} else if app.Keyboard().KeyStatus(input.KeyCode(glfw.KeyE)) == input.KEY_PRESSED { // Pencher à droite
		stage.camera.Roll -= graphic.AngleInDegree(float64(timer.ElapsedTime()) / 5.)
		if stage.camera.Roll < graphic.AngleInDegree(-45.) {
			stage.camera.Roll = graphic.AngleInDegree(-45.)
		}
	} else if stage.camera.Roll.Degree() != 0. {
		if stage.camera.Roll < 0. {
			stage.camera.Roll += graphic.AngleInDegree(float64(timer.ElapsedTime()) / 5.)
			if stage.camera.Roll > 0. {
				stage.camera.Roll = 0.
			}
		} else {
			stage.camera.Roll -= graphic.AngleInDegree(float64(timer.ElapsedTime()) / 5.)
			if stage.camera.Roll < 0. {
				stage.camera.Roll = 0.
			}
		}
	}
	// Compter les FPS
	stage.fpsCounter.Increment(timer.ElapsedTime())
}

func (stage *exampleStage) ProcessEvent(app *engine.Application, event input.Event, timer *engine.Timer) {
	switch event.Source() {
	case input.KEYBOARD:
		if event.(*input.KeyboardEvent).Action() == input.KEY_PRESS {
			switch event.(*input.KeyboardEvent).Key() {
			case glfw.KeyEscape:
				app.Stop()
			case glfw.KeyF11:
				app.SetFullScreen(!app.FullScreen())
			}
		}
	case input.MOUSE_BUTTON:
		if event.(*input.MouseButtonEvent).Action() == input.BUTTON_PRESS {
			app.Mouse().SetCursorHidden(true)
			app.Mouse().SetMoveMode(input.MOUSE_MOVE_RAW)
			app.Mouse().SetCursorLocked(true)
		} else if event.(*input.MouseButtonEvent).Action() == input.BUTTON_RELEASE {
			app.Mouse().SetCursorLocked(false)
			app.Mouse().SetMoveMode(input.MOUSE_MOVE_CURSOR)
			app.Mouse().SetCursorHidden(false)
		}
	case input.MOUSE_MOVE:
		mme := event.(*input.MouseMoveEvent)
		if app.Mouse().CursorLocked() {
			log.Printf("Before - Yaw: %v / Pitch: %v", stage.camera.Yaw.Degree(), stage.camera.Pitch.Degree())
			stage.camera.Yaw = (stage.camera.Yaw - graphic.AngleInDegree(mme.HShift()/10.)).Mod360()
			stage.camera.Pitch = (stage.camera.Pitch + graphic.AngleInDegree(mme.VShift()/10.)).Mod360()
			log.Printf("After - Yaw: %v / Pitch: %v", stage.camera.Yaw.Degree(), stage.camera.Pitch.Degree())
		}
	/*case input.MOUSE_SCROLL:
	mse := event.(*input.MouseScrollEvent)
	stage.camera.Fov = stage.camera.Fov + geometry.AngleInDegree(mse.VShift())
	logs.Infof("New FOV: %0.00f", stage.camera.Fov)*/
	case input.WINDOW:
		if event.(*input.WindowEvent).Action() == input.WINDOW_CLOSED {
			app.Stop()
		}
	}
}
