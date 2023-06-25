package graphic

import "github.com/go-gl/mathgl/mgl32"

// Rectangle is a 2D rectangle.
type Rectangle mgl32.Vec4

// BuildRectangle function builds and returns a rectangle.
func BuildRectangle(x float32, y float32, width float32, height float32) Rectangle {
	return Rectangle{x, y, width, height}
}

// BuildRectFromPosAndDim function builds and returns a rectangle.
func BuildRectFromPosAndDim(pos mgl32.Vec2, dim mgl32.Vec2) Rectangle {
	return Rectangle{pos.X(), pos.Y(), dim.X(), dim.Y()}
}

func (rect *Rectangle) X() float32 {
	return rect[0]
}
func (rect *Rectangle) SetX(x float32) {
	rect[0] = x
}
func (rect *Rectangle) Y() float32 {
	return rect[1]
}
func (rect *Rectangle) SetY(y float32) {
	rect[1] = y
}
func (rect *Rectangle) Width() float32 {
	return rect[2]
}
func (rect *Rectangle) SetWidth(width float32) {
	rect[2] = width
}
func (rect *Rectangle) Height() float32 {
	return rect[3]
}
func (rect *Rectangle) SetHeight(height float32) {
	rect[3] = height
}
func (rect *Rectangle) Pos() mgl32.Vec2 {
	return mgl32.Vec2(rect[0:2])
}
func (rect *Rectangle) SetPos(pos mgl32.Vec2) {
	rect[0] = pos.X()
	rect[1] = pos.Y()
}
func (rect *Rectangle) Dim() mgl32.Vec2 {
	return mgl32.Vec2(rect[2:4])
}
func (rect *Rectangle) SetDim(dim mgl32.Vec2) {
	rect[2] = dim.X()
	rect[3] = dim.Y()
}
