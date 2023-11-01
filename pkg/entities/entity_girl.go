package entities

import (
	"cowboy-gorl/pkg/animation"
	"cowboy-gorl/pkg/entities/proto"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Girl Entity
type GirlEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite rl.Texture2D
    num_frames int32

    curr_frame int32
    sheet_anim *animation.Animation[int32]
}

func NewGirlEntity2D(position rl.Vector2, num_frames int32) *GirlEntity2D {
	new_ent := &GirlEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},
        sprite: rl.LoadTexture("sprites/girl/girl.png"),
        num_frames: num_frames,
        curr_frame: 0,
	}
    new_ent.sheet_anim = animation.CreateAnimation[int32](float32(num_frames)*0.1)
	for i := 0; i < int(new_ent.num_frames); i++ {
		new_ent.sheet_anim.AddKeyframe(&new_ent.curr_frame, float32(i)*0.1, int32(i))
    }
	return new_ent
}

func (ent *GirlEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

    ent.sheet_anim.Play(true, false)
}

func (ent *GirlEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *GirlEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    ent.sheet_anim.Update()
}

func (ent *GirlEntity2D) Draw() {
    frame_width := float32(ent.sprite.Width/ent.num_frames)
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(float32(ent.curr_frame)*frame_width, 0, frame_width, float32(ent.sprite.Height)),
        rl.NewRectangle(ent.Transform.Position.X, ent.Transform.Position.Y, frame_width, float32(ent.sprite.Height)),
        rl.NewVector2(0, 0),
        ent.Transform.Rotation,
        rl.White,
        )
}
