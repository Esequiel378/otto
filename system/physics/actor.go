package physics

import (
	"otto/system"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Physics struct {
	entities map[*actor.PID]EntityRigidBody
}

var _ actor.Receiver = (*Physics)(nil)

func New() actor.Producer {
	return func() actor.Receiver {
		return &Physics{}
	}
}

func (p *Physics) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		ctx.Engine().Subscribe(ctx.PID())
		p.entities = make(map[*actor.PID]EntityRigidBody)
	case EventRigidBodyRegister:
		p.entities[msg.PID] = msg.EntityRigidBody
	case EventRigidBodyUpdate:
		entity, ok := p.entities[msg.PID]
		if !ok {
			return
		}

		entity.Velocity = msg.Velocity
		entity.AngularVelocity = msg.AngularVelocity
		p.entities[msg.PID] = entity
	case system.ServerTick:
		p.Update(ctx)
	}
}

func (p *Physics) Update(ctx *actor.Context) {
	tick, ok := ctx.Message().(system.ServerTick)
	if !ok {
		panic("tick message not found")
	}

	for pid, entity := range p.entities {
		if entity.Velocity.Len() > 0 {
			p.UpdatePosition(ctx, pid, entity, tick.DeltaTime)
		}
	}
}

func (p *Physics) UpdatePosition(ctx *actor.Context, pid *actor.PID, entity EntityRigidBody, deltaTime float64) {
	// Apply velocity to position with frame-rate independent movement speed
	movementSpeed := 10.0 // Increased for faster movement
	entity.Position = entity.Position.Add(entity.Velocity.Mul(movementSpeed * deltaTime))

	// Apply angular velocity to rotation with frame-rate independent rotation speed
	rotationSpeed := 4.0 // Increased for faster camera
	entity.Rotation = entity.Rotation.Add(entity.AngularVelocity.Mul(rotationSpeed * deltaTime))

	// Don't apply damping when there's active input - let the input system control velocity
	// Only apply damping when there's no input (velocity will be zero)
	if entity.Velocity.Len() < 0.01 {
		entity.Velocity = mgl64.Vec3{}
	}

	// Apply damping to angular velocity when there's no rotation input
	if entity.AngularVelocity.Len() < 0.01 {
		entity.AngularVelocity = mgl64.Vec3{}
	}

	p.entities[pid] = entity

	ctx.Send(pid, EventRigidBodyTransform{
		PID:      pid,
		Position: entity.Position,
		Rotation: entity.Rotation,
	})
}
