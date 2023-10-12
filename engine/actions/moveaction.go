package actions

import (
	"github.com/go-gl/mathgl/mgl32"
	"ogl46/engine"
)

type MoveAction struct {
	getPosition     func() mgl32.Vec2
	setPosition     func(mgl32.Vec2)
	computePosition *mgl32.Vec2
	target          mgl32.Vec2
	speed           float32
	end             bool
}

func BuildMoveAction(getPosition func() mgl32.Vec2, setPosition func(mgl32.Vec2), target mgl32.Vec2, speed float32) Action {
	return &MoveAction{
		getPosition:     getPosition,
		setPosition:     setPosition,
		computePosition: nil,
		target:          target,
		speed:           speed,
		end:             false,
	}
}

func (action *MoveAction) Execute(timer *engine.Timer) {
	if action.end {
		return
	}
	if action.computePosition == nil {
		action.computePosition = &mgl32.Vec2{action.getPosition().X(), action.getPosition().Y()}
	}
	timeInS := float32(timer.ElapsedTime()) / 1000
	distance := action.speed * timeInS
	pos := *action.computePosition
	var newPos mgl32.Vec2
	if action.target.Sub(pos).LenSqr() <= distance {
		newPos = action.target
		action.end = true
	} else {
		direction := action.target.Sub(pos).Normalize()
		newPos = mgl32.Vec2{pos.X() + direction.X()*distance, pos.Y() + direction.Y()*distance}
	}
	action.setPosition(newPos)
	action.computePosition = &newPos
}

func (action *MoveAction) IsEnd() bool {
	return action.end
}
