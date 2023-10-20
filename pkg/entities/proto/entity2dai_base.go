package proto

import (
	"cowboy-gorl/pkg/ai"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Entity2DAI = (*BaseEntity2DAI)(nil)

// Base Entity
type BaseEntity2DAI struct {
	// Required fields
	BaseEntity2D
	PlayerEntity Entity2DPlayer
	AiTarget     *ai.AiTarget
    AiSteeringForce rl.Vector2

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func (ent *BaseEntity2DAI) SetAiTarget(target *ai.AiTarget) {
	ent.AiTarget = target
}

func (ent *BaseEntity2DAI) SetAiSteeringForce(force rl.Vector2) {
    ent.AiSteeringForce = force
}

func (ent *BaseEntity2DAI) GetPlayerPosition() rl.Vector2 {
	return ent.PlayerEntity.GetPosition()
}

func (ent *BaseEntity2DAI) CanSeePlayer() bool {
	return true
}

func (ent *BaseEntity2DAI) CanSeePoint(point rl.Vector2) bool {
    return true
}
