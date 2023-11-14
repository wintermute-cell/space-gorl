package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Icebar Entity
type IcebarEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    spriteBackground rl.Texture2D
    spriteForeground rl.Texture2D

    gameState *GameStateHandlerEntity
}

func NewIcebarEntity2D(position rl.Vector2, gameState *GameStateHandlerEntity) *IcebarEntity2D {
	new_ent := &IcebarEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

		// Initialize custom fields here
		// ...
        spriteBackground: rl.LoadTexture("sprites/icebar/icebar_bg.png"),
        spriteForeground: rl.LoadTexture("sprites/icebar/icebar_fg.png"),
        gameState: gameState,
	}
	return new_ent
}

func (ent *IcebarEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *IcebarEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *IcebarEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *IcebarEntity2D) Draw() {
	// Draw logic for the entity
	// ...
}

func (ent *IcebarEntity2D) DrawGUI() {
    rl.DrawTexturePro(
        ent.spriteBackground,
        rl.NewRectangle(0, 0, float32(ent.spriteBackground.Width), float32(ent.spriteBackground.Height)),
        rl.NewRectangle(0, 0, float32(ent.spriteBackground.Width), float32(ent.spriteBackground.Height)),
        rl.Vector2Zero(),
        0, rl.White)
    drawnWidth := float32(ent.spriteForeground.Width) * (ent.gameState.IceMeterValue/100)
    rl.DrawTexturePro(
        ent.spriteForeground,
        rl.NewRectangle(0, 0, drawnWidth, float32(ent.spriteBackground.Height)),
        rl.NewRectangle(0, 0, drawnWidth, float32(ent.spriteBackground.Height)),
        rl.Vector2Zero(),
        0, rl.White)
}
