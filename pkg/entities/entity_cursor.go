package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/input"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Cursor Entity
type CursorEntity2D struct {
	// Required fields
	proto.BaseEntity2D

    sprite rl.Texture2D
	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewCursorEntity2D() *CursorEntity2D {
	new_ent := &CursorEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0.0, Scale: rl.Vector2One()}},

        sprite: rl.LoadTexture("sprites/cursor.png"),
	}
	return new_ent
}

func (ent *CursorEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *CursorEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *CursorEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    ent.SetPosition(input.GetCursorPosition())
}

func (ent *CursorEntity2D) DrawGUI() {
    p := ent.GetPosition()
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewRectangle(float32(int32(p.X)), float32(int32(p.Y)), float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewVector2(3, 3),
        0, rl.White)
}
