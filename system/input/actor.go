package input

import (
	"log"
	"otto/system"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/anthdm/hollywood/actor"
)

// InputActor handles input processing during tick events
type InputActor struct {
	contexts    map[*actor.PID][]Context
	inputStates map[*actor.PID][]bool // Track if each context is currently "pressed"
}

var _ actor.Receiver = (*InputActor)(nil)

func New() actor.Producer {
	return func() actor.Receiver {
		return &InputActor{
			contexts:    make(map[*actor.PID][]Context),
			inputStates: make(map[*actor.PID][]bool),
		}
	}
}

func (ia *InputActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		// Subscribe to tick events when the actor is initialized
		ctx.Engine().Subscribe(ctx.PID())
	case system.Tick:
		// log.Printf("processing %d input contexts", len(ia.contexts))
		ia.processAllInput(ctx)
	case EventRegisterInputs:
		log.Printf("registering %d input contexts", len(msg.Contexts))
		for _, context := range msg.Contexts {
			ia.contexts[context.GetPID()] = append(ia.contexts[context.GetPID()], context)
			ia.inputStates[context.GetPID()] = append(ia.inputStates[context.GetPID()], false)
		}
	}
}

// processAllInput processes all registered input contexts and sends events
func (ia *InputActor) processAllInput(ctx *actor.Context) {
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

	for pid, contexts := range ia.contexts {
		for idx, context := range contexts {
			// Process the context to get current state
			hasInput := context.Process()

			// Get previous state
			states, exists := ia.inputStates[pid]
			if !exists {
				states = make([]bool, len(contexts))
			}

			wasPressed := states[idx]

			// Check if state has changed
			stateChanged := hasInput != wasPressed

			// For continuous input contexts (like camera), always broadcast when there's input
			// For discrete input contexts (like movement), only broadcast on state changes
			shouldBroadcast := stateChanged

			// Broadcast if needed
			if shouldBroadcast {
				event := EventInput{
					Context: context,
				}
				log.Printf("broadcasting input event to %v: %v", pid, event)
				ctx.Send(pid, event)
			}

			// Update state
			ia.inputStates[pid][idx] = stateChanged
		}
	}
}
