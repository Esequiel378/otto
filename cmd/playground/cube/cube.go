package cube

import (
	"otto"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Cube struct {
	*otto.Entity
}

var _ actor.Receiver = (*Cube)(nil)

func NewCube(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		entity := otto.NewEntity(physicsPID, rendererPID, inputPID)
		entity.ModelName = "cube"
		entity.Position = mgl64.Vec3{0, 0, 2} // Position the cube in front of the camera
		return &Cube{Entity: entity}
	}
}

func NewCubeWithPosition(physicsPID, rendererPID, inputPID *actor.PID, position mgl64.Vec3) actor.Producer {
	return func() actor.Receiver {
		entity := otto.NewEntity(physicsPID, rendererPID, inputPID)
		entity.ModelName = "cube"
		entity.Position = position
		entity.Scale = mgl64.Vec3{1, 1, 1}
		return &Cube{Entity: entity}
	}
}
