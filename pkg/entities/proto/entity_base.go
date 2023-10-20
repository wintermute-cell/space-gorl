package proto

import (
	"cowboy-gorl/pkg/util"
)

// This checks at compile time if the interface is implemented
var _ Entity = (*BaseEntity)(nil)

// Base Entity
type BaseEntity struct {
	// Required fields
	entity_manager *EntityManager
	children       []Entity
	parent         Entity
    Name string

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func (ent *BaseEntity) Init() {
	// Required initialization
	ent.entity_manager = NewEntityManager()

	// Initialization logic for the entity
	// ...
}

func (ent *BaseEntity) Deinit() {
	// Required de-initialization
	ent.entity_manager.DisableAllEntities()

	// De-initialization logic for the entity
	// ...
}

func (ent *BaseEntity) Update() {
	// Required update
	ent.entity_manager.UpdateEntities()

	// Update logic for the entity
	// ...
}

func (ent *BaseEntity) FixedUpdate() {

}

func (ent *BaseEntity) Draw() {
	// Draw logic for the entity
	// ...
}

// AddChild adds a child to this entity.
func (ent *BaseEntity) AddChild(child Entity) {
	ent.children = append(ent.children, child)
}

// RemoveChild removes a child from this entity.
func (ent *BaseEntity) RemoveChild(child Entity) {
	idx := util.SliceIndex(ent.children, child)
	if idx > -1 {
		ent.children = util.SliceDelete(ent.children, idx, idx+1)
	}
}

// GetChildren returns the children of this entity.
func (ent *BaseEntity) GetChildren() []Entity {
	return ent.children
}

// GetParent sets the parent of this entity.
func (ent *BaseEntity) GetParent() Entity {
	return ent.parent
}

// SetParent sets the parent of this entity.
func (ent *BaseEntity) SetParent(parent Entity) {
	ent.parent = parent
}

// GetDrawIndex returns the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *BaseEntity) GetDrawIndex() int32 {
	return 0
}

func (ent *BaseEntity) GetName() string {
    if ent.Name == "" {
        ent.Name = "UnnamedEntity"
    }
    return ent.Name
}
