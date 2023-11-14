package entities

import (
    "cowboy-gorl/pkg/entities/proto"
    rl "github.com/gen2brain/raylib-go/raylib"
)


type Item struct {
    Id string
    DisplayName string
    sprite rl.Texture2D
}

// ItemHandler Entity
type ItemHandlerEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewItemHandlerEntity() *ItemHandlerEntity {
	new_ent := &ItemHandlerEntity{}
	return new_ent
}

func (ent *ItemHandlerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *ItemHandlerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *ItemHandlerEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *ItemHandlerEntity) Draw() {
	// Draw logic for the entity
	// ...
}
