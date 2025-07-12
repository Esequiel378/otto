package cube

import (
	"otto/receiver"

	"github.com/anthdm/hollywood/actor"
)

type Cube struct {
	*receiver.Entity
}

var _ actor.Receiver = (*Cube)(nil)

func NewCube(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Cube{Entity: receiver.NewEntity(physicsPID, rendererPID, inputPID)}
	}
}
