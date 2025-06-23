package receiver

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Physics struct {
	entities map[*actor.PID]EntityRigidBody
}

var _ actor.Receiver = (*Physics)(nil)

func NewPhysics() actor.Producer {
	return func() actor.Receiver {
		return &Physics{}
	}
}

func (p *Physics) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		c.Engine().Subscribe(c.PID())
		p.entities = make(map[*actor.PID]EntityRigidBody)
	case EventEntityInitialized:
		p.entities[msg.PID] = msg.EntityRigidBody
	case EventEntityUpdate:
		entity, ok := p.entities[msg.PID]
		if !ok {
			return
		}

		entity.Velocity = msg.Velocity
		p.entities[msg.PID] = entity
	case Tick:
		p.Update(c)
	}
}

func (p *Physics) Update(c *actor.Context) {
	for pid, entity := range p.entities {
		p.ApplyForce(c, pid, entity.Velocity)
	}
}

func (p *Physics) ApplyForce(c *actor.Context, pid *actor.PID, force mgl64.Vec3) {
	entity, ok := p.entities[pid]
	if !ok {
		return
	}
	entity.Velocity = entity.Velocity.Add(force)
	c.Send(pid, EventEntityTransform{
		PID:      pid,
		Position: entity.Position.Add(entity.Velocity),
	})
}
