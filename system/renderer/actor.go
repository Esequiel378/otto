package renderer

import (
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
)

type Renderer interface {
	Render()
}

type Render struct {
	entities map[*actor.PID]physics.EntityRigidBody
}

var _ actor.Receiver = (*Render)(nil)

func New() actor.Producer {
	return func() actor.Receiver {
		return &Render{}
	}
}

func (r *Render) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		ctx.Engine().Subscribe(ctx.PID())
		r.entities = make(map[*actor.PID]physics.EntityRigidBody)
	case EventEntityRegister:
		r.entities[msg.PID] = msg.EntityRigidBody
	case EventEntityRenderUpdate:
		r.entities[msg.PID] = msg.EntityRigidBody
	case RequestEntities:
		entities := make([]physics.EntityRigidBody, 0, len(r.entities))
		for pid := range r.entities {
			entities = append(entities, r.entities[pid])
		}
		ctx.Respond(EntitiesResponse{Entities: entities})
	}
}
