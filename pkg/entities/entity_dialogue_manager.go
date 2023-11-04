package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
)

// DialogueManager Entity
type DialogueManagerEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
    gameState *GameStateHandlerEntity
}

func NewDialogueManagerEntity(gameState *GameStateHandlerEntity) *DialogueManagerEntity {
	new_ent := &DialogueManagerEntity{
        gameState: gameState,
    }
	return new_ent
}

func (ent *DialogueManagerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

    dialoguePanel := NewDialoguePanelEntity2D()
    gem.AddEntity(ent, dialoguePanel)

    ent.gameState.SetCurrentPlayState(PlayStateDialog)
}

func (ent *DialogueManagerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *DialogueManagerEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *DialogueManagerEntity) Draw() {
	// Draw logic for the entity
	// ...
}
