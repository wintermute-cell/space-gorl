package entities

import (
	"cowboy-gorl/pkg/animation"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SnowSplash Entity
type SnowSplashEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
    lifeTimeTimer util.Timer
    sprite rl.Texture2D
    spriteFrames int32
    spriteCurrentFrame int32
    spriteFrameWidth int32

    spriteAnim *animation.Animation[int32]
}

func NewSnowSplashEntity2D(position rl.Vector2, sprite rl.Texture2D) *SnowSplashEntity2D {
	new_ent := &SnowSplashEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

        lifeTimeTimer: *util.NewTimer(1),
        sprite: sprite,
        spriteFrames: 9,
        spriteCurrentFrame: 0,
	}
    new_ent.spriteFrameWidth = new_ent.sprite.Width/new_ent.spriteFrames
    new_ent.spriteAnim = animation.CreateAnimation[int32](float32(new_ent.spriteFrames)*0.1)
	for i := 0; i < int(new_ent.spriteFrames); i++ {
		new_ent.spriteAnim.AddKeyframe(&new_ent.spriteCurrentFrame, float32(i)*0.1, int32(i))
    }
	return new_ent
}

func (ent *SnowSplashEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()
    ent.spriteAnim.Play(false, false)

	// Initialization logic for the entity
	// ...
}

func (ent *SnowSplashEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *SnowSplashEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()
    ent.spriteAnim.Update()

	// Update logic for the entity
	// ...
    if ent.lifeTimeTimer.Check() {
        gem.RemoveEntity(ent)
    }
}

func (ent *SnowSplashEntity2D) Draw() {
	// Draw logic for the entity
	// ...
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(float32(ent.spriteCurrentFrame*ent.spriteFrameWidth), 0, float32(ent.spriteFrameWidth), float32(ent.sprite.Height)),
        rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, float32(ent.spriteFrameWidth), float32(ent.sprite.Height)),
        rl.NewVector2(float32(ent.spriteFrameWidth)/2, float32(ent.sprite.Height)/2),
        0, rl.White)
}

// GetDrawIndex returns the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *SnowSplashEntity2D) GetDrawIndex() int32 {
	return math.MaxInt32
}
