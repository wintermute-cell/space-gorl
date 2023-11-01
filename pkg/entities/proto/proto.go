package proto

import (
	"cowboy-gorl/pkg/ai"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Entity is an interface that every entity in the game should implement
type Entity interface {
	Init()
	Deinit()
	Update()
    FixedUpdate()
	Draw()
    DrawGUI()
	AddChild(child Entity)
	RemoveChild(child Entity)
	GetChildren() []Entity
	GetParent() Entity
	SetParent(parent Entity)
	GetDrawIndex() int32
    GetName() string
}

// -------------------
//
//	ENTITY 2D
//
// -------------------
type Transform2D struct {
	Position rl.Vector2
	Rotation float32
	Scale    rl.Vector2
}

type Entity2D interface {
	Entity
	GetPosition() rl.Vector2
	SetPosition(new_position rl.Vector2)
	GetScale() rl.Vector2
	SetScale(new_size rl.Vector2)
	GetRotation() float32
}

type Entity2DPlayer interface {
    Entity2D
    SendMessage(message string, sender Entity)
}

// -------------------
//
//	ENTITY 2D AI
//
// -------------------
type Entity2DAI interface {
	Entity2D
	ai.AiControllable
}
