package entities

import "cowboy-gorl/pkg/entities/proto"

// Template Entity
type TemplateEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewTemplateEntity() *TemplateEntity {
	new_ent := &TemplateEntity{}
	return new_ent
}

func (ent *TemplateEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *TemplateEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *TemplateEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *TemplateEntity) Draw() {
	// Draw logic for the entity
	// ...
}
