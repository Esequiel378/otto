package physics

import (
	"otto/system"
	"otto/system/camera"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Physics struct {
	entities     map[*actor.PID]EntityRigidBody
	cameraPID    *actor.PID
	latestCamera *camera.Camera
}

var _ actor.Receiver = (*Physics)(nil)

func New(cameraPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Physics{
			cameraPID: cameraPID,
		}
	}
}

func (p *Physics) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		c.Engine().Subscribe(c.PID())
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
	case EventCameraUpdate:
		p.latestCamera = &msg.Camera
	case system.Tick:
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

	// Calculate movement direction based on camera orientation if available
	var movementDirection mgl64.Vec3
	if p.latestCamera != nil {
		// Use camera vectors to transform velocity into world space
		front := p.latestCamera.GetFrontVector()
		right := p.latestCamera.GetRightVector()
		up := p.latestCamera.GetUpVector()

		// Transform velocity from camera-relative to world coordinates
		movementDirection = right.Mul(entity.Velocity.X()).
			Add(up.Mul(entity.Velocity.Y())).
			Add(front.Mul(entity.Velocity.Z()))
	} else {
		// Fallback to direct velocity if no camera available
		movementDirection = entity.Velocity
	}

	entity.Position = entity.Position.Add(movementDirection.Mul(movementSpeed))

	// Apply angular velocity to rotation
	rotationSpeed := 0.1 // Fixed rotation speed per frame
	entity.Rotation = entity.Rotation.Add(entity.AngularVelocity.Mul(rotationSpeed))

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

	c.Send(pid, EventRigidBodyTransform{
		PID:      pid,
		Position: entity.Position,
		Rotation: entity.Rotation,
	})
}
