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

		// Update velocity components, preserving existing components when not provided
		if msg.Velocity.X() != 0 {
			entity.Velocity = mgl64.Vec3{msg.Velocity.X(), entity.Velocity.Y(), entity.Velocity.Z()}
		}
		if msg.Velocity.Y() != 0 {
			entity.Velocity = mgl64.Vec3{entity.Velocity.X(), msg.Velocity.Y(), entity.Velocity.Z()}
		}
		if msg.Velocity.Z() != 0 {
			entity.Velocity = mgl64.Vec3{entity.Velocity.X(), entity.Velocity.Y(), msg.Velocity.Z()}
		}

		// Update angular velocity
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
		// Apply gravity to all entities
		p.ApplyGravity(&entity, tick.DeltaTime)
		p.entities[pid] = entity

		// Always update position if there's any velocity (including from gravity)
		p.UpdatePosition(ctx, pid, entity, tick.DeltaTime)
	}
}

func (p *Physics) ApplyGravity(entity *EntityRigidBody, deltaTime float64) {
	// Apply gravity acceleration (9.8 m/s², scaled for game world)
	gravityAcceleration := -9.8 * 0.2 // Further reduced for gentler fall speed
	gravityVelocity := mgl64.Vec3{0, gravityAcceleration * deltaTime, 0}

	// Add gravity to existing velocity
	entity.Velocity = entity.Velocity.Add(gravityVelocity)
}

func (p *Physics) UpdatePosition(ctx *actor.Context, pid *actor.PID, entity EntityRigidBody, deltaTime float64) {
	// Apply velocity to position with frame-rate independent movement speed
	movementSpeed := 7.0
	newPosition := entity.Position.Add(entity.Velocity.Mul(movementSpeed * deltaTime))

	// Calculate the bottom of the entity based on its scale
	// For cubes and other entities, the bottom is at position.Y - (scale.Y / 2)
	entityBottom := newPosition.Y() - (entity.Scale.Y() / 2)

	// Floor collision detection at y=0
	// Check if the bottom of the entity would go below the floor
	if entityBottom < 0 {
		// Calculate the correct position so the bottom of the entity is at y=0
		correctedY := entity.Scale.Y() / 2
		newPosition = mgl64.Vec3{newPosition.X(), correctedY, newPosition.Z()}

		// Handle collision based on entity type
		if entity.Velocity.Y() < 0 {
			if entity.EntityType == "player" {
				// Player should not bounce - just stop
				entity.Velocity = mgl64.Vec3{entity.Velocity.X(), 0, entity.Velocity.Z()}
			} else {
				// Other entities (cubes) should bounce
				bounceFactor := 0.3 // 30% bounce (70% energy loss)
				entity.Velocity = mgl64.Vec3{entity.Velocity.X(), -entity.Velocity.Y() * bounceFactor, entity.Velocity.Z()}
			}
		}
	}

	entity.Position = newPosition

	// Apply angular velocity to rotation with frame-rate independent rotation speed
	rotationSpeed := 4.0 // Increased for faster camera
	entity.Rotation = entity.Rotation.Add(entity.AngularVelocity.Mul(rotationSpeed * deltaTime))

	// Apply damping to horizontal velocity only when there's no input
	// Allow gravity to continue acting (vertical velocity)
	horizontalVelocity := mgl64.Vec3{entity.Velocity.X(), 0, entity.Velocity.Z()}
	if horizontalVelocity.Len() < 0.01 {
		// Only damp horizontal velocity, preserve vertical velocity from gravity
		entity.Velocity = mgl64.Vec3{0, entity.Velocity.Y(), 0} // Keep only vertical velocity
	} else {
		// When there's horizontal movement, apply gentle damping ONLY to horizontal components
		// NEVER touch the vertical velocity (gravity)
		dampingFactor := 0.95 // 5% damping per frame
		entity.Velocity = mgl64.Vec3{
			entity.Velocity.X() * dampingFactor,
			entity.Velocity.Y(), // Keep gravity completely unaffected
			entity.Velocity.Z() * dampingFactor,
		}
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

	// Send ground state for player entities
	if entity.EntityType == "player" {
		// Check if player is on ground (bottom of entity is at floor level)
		entityBottom := entity.Position.Y() - (entity.Scale.Y() / 2)
		isOnGround := entityBottom <= 0.01 && entity.Velocity.Y() <= 0.01

		ctx.Send(pid, EventGroundState{
			PID:        pid,
			IsOnGround: isOnGround,
		})
	}
}
