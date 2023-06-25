package graphic

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log/slog"
	"ogl46/engine/ogl"
)

// renderer2dVertexShader is the 2D vertex shader.
var renderer2dVertexShader = `
#version 460 core

layout (location = 0) in vec2 position;     // Vertex Array - Position 
layout (location = 1) in vec2 textPosition; // Vertex Array - Texture coordinates (in)

out vec2 TexCoords;       // Texture coordinates (out)

uniform mat4 projection;  // Orthogonal projection for 2d rendering
uniform mat4 coordinates; // Image coordinate on screen

void main()
{
    gl_Position = projection * coordinates * vec4(position.xy, 0.0, 1.0);
    TexCoords = textPosition;
}
`

// renderer2dFragmentShader is the 2D fragment shader.
var renderer2dFragmentShader = `
#version 460 core

in vec2 TexCoords; // Texture coordinates
out vec4 color;    // Pixel color

uniform sampler2D image;      // Texture
uniform vec4 spriteColor;     // User color
uniform vec2 texturePosition; // Position on texture
uniform vec2 textureRatio;    // Texture ratio (for resize)

void main()
{
	// compute coordinates on texture
	vec2 coordinates = (TexCoords * textureRatio) + texturePosition;
	// mix user color and texture color
    color = spriteColor * texture(image, coordinates);
}
`

// renderer2dVertices contains 2D square coordinates (and texture coordinates).
var renderer2dVertices = []float32{
	// pos    // tex
	0.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 0.0,

	0.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
}

const TabulationSize = 4

// Renderer2d contains methods to draw a 2D scene.
type Renderer2d struct {
	// 2D shaderProgram handler.
	shaderProgram *ogl.ShaderProgram
	//  square vao (Vertex Array Object) handler.
	vao uint32
	//  square vbo (Vertex Buffer Object) handler.
	vbo uint32

	// State before 2D view (to restore state after 2d rendering)
	keepPreviousStateBlend     bool
	keepPreviousStateDepthTest bool
}

// Initialize renderer (create vertex array and shader).
func (renderer *Renderer2d) Initialize() error {
	renderer.initializeVao()
	return renderer.initializeShader()
}

// initializeShader compiles vertex shader and fragment et creates shader program.
func (renderer *Renderer2d) initializeShader() error {
	slog.Debug("Renderer2d - 2D shader creation")
	// Créer le shader
	vertexShader, err := ogl.NewShaderFromSource(renderer2dVertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragmentShader, err := ogl.NewShaderFromSource(renderer2dFragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	renderer.shaderProgram, err = ogl.NewShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}
	return nil
}

// initializeVao defines vertex array.
func (renderer *Renderer2d) initializeVao() {
	slog.Debug("Renderer2d - 2D vao creation")
	// Generate VAO and VBO
	gl.GenVertexArrays(1, &renderer.vao)
	gl.GenBuffers(1, &renderer.vbo)
	// Active VAO
	gl.BindVertexArray(renderer.vao)

	// Copy all square data in VBO
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(renderer2dVertices)*4, gl.Ptr(renderer2dVertices), gl.STATIC_DRAW) // 4 = float32 size

	var offset uintptr = 0       // Beginning offset
	var stride int32 = 2*4 + 2*4 // Vertex size = (2 coordinates + 2 texture coordinates) * 4 float32 size

	// position coordinates
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, stride, offset) // 2 = coordinates number
	gl.EnableVertexAttribArray(0)
	offset += 2 * 4 // Move offset by coordinates number

	// texture coordinates
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, offset) // 2 = coordinates number
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4 // Move offset by coordinates number

	// Unbind VAO after VAO creation
	gl.BindVertexArray(0)
}

// Release renderer (release VAO and shader).
func (renderer *Renderer2d) Release() {
	renderer.releaseVao()
	if renderer.shaderProgram != nil {
		renderer.releaseShader()
	}
}

func (renderer *Renderer2d) releaseShader() {
	slog.Debug("Renderer2d - 2D shader destruction")
	renderer.shaderProgram.Delete()
}

func (renderer *Renderer2d) releaseVao() {
	if renderer.vao != 0 {
		slog.Debug("Renderer2d - 2D vertex array destruction")
		gl.DeleteBuffers(1, &renderer.vbo)
		renderer.vbo = 0
		gl.DeleteVertexArrays(1, &renderer.vao)
		renderer.vao = 0
	}
}

// Begin method initializes 2D rendering (active shader and OpenGL states)
func (renderer *Renderer2d) Begin() {
	// Disable depth test
	renderer.keepPreviousStateDepthTest = gl.IsEnabled(gl.DEPTH_TEST)
	if renderer.keepPreviousStateDepthTest {
		gl.Disable(gl.DEPTH_TEST)
	}

	// Active transparence
	renderer.keepPreviousStateBlend = gl.IsEnabled(gl.BLEND)
	if !renderer.keepPreviousStateBlend {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}

	// Active shader
	renderer.shaderProgram.Use()

	// Initialize orthogonal projection
	projection2d := mgl32.Ortho2D(0., 800., 600., 0.)
	renderer.shaderProgram.UniformMatrix4fv("projection", &projection2d)
}

// End method finalizes 2D rendering (restore OpenGL states)
func (renderer *Renderer2d) End() {
	renderer.shaderProgram.Unuse()
	if renderer.keepPreviousStateDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	}
	if !renderer.keepPreviousStateBlend {
		gl.Disable(gl.BLEND)
	}
}

