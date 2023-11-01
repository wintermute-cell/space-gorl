package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// EffectLayer Entity
type EffectLayerEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite rl.Texture2D
    isSoftlight bool
    shader rl.Shader
    shaderTex1Loc int32
    tint_color rl.Color
    premultiplied bool

    renderTex rl.RenderTexture2D
}

func NewEffectLayerEntity2D(position rl.Vector2, sprite rl.Texture2D, isSoftlight bool, opacity uint8, premultiplied bool) *EffectLayerEntity2D {
	new_ent := &EffectLayerEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},
        sprite: sprite,
        tint_color: rl.NewColor(255, 255, 255, opacity),
        isSoftlight: isSoftlight,
        premultiplied: premultiplied,
	}

    if new_ent.isSoftlight {
        new_ent.renderTex = rl.LoadRenderTexture(sprite.Width, sprite.Height)
        new_ent.shader = rl.LoadShader("", "shaders/blend-soft-light.glsl")
        new_ent.shaderTex1Loc = rl.GetShaderLocation(new_ent.shader, "texture1")
    }
	return new_ent
}

func (ent *EffectLayerEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *EffectLayerEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *EffectLayerEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *EffectLayerEntity2D) Draw() {
    p := ent.GetPosition()
    y_flipped := float32(render.Rs.CurrentStage.Target.Texture.Height) - (p.Y + float32(ent.sprite.Height))
    if ent.isSoftlight && ent.shader != (rl.Shader{}) {
        // save the current render target
        render.PauseTargetTex()
        rl.BeginTextureMode(ent.renderTex)
        rl.DrawTexturePro(
            render.Rs.CurrentStage.Target.Texture,
            rl.NewRectangle(p.X, y_flipped, float32(ent.sprite.Width), float32(ent.sprite.Height)),
            rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
            rl.NewVector2(0, 0),
            0, rl.White)
        rl.EndTextureMode()
        render.ContinueTargetTex()

        rl.BeginShaderMode(ent.shader)
        rl.SetShaderValueTexture(ent.shader, ent.shaderTex1Loc, ent.renderTex.Texture)
    }
    if ent.premultiplied {
        rl.SetBlendMode(int32(rl.BlendAlphaPremultiply))
    }
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewRectangle(p.X, p.Y, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewVector2(0, 0),
        ent.Transform.Rotation,
        ent.tint_color)
    if ent.premultiplied {
        rl.SetBlendMode(int32(rl.BlendAlpha))
    }

    if ent.isSoftlight && ent.shader != (rl.Shader{}) {
        rl.EndShaderMode()
    }
}
