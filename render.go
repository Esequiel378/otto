package main

import (
	"github.com/anthdm/hollywood/actor"
)

type Renderer interface {
	Render()
}

type Render struct {
	entities map[*actor.PID]EntityRigidBody
}

var _ actor.Receiver = (*Render)(nil)

func NewRender() actor.Producer {
	return func() actor.Receiver {
		return &Render{}
	}
}

func (r *Render) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		c.Engine().Subscribe(c.PID())
		r.entities = make(map[*actor.PID]EntityRigidBody)
	case EventEntityInitialized:
		r.entities[msg.PID] = msg.EntityRigidBody
	case EventEntityRenderUpdate:
		r.entities[msg.PID] = msg.EntityRigidBody
	case RequestEntities:
		entities := make([]EntityRigidBody, 0, len(r.entities))
		for pid := range r.entities {
			entities = append(entities, r.entities[pid])
		}
		c.Respond(EntitiesResponse{Entities: entities})
	}
}
