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
