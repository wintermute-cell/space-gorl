package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Template Entity
type TemplateEntity2DAI struct {
	// Required fields
	proto.BaseEntity2DAI

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewTemplateEntity2DAI(position rl.Vector2, rotation float32, scale rl.Vector2, player_ref proto.Entity2DPlayer) *TemplateEntity2DAI {
	new_ent := &TemplateEntity2DAI{
		BaseEntity2DAI: proto.BaseEntity2DAI{
			BaseEntity2D: proto.BaseEntity2D{
				Transform: proto.Transform2D{Position: position, Rotation: rotation, Scale: scale},
			},
			PlayerEntity: player_ref,
		},

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *TemplateEntity2DAI) Init() {
	// Required initialization
	ent.BaseEntity2DAI.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *TemplateEntity2DAI) Deinit() {
	// Required de-initialization
	ent.BaseEntity2DAI.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *TemplateEntity2DAI) Update() {
	// Required update
	ent.BaseEntity2DAI.Update()

	// Update logic for the entity
	// ...
}

func (ent *TemplateEntity2DAI) Draw() {
	// Draw logic for the entity
	// ...
}
