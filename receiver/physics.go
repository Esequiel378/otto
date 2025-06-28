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
		p.UpdatePosition(c, pid, entity)
	}
}

func (p *Physics) UpdatePosition(c *actor.Context, pid *actor.PID, entity EntityRigidBody) {
	// Apply velocity to position with a fixed movement speed
	movementSpeed := 0.1 // Fixed movement speed per frame
	entity.Position = entity.Position.Add(entity.Velocity.Mul(movementSpeed))

	// Don't apply damping when there's active input - let the input system control velocity
	// Only apply damping when there's no input (velocity will be zero)
	if entity.Velocity.Len() < 0.01 {
		entity.Velocity = mgl64.Vec3{}
	}

	p.entities[pid] = entity

	c.Send(pid, EventEntityTransform{
		PID:      pid,
		Position: entity.Position,
	})
}
