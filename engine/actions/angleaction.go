package actions

import (
	"math"
	"ogl46/engine"
	"ogl46/engine/graphic"
)

type AngleAction struct {
	getAngle func() graphic.Angle
	setAngle func(graphic.Angle)
	target   graphic.Angle
	speed    graphic.Angle
	end      bool
}

func BuildAngleAction(getAngle func() graphic.Angle, setAngle func(graphic.Angle), target graphic.Angle, speed graphic.Angle) Action {
	return &AngleAction{
		getAngle: getAngle,
		setAngle: setAngle,
		target:   target,
		speed:    speed,
		end:      false,
	}
}

func (action *AngleAction) Execute(timer *engine.Timer) {
	if action.end {
		return
	}
	timeInS := float64(timer.ElapsedTime()) / 1000
	distance := action.speed.Radian() * timeInS
	currentAngle := action.getAngle().Radian()
	var newAngle graphic.Angle
	if math.Abs(action.target.Radian()-currentAngle) <= math.Abs(distance) {
		newAngle = action.target
		action.end = true
	} else {
		newAngle = graphic.AngleInRadian(currentAngle + distance)
	}
	action.setAngle(newAngle)
}

func (action *AngleAction) IsEnd() bool {
	return action.end
}
