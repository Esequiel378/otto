package receiver

import (
	"otto/manager"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// InputActor handles input processing during tick events
type InputActor struct {
	contexts    map[string]manager.InputContext
	inputStates map[string]bool // Track if each context is currently "pressed"
}

var _ actor.Receiver = (*InputActor)(nil)

func NewInputActor() actor.Producer {
	return func() actor.Receiver {
		return &InputActor{
			contexts:    make(map[string]manager.InputContext),
			inputStates: make(map[string]bool),
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

	for contextType, context := range ia.contexts {
		// Process the context to get current state
		hasInput := context.Process()

		// Get previous state
		wasPressed, exists := ia.inputStates[contextType]
		if !exists {
			wasPressed = false
		}

		// Check if state has changed
		stateChanged := hasInput != wasPressed

		// For continuous input contexts (like camera), always broadcast when there's input
		// For discrete input contexts (like movement), only broadcast on state changes
		shouldBroadcast := stateChanged

		// Special handling for camera input - always broadcast when there's movement
		if _, ok := context.(*manager.InputCameraControl); ok {
			camera := context.(*manager.InputCameraControl)
			shouldBroadcast = hasInput || (camera.Rotation != (mgl64.Vec2{}) || camera.Zoom != 0.0)
		}

		// Broadcast if needed
		if shouldBroadcast {
			event := manager.InputEvent{
				Context: context,
			}
			c.Engine().BroadcastEvent(event)
		}

		// Update state
		ia.inputStates[contextType] = hasInput
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
