package input

import (
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Event Type
type EventType int32

const (
	EventTypeDown EventType = iota
	EventTypePressed
	EventTypeReleased
)

// Actions
type Action int32

const (
	ActionMoveUp Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
    ActionClickDown
    ActionClickHeld
    ActionClickUp
	ActionReload
	ActionDash
	// Add other actions as needed
)

// Trigger Type
type TriggerType int32

const (
	TriggerTypeKey TriggerType = iota
	TriggerTypeMouse
	TriggerTypeGamepad
)

// Trigger definition
type Trigger struct {
	Type          TriggerType
	Event         EventType
	Key           int32
	MouseButton   int32
	GamepadButton int32
}

var ActionMap = map[Action][]Trigger{
	ActionMoveUp: {
		{Type: TriggerTypeKey, Event: EventTypeDown, Key: rl.KeyW},
	},
	ActionMoveDown: {
		{Type: TriggerTypeKey, Event: EventTypeDown, Key: rl.KeyS},
	},
	ActionMoveLeft: {
		{Type: TriggerTypeKey, Event: EventTypeDown, Key: rl.KeyA},
	},
	ActionMoveRight: {
		{Type: TriggerTypeKey, Event: EventTypeDown, Key: rl.KeyD},
	},
	ActionClickDown: {
		{Type: TriggerTypeMouse, Event: EventTypePressed, MouseButton: rl.MouseLeftButton},
	},
	ActionClickHeld: {
		{Type: TriggerTypeMouse, Event: EventTypeDown, MouseButton: rl.MouseLeftButton},
	},
	ActionClickUp: {
		{Type: TriggerTypeMouse, Event: EventTypeReleased, MouseButton: rl.MouseLeftButton},
	},
    ActionReload: {
        {Type: TriggerTypeKey, Event: EventTypePressed, Key: rl.KeyR},
    },
    ActionDash: {
        {Type: TriggerTypeKey, Event: EventTypePressed, Key: rl.KeySpace},
    },
	// Add other action-trigger mappings
}

// UpdateActionTrigger allows customizing the trigger for a specific action at
// runtime
func UpdateActionTrigger(action Action, triggers []Trigger) {
	ActionMap[action] = triggers
}

// Triggered returns whether the given action has been triggered
func Triggered(action Action) bool {
	for _, trigger := range ActionMap[action] {
		switch trigger.Type {
		case TriggerTypeKey:
			switch trigger.Event {
			case EventTypeDown:
				if rl.IsKeyDown(trigger.Key) {
					return true
				}
			case EventTypePressed:
				if rl.IsKeyPressed(trigger.Key) {
					return true
				}
			case EventTypeReleased:
				if rl.IsKeyReleased(trigger.Key) {
					return true
				}
			}
		case TriggerTypeMouse:
			switch trigger.Event {
			case EventTypeDown:
				if rl.IsMouseButtonDown(trigger.MouseButton) {
					return true
				}
			case EventTypePressed:
				if rl.IsMouseButtonPressed(trigger.MouseButton) {
					return true
				}
			case EventTypeReleased:
				if rl.IsMouseButtonReleased(trigger.MouseButton) {
					return true
				}
			}
		case TriggerTypeGamepad:
			// Implement the checks for gamepad buttons using a similar pattern
		}
	}
	return false
}

// GetMovementVector returns a vector representing the current normalized
// movement direction
func GetMovementVector() rl.Vector2 {
	dir := rl.Vector2Zero()

	if Triggered(ActionMoveUp) {
		dir.Y -= 1
	}
	if Triggered(ActionMoveDown) {
		dir.Y += 1
	}
	if Triggered(ActionMoveLeft) {
		dir.X -= 1
	}
	if Triggered(ActionMoveRight) {
		dir.X += 1
	}

	return util.Vector2NormalizeSafe(dir)
}

// GetCursorPosition returns the current cursor position in screen space
func GetCursorPosition() rl.Vector2 {
	// TODO: add gamepad stuff here. position for gamepad has to be tracked over time.
	return rl.GetMousePosition()
}

// Returns true if the the action was triggered while the mouse cursor was
// inside the given area, in world space coordinates.
func TriggeredInArea(action Action, area rl.Rectangle) bool {
    if rl.CheckCollisionPointRec(render.ScreenToWorldPoint(GetCursorPosition()), area) {
        if Triggered(action) {
            return true
        }
    }
    return false
}
