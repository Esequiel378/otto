package cube

import (
	"otto/system"

	"github.com/anthdm/hollywood/actor"
)

type Cube struct {
	*system.Entity
}

var _ actor.Receiver = (*Cube)(nil)

func NewCube(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Cube{Entity: system.NewEntity(physicsPID, rendererPID, inputPID)}
	}
}
