package proto

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Entity2D = (*BaseEntity2D)(nil)

// Base Entity
type BaseEntity2D struct {
	// Required fields
	BaseEntity
	entity_manager *EntityManager
	Transform      Transform2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func (ent *BaseEntity2D) GetPosition() rl.Vector2 {
	p := ent.Transform.Position
	// if the parent is an entity2D too, use its position as a base
	switch e := ent.BaseEntity.GetParent().(type) {
	case Entity2D:
		p = rl.Vector2Add(p, e.GetPosition())
	}
	return p
}

func (ent *BaseEntity2D) SetPosition(new_position rl.Vector2) {
	// if the parent is an entity2D too, subtract its position to make new_position relative
	switch e := ent.BaseEntity.GetParent().(type) {
	case Entity2D:
		new_position = rl.Vector2Subtract(new_position, e.GetPosition())
	}
	ent.Transform.Position = new_position
}

func (ent *BaseEntity2D) GetScale() rl.Vector2 {
	return ent.Transform.Scale
}

func (ent *BaseEntity2D) SetScale(new_scale rl.Vector2) {
	ent.Transform.Scale = new_scale
}

func (ent *BaseEntity2D) GetRotation() float32 {
	return ent.Transform.Rotation
}

func (ent *BaseEntity2D) SetRotation(new_rotation float32) {
	ent.Transform.Rotation = new_rotation
}

func (ent *BaseEntity2D) GetTransform() *Transform2D {
	return &ent.Transform
}
