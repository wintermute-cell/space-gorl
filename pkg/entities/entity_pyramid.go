package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/messaging"
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/util"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Pyramid Entity
type PyramidEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    Hitpoints int32
    Score int32
    dmgReceiver *messaging.Receiver[messaging.DamageFromEnemyMessage]
    scoreReceiver *messaging.Receiver[messaging.GainScoreMessage]

    hardpoints *HardpointsEntity2D
}

func NewPyramidEntity2D() *PyramidEntity2D {
	new_ent := &PyramidEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.NewVector2(160, 160), Rotation: 0.0, Scale: rl.Vector2One()}},

        Hitpoints: 100,
        Score: 0,
        dmgReceiver: messaging.NewReceiver[messaging.DamageFromEnemyMessage](""),
        scoreReceiver: messaging.NewReceiver[messaging.GainScoreMessage](""),
	}
    new_ent.hardpoints = NewHardpointsEntity2D()
    gem.AddEntity(new_ent, new_ent.hardpoints)
	return new_ent
}

func (ent *PyramidEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

}

func (ent *PyramidEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *PyramidEntity2D) FixedUpdate() {
}

func (ent *PyramidEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    //if input.Triggered(input.ActionClickHeld) {
    //    ent.testTurret.Shoot(input.GetCursorPosition())
    //}

    //if input.Triggered(input.ActionClickUp) {
    //    ent.testTurret.StopShooting()
    //}

    for msg, ok := ent.dmgReceiver.GetNextMessage(); ok; msg, ok = ent.dmgReceiver.GetNextMessage() {
        ent.Hitpoints -= msg.Amount
        render.CameraShake(float32(msg.Amount) / 10)
    }

    for msg, ok := ent.scoreReceiver.GetNextMessage(); ok; msg, ok = ent.scoreReceiver.GetNextMessage() {
        pos := rl.NewVector2(10, 40)
        if msg.FromKill {
            pos = msg.KilledPosition
        }
        pos = rl.Vector2Subtract(pos, ent.GetPosition()) // counteract relative positioning
        popup := NewTextPopupEntity2D(fmt.Sprintf("+%v", msg.Amount), pos, 0.0, rl.Vector2One(), 0.3)
        gem.AddEntity(ent, popup)
        ent.Score += msg.Amount
    }
}

func (ent *PyramidEntity2D) DrawGUI() {
    rl.DrawText(fmt.Sprintf("Health: %v", ent.Hitpoints), 10, 10, 8, rl.White)
    rl.DrawText(fmt.Sprintf("Score: %v", ent.Score), 10, 20, 8, rl.White)

}

func (ent *PyramidEntity2D) Draw() {
    pyramidPoints := util.GetRotatedPyramidPoints(rl.GetTime(), 10, 16)

    // Assuming the points are in the order: top, base1, base2, base3, base4
    if len(pyramidPoints) >= 5 {
        top := rl.Vector2Add(pyramidPoints[0], ent.GetPosition())
        base1 := rl.Vector2Add(pyramidPoints[1], ent.GetPosition())
        base2 := rl.Vector2Add(pyramidPoints[2], ent.GetPosition())
        base3 := rl.Vector2Add(pyramidPoints[3], ent.GetPosition())
        base4 := rl.Vector2Add(pyramidPoints[4], ent.GetPosition())

        // Draw lines for the pyramid base
        rl.DrawLineV(base1, base2, rl.SkyBlue)
        rl.DrawLineV(base2, base3, rl.SkyBlue)
        rl.DrawLineV(base3, base4, rl.SkyBlue)
        rl.DrawLineV(base4, base1, rl.SkyBlue)

        // Draw lines from the pyramid top to each base vertex
        rl.DrawLineV(top, base1, rl.SkyBlue)
        rl.DrawLineV(top, base2, rl.SkyBlue)
        rl.DrawLineV(top, base3, rl.SkyBlue)
        rl.DrawLineV(top, base4, rl.SkyBlue)
    }
}
