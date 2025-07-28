package input

import (
	"log"
	"otto/system"

	"github.com/anthdm/hollywood/actor"
)

// InputActor handles input processing during tick events
type InputActor struct {
	contexts      map[*actor.PID][]Context
	inputStates   map[*actor.PID][]bool // Track if each context is currently "pressed"
	inputProvider InputProvider
}

var _ actor.Receiver = (*InputActor)(nil)

func New() actor.Producer {
	return func() actor.Receiver {
		return &InputActor{
			contexts:      make(map[*actor.PID][]Context),
			inputStates:   make(map[*actor.PID][]bool),
			inputProvider: NewImGuiProvider(),
		}
	}
}

func (ia *InputActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case system.TickInput:
		ia.processAllInput(ctx, msg.DeltaTime)
	case EventRegisterInputs:
		for _, context := range msg.Contexts {
			ia.contexts[context.GetPID()] = append(ia.contexts[context.GetPID()], context)
			ia.inputStates[context.GetPID()] = append(ia.inputStates[context.GetPID()], false)
		}
	}
}

// processAllInput processes all registered input contexts and sends events
func (ia *InputActor) processAllInput(ctx *actor.Context, deltaTime float64) {
	// Update input state from the provider
	if err := ia.inputProvider.Update(); err != nil {
		log.Printf("Input provider update error: %v", err)
		return
	}

	// Get the current input state
	inputState := ia.inputProvider.GetInputState()

	for pid, contexts := range ia.contexts {
		for idx, context := range contexts {
			// Process the context to get current state, passing capture information
			hasInput := context.Process(deltaTime, inputState, inputState.WantCaptureKeyboard(), inputState.WantCaptureMouse())

			// Get previous state
			states, exists := ia.inputStates[pid]
			if !exists {
				states = make([]bool, len(contexts))
			}

			wasPressed := states[idx]

			// Send event if there's input OR if the context is still active (for continuous input)
			// This ensures smooth camera movement even when mouse delta is zero
			if hasInput || (wasPressed && !hasInput) {
				event := EventInput{
					Context: context,
				}
				ctx.Send(pid, event)
			}

			// Update state - store the current input state, not whether it changed
			ia.inputStates[pid][idx] = hasInput
		}
	}
}
