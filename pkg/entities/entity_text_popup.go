package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TextPopup Entity
type TextPopupEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    text string
    lifetimeTimer *util.Timer
    isFadingOut bool
    fadeAlpha float32
}

func NewTextPopupEntity2D(text string, position rl.Vector2, rotation float32, scale rl.Vector2, durationSeconds float32) *TextPopupEntity2D {
	new_ent := &TextPopupEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: rotation, Scale: scale}},

        text: text,
        lifetimeTimer: util.NewTimer(durationSeconds),
        fadeAlpha: 1.0,
	}
	return new_ent
}

func (ent *TextPopupEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *TextPopupEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *TextPopupEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    if ent.lifetimeTimer.Check() {
        ent.isFadingOut = true
    }

    // fade out text
    if ent.isFadingOut {
        ent.fadeAlpha -= rl.GetFrameTime()
        ent.fadeAlpha = util.Max(ent.fadeAlpha, 0)

        // remove entity
        if ent.fadeAlpha <= 0.01 {
            gem.RemoveEntity(ent)
        }
    }
}

func (ent *TextPopupEntity2D) Draw() {
    rl.DrawText(ent.text, int32(ent.GetPosition().X), int32(ent.GetPosition().Y), 8, rl.NewColor(255, 255, 255, uint8(255.0*ent.fadeAlpha)))
}
