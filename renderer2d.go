package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"ogl46/assets"
)

var renderer2dVertexShader = `
#version 410 core

layout (location = 0) in vec2 position;   // position
layout (location = 1) in vec2 textPosition; // coordonnées de la texture

out vec2 TexCoords;

uniform mat4 coordinates;
uniform mat4 projection;

void main()
{
    gl_Position = projection * coordinates * vec4(position.xy, 0.0, 1.0);
    TexCoords = textPosition;
}
`

var renderer2dFragmentShader = `
#version 410 core

in vec2 TexCoords;
out vec4 color;

uniform sampler2D image;
uniform vec3 spriteColor;

void main()
{
    color = vec4(spriteColor, 1.0) * texture(image, TexCoords);
	// color = vec4(spriteColor, 1.0) + texture(image, TexCoords);
}
`

var renderer2dVertices = []float32{
	// pos    // tex
	0.0, 1.0, 0.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
	0.0, 0.0, 0.0, 0.0,

	0.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 0.0,
}

type Renderer2d struct {
	// Shader de rendu 2D
	shaderProgram *assets.ShaderProgram
	// vao (Vertex Array Object)
	vao uint32
	// vbo (Vertex Buffer Object)
	vbo uint32

	// Etats AVANT passage en 2D
	keepPreviousStateBlend     bool
	keepPreviousStateDepthTest bool
}

func (renderer *Renderer2d) Initialize() error {
	renderer.initializeVao()
	return renderer.initializeShader()
}

func (renderer *Renderer2d) initializeShader() error {
	// Créer le shader
	vertexShader, err := assets.NewShaderFromSource(renderer2dVertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragmentShader, err := assets.NewShaderFromSource(renderer2dFragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	renderer.shaderProgram, err = assets.NewShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	// Configurer le shader
	//renderer.shaderProgram.Uniform1i("image", 0)
	//projection2d := mgl32.Ortho2D(0., 800., 600., 0.)
	//renderer.shaderProgram.UniformMatrix4fv("projection", &projection2d)

	return nil
}

func (renderer *Renderer2d) initializeVao() {
	// Générer les objets OpenGL
	gl.GenVertexArrays(1, &renderer.vao)
	gl.GenBuffers(1, &renderer.vbo)
	// Activer le vao
	gl.BindVertexArray(renderer.vao)
	// Copier les données dans le vbo
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(renderer2dVertices)*4, gl.Ptr(renderer2dVertices), gl.STATIC_DRAW) // 4 = float32 size

	var offset uintptr = 0       // Au début, l'offset est au début
	var stride int32 = 2*4 + 2*4 // Taille des données d'un vertex = 2 coordonnées + 2 coordonnées de texture, 4 = taille d'un float32

	// positions
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, stride, offset) // 2 = nb coordonnées
	gl.EnableVertexAttribArray(0)
	offset += 2 * 4 // Décalage de l'offset des 2 coordonnées

	// coordonnées
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, offset) // 2 = nb coordonnées
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4 // Décalage de l'offset des 2 coordonnées de texture

	// Unbind du VAO à la fin de la création
	gl.BindVertexArray(0)
}

func (renderer *Renderer2d) Release() {
	renderer.releaseVao()
	if renderer.shaderProgram != nil {
		renderer.releaseShader()
	}
}

func (renderer *Renderer2d) releaseShader() {
	renderer.shaderProgram.Delete()
}

func (renderer *Renderer2d) releaseVao() {
	if renderer.vao != 0 {
		log.Printf("2D vertex array destruction")
		gl.DeleteBuffers(1, &renderer.vbo)
		renderer.vbo = 0
		gl.DeleteVertexArrays(1, &renderer.vao)
		renderer.vao = 0
	}
}

// Begin met en place / initialise la vue 2D pour le rendu
func (renderer *Renderer2d) Begin() {
	// Test de profondeur : pas en 2D
	renderer.keepPreviousStateDepthTest = gl.IsEnabled(gl.DEPTH_TEST)
	if renderer.keepPreviousStateDepthTest {
		gl.Disable(gl.DEPTH_TEST)
	}

	// Activer la transparence
	renderer.keepPreviousStateBlend = gl.IsEnabled(gl.BLEND)
	if !renderer.keepPreviousStateBlend {
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}

	renderer.shaderProgram.Use()

	projection2d := mgl32.Ortho2D(0., 800., 600., 0.)
	//projection2d := mgl32.Ortho(0.0, 800.0, 600.0, 0.0, -1.0, 1.0)
	renderer.shaderProgram.UniformMatrix4fv("projection", &projection2d)
}

// End termine le rendu de la scène 2D
func (renderer *Renderer2d) End() {
	renderer.shaderProgram.Unuse()
	if renderer.keepPreviousStateDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	}
	if !renderer.keepPreviousStateBlend {
		gl.Disable(gl.BLEND)
	}
}

func (renderer *Renderer2d) DrawSprite(texture *assets.Texture /*, glm::vec2 position, glm::vec2 size, float rotate, glm::vec3 color*/) {
	// Activer la texture sur l'unité de texture 0 et l'indiquer au Shader
	texture.Bind()
	renderer.shaderProgram.Uniform1i("image", 0)

	// Positionner la couleur
	renderer.shaderProgram.Uniform3f("spriteColor", 1., 0.5, 1.)

	// Positionner le modèle
	//glm::mat4 model = glm::mat4(1.0f);
	//model = glm::translate(model, glm::vec3(position, 0.0f));
	model := mgl32.Translate3D(100., 100., 0.)

	// Rotation
	//model = glm::translate(model, glm::vec3(0.5f * size.x, 0.5f * size.y, 0.0f)); // move origin of rotation to center of quad
	//model = glm::rotate(model, glm::radians(rotate), glm::vec3(0.0f, 0.0f, 1.0f)); // then rotate
	//model = glm::translate(model, glm::vec3(-0.5f * size.x, -0.5f * size.y, 0.0f)); // move origin back

	// Dimensionner le sprite
	model = model.Mul4(mgl32.Scale3D(100., 100., 1.)) // last scale

	// Et passer la transformation au shader
	renderer.shaderProgram.UniformMatrix4fv("coordinates", &model)

	// Faire le rendu du VAO
	gl.BindVertexArray(renderer.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
}
