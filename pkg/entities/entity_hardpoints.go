package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/input"
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"cowboy-gorl/pkg/weapons"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Hardpoints Entity
type HardpointsEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    turrets []*weapons.Turret
}

func NewHardpointsEntity2D() *HardpointsEntity2D {
	new_ent := &HardpointsEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0, Scale: rl.Vector2One()}},

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *HardpointsEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

    ent.turrets = append(ent.turrets, weapons.NewTurret(
        ent.GetPosition(),
        weapons.NewProjectileConfig(
            230, // projectile speed
            3, // projectile size
            200, // projectile max range
            "audio/sounds/weapons/chaingun1.ogg",
            rl.LoadTexture("sprites/projectiles/ballistic1.png"),
            weapons.ProjectileTypeBallistic,
            false, // is continuous beam
            ),
        4.0,
        0.1))
}

func (ent *HardpointsEntity2D) Deinit() {
	// Required de-initialization
    ent.BaseEntity2D.Deinit()

    // De-initialization logic for the entity
    // ...
}

func (ent *HardpointsEntity2D) FixedUpdate() {
    for _, trrt := range ent.turrets {
        trrt.FixedUpdate(ent.GetPosition(), input.GetCursorPosition())
    }

    hits := physics.Circlecast(ent.GetPosition(), 100, physics.CollisionCategoryAll)
    logging.Debug("%v", len(hits))
    //for _, hit := range hits {
    //    logging.Debug("%v", hit)
    //}
}
func (ent *HardpointsEntity2D) Update() {
    // Required update
    ent.BaseEntity2D.Update()

    if input.Triggered(input.ActionClickHeld) {
        for _, trrt := range ent.turrets {
            trrt.Shoot(input.GetCursorPosition())
        }
    }

    if input.Triggered(input.ActionClickUp) {
        for _, trrt := range ent.turrets {
            trrt.StopShooting()
        }
    }
}

func (ent *HardpointsEntity2D) Draw() {
    // HARDPOINTS
    // Get the current time
    rl.DrawCircleV(ent.GetPosition(), 100, rl.NewColor(255, 0, 0, 60))

    currentTime := rl.GetTime()

    // Generate hardpoints
    hardpoints := util.DistributePointsOnCircleSection(195, 345, 5, 30, ent.GetPosition())

    for i, point := range hardpoints {
        // Calculate a unique phase offset for each point
        phaseOffset := float64(i) * math.Pi / 2.0 // Half Pi offset for each point

        // Calculate the vertical bounce amount (adjust amplitude and frequency as desired)
        bounceAmplitude := float32(3) // Max vertical movement amount
        bounceFrequency := float64(2) // Speed of the bounce
        bounce := bounceAmplitude * float32(math.Sin(currentTime*bounceFrequency + phaseOffset))

        // Apply the bounce to the y-coordinate
        bouncingPoint := rl.Vector2{X: point.X, Y: point.Y + bounce}

        // Draw the point
        rl.DrawCircleV(bouncingPoint, 2, rl.Purple)
    }

    for _, trrt := range ent.turrets {
        trrt.Draw()
    }
}
