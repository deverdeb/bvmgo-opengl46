package graphic

import "github.com/go-gl/mathgl/mgl32"

// Color représente une couleur sur 32 bits
type Color mgl32.Vec4

var White = CreateColorRVB(1., 1., 1.)
var Red = CreateColorRVB(1., 0., 0.)
var Green = CreateColorRVB(0., 1., 0.)
var Blue = CreateColorRVB(0., 0., 1.)
var Black = CreateColorRVB(0., 0., 0.)

var GrayBright = CreateColorRVB(0.9, 0.9, 0.9)
var GrayLight = CreateColorRVB(0.7, 0.7, 0.7)
var Gray = CreateColorRVB(0.5, 0.5, 0.5)
var GrayDark = CreateColorRVB(0.2, 0.2, 0.2)

func CreateColorRVBA(red float32, green float32, blue float32, alpha float32) Color {
	return Color{red, green, blue, alpha}
}

func CreateColorRVB(red float32, green float32, blue float32) Color {
	return CreateColorRVBA(red, green, blue, 1.)
}

func (color *Color) Red() float32 {
	return color[0]
}
func (color *Color) SetRed(c float32) {
	color[0] = c
}
func (color *Color) Green() float32 {
	return color[1]
}
func (color *Color) SetGreen(c float32) {
	color[1] = c
}
func (color *Color) Blue() float32 {
	return color[2]
}
func (color *Color) SetBlue(c float32) {
	color[2] = c
}
func (color *Color) Alpha() float32 {
	return color[3]
}
func (color *Color) SetAlpha(c float32) {
	color[3] = c
}