// DrawSprite draws a picture
func (renderer *Renderer2d) DrawSprite(texture *Texture, targetPosition mgl32.Vec2) {
	// Screen / target zone
	targetDimension := mgl32.Vec2{float32(texture.Width()), float32(texture.Height())}
	// Texture zone
	sourcePosition := mgl32.Vec2{0., 0.}
	sourceDimension := targetDimension

	renderer.DrawSpriteExWithRotateAndColor(texture, sourcePosition, sourceDimension,
		targetPosition, targetDimension, mgl32.Vec2{}, 0., White)
}

// DrawSpriteFromRect draws a picture
func (renderer *Renderer2d) DrawSpriteFromRect(texture *Texture, sourceRectangle Rectangle, targetRectangle Rectangle) {
	renderer.DrawSpriteExWithRotateAndColor(texture, sourceRectangle.Pos(), sourceRectangle.Dim(),
		targetRectangle.Pos(), targetRectangle.Dim(), mgl32.Vec2{}, 0., White)
}

// DrawSpriteFromRectWithColor colorizes and draws a picture
func (renderer *Renderer2d) DrawSpriteFromRectWithColor(texture *Texture, sourceRectangle Rectangle, targetRectangle Rectangle, color Color) {
	renderer.DrawSpriteExWithRotateAndColor(texture, sourceRectangle.Pos(), sourceRectangle.Dim(),
		targetRectangle.Pos(), targetRectangle.Dim(), mgl32.Vec2{}, 0., color)
}

// DrawSpriteExWithRotateAndColor draws part of texture at position and applies a color and a rotation.
//
//	sourcePosition et sourceDimension définissent la zone à rendre au niveau de la texture.
//	targetPosition et targetDimension définissent la zone où effectuer le rendu à l'écran.
//	rotationCenter indique le centre de rotation à utiliser (ratio relatif à la taille de la zone à dessiner : 0. = à gauche ou en haut / 1. = à droite ou en bas)
//	rotationInDegree est la rotation à appliquer (angle en degrés)
//	color est la couleur à appliquer
func (renderer *Renderer2d) DrawSpriteExWithRotateAndColor(texture *Texture, sourcePosition mgl32.Vec2, sourceDimension mgl32.Vec2,
	targetPosition mgl32.Vec2, targetDimension mgl32.Vec2, rotationCenter mgl32.Vec2, rotation Angle, color Color) {

	// Activer la texture sur l'unité de texture 0 et l'indiquer au Shader
	texture.Bind()
	renderer.shaderProgram.Uniform1i("image", 0)

	// Positionner la couleur
	renderer.shaderProgram.UniformVector4f("spriteColor", mgl32.Vec4(color))

	// Se positionner sur la texture (écart avant / après)
	textRatioInfo := mgl32.Vec2{
		sourceDimension.X() / float32(texture.Width()),
		sourceDimension.Y() / float32(texture.Height()),
	}
	textCoordInfo := mgl32.Vec2{
		sourcePosition.X() / float32(texture.Width()),
		sourcePosition.Y() / float32(texture.Height()),
	}

	// Et passer les infos de texture au shader
	renderer.shaderProgram.UniformVector2f("texturePosition", textCoordInfo)
	renderer.shaderProgram.UniformVector2f("textureRatio", textRatioInfo)

	// Positionner le modèle
	//glm::mat4 model = glm::mat4(1.0f);
	//model = glm::translate(model, glm::vec3(targetPosition, 0.0f));
	model := mgl32.Translate3D(targetPosition.X(), targetPosition.Y(), 0.)

	// Rotation
	if rotation != 0. {
		model = model.Mul4(mgl32.Translate3D(rotationCenter.X(), rotationCenter.Y(), 0.))   // move origin of rotation to center of quad
		model = model.Mul4(mgl32.HomogRotate3DZ(-float32(rotation)))                        // then rotate
		model = model.Mul4(mgl32.Translate3D(-rotationCenter.X(), -rotationCenter.Y(), 0.)) // move origin back
	}

	// Dimensionner le sprite
	model = model.Mul4(mgl32.Scale3D(targetDimension.X(), targetDimension.Y(), 1.)) // last scale

	// Et passer la transformation au shader
	renderer.shaderProgram.UniformMatrix4fv("coordinates", &model)

	// Faire le rendu du VAO
	gl.BindVertexArray(renderer.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}

func (renderer *Renderer2d) DrawText(font *Font, text string, position mgl32.Vec2, color Color) {
	// Activer la texture
	font.Texture().Bind()

	destination := BuildRectangle(position.X(), position.Y(), float32(font.CharacterWidth()), float32(font.CharacterHeight()))
	for _, character := range text {
		switch character {
		case '\n':
			destination.SetX(position.X())
			destination.SetY(destination.Y() + float32(font.CharacterHeight()))
		case '\t': // Tabulation
			for iter := 0; iter < TabulationSize; iter++ {
				renderer.drawCharacter(font, ' ', destination, color)
				destination.SetX(destination.X() + float32(font.CharacterWidth()))
			}
		case '\r', '\b', '\a':
			// Ignorer...
		default:
			renderer.drawCharacter(font, character, destination, color)
			destination.SetX(destination.X() + float32(font.CharacterWidth()))
		}
	}
}

func (renderer *Renderer2d) drawCharacter(font *Font, character int32, destination Rectangle, color Color) {
	source := font.CharacterRectangle(character)
	renderer.DrawSpriteFromRectWithColor(font.Texture(), source, destination, color)
}
