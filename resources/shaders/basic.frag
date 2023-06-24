#version 410 core

in vec2 TexCoord;

out vec4 color;

uniform sampler2D texture0;
uniform sampler2D texture1;

void main()
{
    // mix the two textures together 
    color = mix(texture(texture0, TexCoord), texture(texture1, TexCoord), 0.5);
}