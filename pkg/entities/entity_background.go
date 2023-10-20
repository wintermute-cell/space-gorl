package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Background Entity
type BackgroundEntity2D struct {
	// Required fields
	proto.BaseEntity2D

    sprite rl.Texture2D
}

func NewBackgroundEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *BackgroundEntity2D {
	new_ent := &BackgroundEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: rotation, Scale: scale}},

        sprite: rl.LoadTexture(":q"),
	}
	return new_ent
}

func (ent *BackgroundEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *BackgroundEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *BackgroundEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *BackgroundEntity2D) Draw() {
	// Draw logic for the entity
	// ...
}
