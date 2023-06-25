package graphic

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Camera représente une caméra pour une scène 3D.
type Camera struct {
	// position représente la position de la caméra dans la scène
	Position mgl32.Vec3
	// yaw représente l'angle gauche / droite de la caméra (rotation par rapport à l'axe des y).
	Yaw Angle
	// pitch représente l'angle haut / bas de la caméra (rotation par rapport à l'axe des x).
	Pitch Angle
	// roll représente l'inclinaison de la caméra (rotation par rapport à l'axe des z).
	Roll Angle
	// previousYaw indique l'ancienne valeur de Yaw
	yawPrevious Angle
	yawCos      float32
	yawSin      float32
}

// CreateCamera retourne une caméra
func CreateCamera() Camera {
	camera := Camera{
		Position: mgl32.Vec3{0., 0., 0.},
		Yaw:      0.,
		Pitch:    0.,
		Roll:     0.,
	}
	camera.computeYawCosSin()
	return camera
}

// GetTransform returns the matrix to transform from world coordinates to this camera's coordinates.
func (camera *Camera) GetTransform() mgl32.Mat4 {
	mat := mgl32.Ident4()
	// Positionnement de la caméra
	if camera.Yaw != 0. {
		mat = mgl32.HomogRotate3D(-float32(camera.Yaw), mgl32.Vec3{0.0, 1.0, 0.0})
	}
	if camera.Pitch != 0. || camera.Roll != 0. {
		if camera.yawPrevious != camera.Yaw {
			camera.computeYawCosSin()
		}
		cos := camera.yawCos
		sin := camera.yawSin
		if camera.Pitch != 0. {
			// La rotation verticale est en fonction de l'axe horizontale
			mat = mat.Mul4(mgl32.HomogRotate3D(-float32(camera.Pitch), mgl32.Vec3{-cos, 0.0, sin}))
		}
		if camera.Roll != 0. {
			// La rotation sur soi est dépendant de la direction
			mat = mat.Mul4(mgl32.HomogRotate3D(-float32(camera.Roll), mgl32.Vec3{sin, 0.0, cos}))
		}
	}
	mat = mat.Mul4(mgl32.Translate3D(-camera.Position.X(), -camera.Position.Y(), -camera.Position.Z()))
	return mat
}

func (camera *Camera) MoveForward(distance float32) {
	if camera.yawPrevious != camera.Yaw {
		camera.computeYawCosSin()
	}
	camera.Position = mgl32.Vec3{
		camera.Position.X() - camera.yawSin*distance,
		camera.Position.Y(),
		camera.Position.Z() - camera.yawCos*distance,
	}
}
func (camera *Camera) MoveBack(distance float32) {
	camera.MoveForward(-distance)
}
func (camera *Camera) MoveLeft(distance float32) {
	if camera.yawPrevious != camera.Yaw {
		camera.computeYawCosSin()
	}
	camera.Position = mgl32.Vec3{
		camera.Position.X() - camera.yawCos*distance,
		camera.Position.Y(),
		camera.Position.Z() + camera.yawSin*distance,
	}
}
func (camera *Camera) MoveRight(distance float32) {
	camera.MoveLeft(-distance)
}

func (camera *Camera) computeYawCosSin() {
	camera.yawPrevious = camera.Yaw
	camera.yawCos = float32(camera.Yaw.Cos())
	camera.yawSin = float32(camera.Yaw.Sin())
}
