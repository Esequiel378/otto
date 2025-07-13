package renderer

import (
	"otto/system"
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
)

type Renderer interface {
	Render()
}

type Render struct {
	camera   system.Camera
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
		r.entities = make(map[*actor.PID]physics.EntityRigidBody)
		r.camera = system.Camera{}
	case EventEntityRegister:
		r.entities[msg.PID] = msg.EntityRigidBody
	case EventEntityRenderUpdate:
		r.entities[msg.PID] = msg.EntityRigidBody
	case EventUpdateCamera:
		r.camera = msg.Camera
	case RequestEntities:
		entities := make([]physics.EntityRigidBody, 0, len(r.entities))
		for pid := range r.entities {
			entities = append(entities, r.entities[pid])
		}
		ctx.Respond(EntitiesResponse{Entities: entities, Camera: r.camera})
	}
}
