package graphic

import (
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// Angle (value in radians)
type Angle float64

// AngleInDegree crée un angle à partie d'une valeur en degrés
func AngleInDegree(valueDegree float64) Angle {
	return Angle(mgl64.DegToRad(valueDegree))
}

// AngleInRadian crée un angle à partie d'une valeur en radians
func AngleInRadian(valueRadian float64) Angle {
	return Angle(valueRadian)
}

// Degree retourne la valeur de l'angle en degrés
func (angle Angle) Degree() float64 {
	return mgl64.RadToDeg(float64(angle))
}

// Radian retourne la valeur de l'angle en radians
func (angle Angle) Radian() float64 {
	return float64(angle)
}

// Cos retourne le cosinus de l'angle
func (angle Angle) Cos() float64 {
	return math.Cos(angle.Radian())
}

// Sin retourne le Sin de l'angle
func (angle Angle) Sin() float64 {
	return math.Sin(angle.Radian())
}

// Mod360 normalize angle (return angle between 0° and 360°).
func (angle Angle) Mod360() Angle {
	if angle.Degree() >= 0 && angle.Degree() < 360 {
		return angle
	} else {
		angleInDeg360 := math.Mod(angle.Degree(), 360.)
		if angleInDeg360 < 0 {
			angleInDeg360 += 360
		}
		return AngleInDegree(angleInDeg360)
	}
}
