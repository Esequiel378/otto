# Extensible Input System

This document describes how to use the extensible input system that allows developers to define custom input contexts and broadcast them for anyone to receive and process.

## Overview

The input system consists of:
- **InputContext Interface**: Defines the contract for input contexts
- **InputManager**: Manages multiple input contexts and broadcasts events
- **InputEvent**: Represents broadcasted input events
- **Pre-built Input Contexts**: Ready-to-use input contexts for common scenarios

## How It Works

1. **Input Contexts**: Developers define input contexts that implement the `InputContext` interface
2. **Registration**: Input contexts are registered with the `InputManager`
3. **Processing**: The manager processes all contexts each frame and broadcasts events when input is detected
4. **Receiving**: Any actor can listen for `InputEvent` messages and handle them accordingly

## Creating Custom Input Contexts

To create a custom input context, implement the `InputContext` interface:

```go
type MyCustomInput struct {
    // Your input data fields
    Value string
}

// GetType returns a unique identifier for this input context
func (m *MyCustomInput) GetType() string {
    return "my_custom_input"
}

// Process handles the input and returns true if any input was detected
func (m *MyCustomInput) Process() bool {
    m.Value = "" // Reset value
    
    // Check for input (example using ImGui)
    if imgui.IsKeyDown(imgui.KeyX) {
        m.Value = "x_pressed"
        return true
    }
    
    return false
}
```

## Registering Input Contexts

Register your input contexts with the manager:

```go
inputManager := receiver.NewInputManager(engine)

// Register built-in contexts
inputManager.RegisterContext(&receiver.InputPlayerMovement{})
inputManager.RegisterContext(&receiver.InputCameraControl{})

// Register custom contexts
inputManager.RegisterContext(&MyCustomInput{})
```

## Receiving Input Events

Actors can listen for broadcasted input events:

```go
func (a *MyActor) Receive(c *actor.Context) {
    switch msg := c.Message().(type) {
    case receiver.InputEvent:
        if msg.Type == "my_custom_input" && msg.Active {
            if input, ok := msg.Data.(*MyCustomInput); ok {
                // Handle the input
                fmt.Printf("Received input: %s\n", input.Value)
            }
        }
    }
}
```

## Built-in Input Contexts

### InputPlayerMovement
Handles WASD movement with space/shift for vertical movement.

**Keys:**
- W: Forward (Z+)
- S: Backward (Z-)
- A: Left (X-)
- D: Right (X+)
- Space: Up (Y+)
- Left Shift: Down (Y-)

### InputCameraControl
Handles camera rotation and zoom.

**Keys:**
- Right Mouse + I: Pitch up
- Right Mouse + K: Pitch down
- Right Mouse + J: Yaw left
- Right Mouse + L: Yaw right
- =: Zoom in
- -: Zoom out

### InputUIInteraction
Handles UI-related input like mouse position and key presses.

**Keys:**
- Escape: Escape key
- Enter: Enter key
- Tab: Tab key

### InputGameActions
Handles game-specific actions.

**Keys:**
- F: Interact
- R: Reload
- E: Use
- Q: Drop

## Example Usage

Here's a complete example of creating and using a custom input context:

```go
// Define your input context
type InventoryInput struct {
    OpenInventory bool
    NextItem      bool
    PrevItem      bool
}

func (i *InventoryInput) GetType() string {
    return "inventory"
}

func (i *InventoryInput) Process() bool {
    i.OpenInventory = false
    i.NextItem = false
    i.PrevItem = false
    
    if imgui.IsKeyPressed(imgui.KeyI) {
        i.OpenInventory = true
    }
    if imgui.IsKeyPressed(imgui.KeyRight) {
        i.NextItem = true
    }
    if imgui.IsKeyPressed(imgui.KeyLeft) {
        i.PrevItem = true
    }
    
    return i.OpenInventory || i.NextItem || i.PrevItem
}

// Register it
inputManager.RegisterContext(&InventoryInput{})

// Handle it in an actor
func (a *InventoryActor) Receive(c *actor.Context) {
    switch msg := c.Message().(type) {
    case receiver.InputEvent:
        if msg.Type == "inventory" && msg.Active {
            if input, ok := msg.Data.(*InventoryInput); ok {
                if input.OpenInventory {
                    a.ToggleInventory()
                }
                if input.NextItem {
                    a.SelectNextItem()
                }
                if input.PrevItem {
                    a.SelectPrevItem()
                }
            }
        }
    }
}
```

## Benefits

1. **Decoupling**: Input handling is separated from game logic
2. **Extensibility**: Easy to add new input contexts without modifying existing code
3. **Broadcasting**: Multiple actors can respond to the same input
4. **Type Safety**: Input contexts are strongly typed
5. **Reusability**: Input contexts can be reused across different parts of the application 