package assets

import (
	"log/slog"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type VertexArray struct {
	// vao (Vertex Array Object)
	vao uint32
	// vbo (Vertex Buffer Object)
	vbo uint32

	mode uint32
}

func NewVertexArrayWithTexture(vertices []float32, mode uint32) *VertexArray {
	slog.Debug("vertex array creation")

	// https://learnopengl.com/code_viewer_gh.php?code=src/1.getting_started/4.2.textures_combined/textures_combined.cpp

	// Créer le VAO (Vertex Array Object)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	// Créer le VBO (Vertex Buffer Object)
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(vao)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW) // 4 = float32 size

	var stride int32 = 3*4 + 2*4 // Taille des données d'un vertex = 3 coordonnées + 2 coordonnées de texture
	var offset uintptr = 0       // Au début, l'offset est au début

	// position
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, offset) // 3 = nb coordonnées
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4 // Décalage de l'offset des 3 coordonnées

	// color
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, offset) // 2 = nb coordonnées
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4 // Décalage de l'offset des 2 coordonnées de texture

	// Unbind du VAO à la fin de la création
	gl.BindVertexArray(0)

	return &VertexArray{
		vao:  vao,
		vbo:  vbo,
		mode: mode,
	}
}

func (vertexArray *VertexArray) Bind() {
	gl.BindVertexArray(vertexArray.vao)
}

func (vertexArray *VertexArray) UnBind() {
	gl.BindVertexArray(0)
}

func (vertexArray *VertexArray) Release() {
	if vertexArray.vao != 0 {
		slog.Debug("vertex array destruction")
		gl.DeleteBuffers(1, &vertexArray.vbo)
		vertexArray.vbo = 0
		gl.DeleteVertexArrays(1, &vertexArray.vao)
		vertexArray.vao = 0
	}
}

func (vertexArray *VertexArray) Draw() {
	vertexArray.Bind()
	//gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.DrawArrays(vertexArray.mode, 0, 36)
}
