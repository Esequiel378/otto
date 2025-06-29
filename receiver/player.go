package receiver

import (
	"otto/manager"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	Entity
	inputPID *actor.PID
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: Entity{
				position:    mgl64.Vec3{0, 0, 0}, // Start at origin
				velocity:    mgl64.Vec3{0, 0, 0},
				scale:       mgl64.Vec3{1, 1, 1}, // Make it visible size
				rotation:    mgl64.Vec3{0, 0, 0},
				physicsPID:  physicsPID,
				rendererPID: rendererPID,
			},
			inputPID: inputPID,
		}
	}
}

func (p *Player) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		// Subscribe to input events
		c.Engine().Subscribe(c.PID())
		// Register our input context with the input actor
		c.Send(p.inputPID, RegisterInputContext{
			Context: &InputPlayerMovement{},
		})
		// TODO: Subscribe to physics events by sending a message to the physics manager
		// TODO: Subscribe to renderer events by sending a message to the renderer manager
	case manager.InputEvent:
		switch ctx := msg.Context.(type) {
		case *InputPlayerMovement:
			c.Send(p.physicsPID, EventEntityUpdate{
				PID:      c.PID(),
				Velocity: ctx.Velocity,
			})
		}
	}

	p.Entity.Receive(c)
}

type InputPlayerMovement struct {
	Velocity mgl64.Vec3
}

var _ manager.InputContext = (*InputPlayerMovement)(nil)

// GetType returns the type identifier for player movement input
func (h *InputPlayerMovement) GetType() string {
	return "player_movement"
}

// Process handles player movement input
func (h *InputPlayerMovement) Process() bool {
	h.Velocity = mgl64.Vec3{} // Reset velocity

	// Process keyboard input - Allow multiple keys to be pressed simultaneously
	if imgui.IsKeyDown(imgui.KeyW) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, 1}) // Forward (Z+)
	}

	if imgui.IsKeyDown(imgui.KeyS) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, -1}) // Backward (Z-)
	}

	if imgui.IsKeyDown(imgui.KeyA) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{1, 0, 0}) // Left (X+)
	}

	if imgui.IsKeyDown(imgui.KeyD) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{-1, 0, 0}) // Right (X-)
	}

	if imgui.IsKeyDown(imgui.KeySpace) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, -1, 0}) // Up (Y-)
	}

	if imgui.IsKeyDown(imgui.KeyLeftShift) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 1, 0}) // Down (Y+)
	}

	// Normalize the velocity to prevent faster diagonal movement
	if h.Velocity.Len() > 0 {
		h.Velocity = h.Velocity.Normalize()
	}

	// Return true if any movement key is pressed (for state tracking)
	// The actual velocity is stored in the context for processing
	return h.Velocity != (mgl64.Vec3{})
}

// Legacy Handle method for backward compatibility
func (h *InputPlayerMovement) Handle() {
	h.Process()
}
