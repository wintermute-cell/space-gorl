package entities

import (
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/messaging"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Enemy Entity
type EnemyEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    collider *physics.Collider
    data EnemyData
    isAlive bool

    pyramidPosition rl.Vector2
}

func NewEnemyEntity2D(position rl.Vector2, rotation float32, scale rl.Vector2) *EnemyEntity2D {
	new_ent := &EnemyEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: rotation, Scale: scale}},
        isAlive: true,
        pyramidPosition: rl.NewVector2(160, 160),
	}
	return new_ent
}

func (ent *EnemyEntity2D) SetData(data EnemyData) {
    ent.data = data
    ent.isAlive = true
    
    radius := (data.Size.X + data.Size.Y)/4
    callbacks := make(map[physics.CollisionCategory]physics.CollisionCallback)
    callbacks[physics.CollisionCategoryBullet] = func(values ...any) {
        if len(values) != 2 {
            return
        }
        
        if id, ok := values[0].(string); !ok || id != "apply-damage" {
            return
        }

        damage, ok := values[1].(float32)
        if !ok {
            return
        }
        ent.data.Hitpoints -= damage
    }
    ent.collider = physics.NewCircleCollider(ent.GetPosition(), radius, physics.BodyTypeDynamic).SetFixedRotation(false).SetCallbacks(callbacks)
    ent.collider.SetLinearDamping(data.LinearDamping)
    ent.collider.SetAngularDamping(data.AngularDamping)
}

func (ent *EnemyEntity2D) GetData() EnemyData {
    return ent.data
}

func (ent *EnemyEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *EnemyEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *EnemyEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    // match the entity position to the collider
    ent.SetPosition(ent.collider.GetPosition())
    ent.SetRotation(float32(ent.collider.GetB2Body().GetAngle())*rl.Rad2deg + 90)

    dist := rl.Vector2Distance(ent.pyramidPosition, ent.GetPosition())
    if dist <= 10 {
        ent.data.Hitpoints -= 1
        messaging.SendMessage[messaging.DamageFromEnemyMessage]("", messaging.DamageFromEnemyMessage{Amount: 1})
    }

    if ent.data.Hitpoints <= 0 {
        physics.DestroyCollider(ent.collider)
        soundSide := util.Clamp((ent.GetPosition().X/render.Rs.RenderResolution.X), 0.05, 0.95)
        logging.Debug("%v", soundSide)
        audio.PlaySoundExV2("audio/sounds/explosions/distant1.ogg", 1.0, 1.0, soundSide, 0.1)
        ent.isAlive = false
        messaging.SendMessage[messaging.GainScoreMessage]("", messaging.GainScoreMessage{Amount: ent.data.ScoreWorth, FromKill: true, KilledPosition: ent.GetPosition()})
    }
}

func (ent *EnemyEntity2D) FixedUpdate() {
	// Required update
	ent.BaseEntity2D.FixedUpdate()

    ent.data.MovementPattern(ent.collider)
}

func (ent *EnemyEntity2D) Draw() {
    rl.DrawTexturePro(
        ent.data.Sprite,
        rl.NewRectangle(0, 0, float32(ent.data.Sprite.Width), float32(ent.data.Sprite.Height)),
        rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, ent.data.Size.X, ent.data.Size.Y),
        rl.NewVector2(ent.data.Size.X/2, ent.data.Size.Y/2),
        ent.GetRotation(), rl.Blue)
}

//func (ent *EnemyEntity2D) Draw() {
//    position := ent.GetPosition()
//    rotation := ent.GetRotation() - 90 // Rotation in degrees
//
//    // Define size and length of the triangle
//    size := float32(3) // Half the width of the triangle base
//    length := float32(7) // Length from base to tip
//
//    // Calculate the tip of the triangle
//    tip := rl.Vector2{
//        X: position.X + length * float32(math.Cos(float64(rotation)*math.Pi/180)),
//        Y: position.Y + length * float32(math.Sin(float64(rotation)*math.Pi/180)),
//    }
//
//    // Calculate the two base points
//    baseAngle := float32(math.Pi) / 180 * (rotation + 90) // Perpendicular to the rotation
//    base1 := rl.Vector2{
//        X: position.X + size * float32(math.Cos(float64(baseAngle))),
//        Y: position.Y + size * float32(math.Sin(float64(baseAngle))),
//    }
//    base2 := rl.Vector2{
//        X: position.X + size * float32(math.Cos(float64(baseAngle+math.Pi))),
//        Y: position.Y + size * float32(math.Sin(float64(baseAngle+math.Pi))),
//    }
//
//    col := rl.Red
//    // Draw the triangle
//    rl.DrawLineV(tip, base1,    col)
//    rl.DrawLineV(tip, base2,    col)
//    rl.DrawLineV(base1, base2,  col)
//}
