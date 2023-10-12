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
    color = texture(image, coordinates) * spriteColor;
}
