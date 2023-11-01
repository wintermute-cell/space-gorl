package entities

import (
	"cowboy-gorl/pkg/entities/proto"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Snowpile Entity
type SnowpileEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite rl.Texture2D
    gameState *GameStateHandlerEntity
}

func NewSnowpileEntity2D(position rl.Vector2, gameState *GameStateHandlerEntity) *SnowpileEntity2D {
	new_ent := &SnowpileEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

        sprite: rl.LoadTexture("sprites/snowpile/snowpile.png"),
        gameState: gameState,
	}
	return new_ent
}

func (ent *SnowpileEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *SnowpileEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *SnowpileEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    mp := rl.GetMousePosition()
    poly := []rl.Vector2{
        {X: 0,  Y: 131 + 16},
        {X: 37, Y: 91 + 16},
        {X: 63, Y: 80  + 16},
        {X: 127, Y: 111 + 16},
        {X: 127, Y: 153 + 16},
        {X: 0, Y: 153 + 16},
    }
    if rl.CheckCollisionPointPoly(mp, poly, 6) {
        ent.gameState.MouseHoversSnowpile = true
    } else {
        ent.gameState.MouseHoversSnowpile = false
    }

	// Update logic for the entity
	// ...
}

func (ent *SnowpileEntity2D) Draw() {
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewRectangle(ent.Transform.Position.X, ent.Transform.Position.Y, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewVector2(0, 0),
        ent.Transform.Rotation,
        rl.White,
        )
}
