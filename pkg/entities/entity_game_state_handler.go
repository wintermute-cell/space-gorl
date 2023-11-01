package entities

import "cowboy-gorl/pkg/entities/proto"

type PlayState int32
const (
    PlayStateNone PlayState = iota
    PlayStateDigging
    PlayStateDialog
)

// GameStateHandler Entity
type GameStateHandlerEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
    currentPlayState PlayState

    // Game State Fields
    MouseHoversSnowpile bool
}

func NewGameStateHandlerEntity() *GameStateHandlerEntity {
	new_ent := &GameStateHandlerEntity{
        BaseEntity: proto.BaseEntity {
            Name: "PlayStateHandler",
        },
        currentPlayState: PlayStateDigging,

        MouseHoversSnowpile: false,
    }
	return new_ent
}

func (ent *GameStateHandlerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *GameStateHandlerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *GameStateHandlerEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *GameStateHandlerEntity) Draw() {
	// Draw logic for the entity
	// ...
}
