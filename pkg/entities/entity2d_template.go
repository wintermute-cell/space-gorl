package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Template Entity
type TemplateEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewTemplateEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *TemplateEntity2D {
	new_ent := &TemplateEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: rotation, Scale: scale}},

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *TemplateEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *TemplateEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *TemplateEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *TemplateEntity2D) Draw() {
	// Draw logic for the entity
	// ...
}
