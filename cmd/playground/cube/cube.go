package cube

import (
	"otto"

	"github.com/anthdm/hollywood/actor"
)

type Cube struct {
	*otto.Entity
}

var _ actor.Receiver = (*Cube)(nil)

func NewCube(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Cube{Entity: otto.NewEntity(physicsPID, rendererPID, inputPID)}
	}
}
