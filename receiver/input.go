package receiver

import (
	"otto/manager"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// InputActor handles input processing during tick events
type InputActor struct {
	contexts map[string]manager.InputContext
}

var _ actor.Receiver = (*InputActor)(nil)

func NewInputActor() actor.Producer {
	return func() actor.Receiver {
		return &InputActor{
			contexts: make(map[string]manager.InputContext),
		}
	}
}

func (ia *InputActor) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		// Subscribe to tick events when the actor is initialized
		c.Engine().Subscribe(c.PID())
	case Tick:
		// Process input during each tick
		ia.processAllInput(c)
	case RegisterInputContext:
		// Register an input context from an entity
		ia.contexts[msg.Context.GetType()] = msg.Context
	case UnregisterInputContext:
		// Unregister an input context
		delete(ia.contexts, msg.ContextType)
	}
}

// processAllInput processes all registered input contexts and sends events
func (ia *InputActor) processAllInput(c *actor.Context) {
	// Wrap everything in panic recovery to handle ImGui context issues
	defer func() {
		if r := recover(); r != nil {
			// ImGui context is invalid or destroyed, silently return
			return
		}
	}()

	// Don't process input if ImGui wants to capture it
	if ctx := imgui.CurrentContext(); ctx == nil {
		return
	}

	io := imgui.CurrentIO()
	if io == nil {
		return
	}

	if io.WantCaptureKeyboard() {
		return
	}

	for _, context := range ia.contexts {
		// Process the context to get current state
		hasInput := context.Process()

		// Always broadcast if there's input
		// For movement contexts, also broadcast when input stops (velocity becomes zero)
		shouldBroadcast := hasInput

		if movement, ok := context.(*InputPlayerMovement); ok {
			// For movement, broadcast if there's input OR if velocity is zero (stopping)
			shouldBroadcast = hasInput || movement.Velocity == (mgl64.Vec3{})
		}

		if shouldBroadcast {
			event := manager.InputEvent{
				Context: context,
			}
			c.Engine().BroadcastEvent(event)
		}
	}
}

// RegisterInputContext is sent by entities to register their input contexts
type RegisterInputContext struct {
	Context manager.InputContext
}

// UnregisterInputContext is sent by entities to unregister their input contexts
type UnregisterInputContext struct {
	ContextType string
}
