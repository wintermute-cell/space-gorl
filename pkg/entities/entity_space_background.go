package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/render"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SpaceBackground Entity
type SpaceBackgroundEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    spaceShader rl.Shader
    spaceShaderSchemeLoc int32
    spaceShaderTimeLoc int32
    colorschemeTex rl.Texture2D
}

func NewSpaceBackgroundEntity2D() *SpaceBackgroundEntity2D {
	new_ent := &SpaceBackgroundEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0.0, Scale: rl.Vector2One()}},

        spaceShader: rl.LoadShader("", "shaders/space.glsl"),
        colorschemeTex: rl.LoadTexture("sprites/colorscheme.png"),
	}
    new_ent.spaceShaderSchemeLoc = rl.GetShaderLocation(new_ent.spaceShader, "colorscheme")
    new_ent.spaceShaderTimeLoc = rl.GetShaderLocation(new_ent.spaceShader, "time")
	return new_ent
}

func (ent *SpaceBackgroundEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *SpaceBackgroundEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *SpaceBackgroundEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *SpaceBackgroundEntity2D) Draw() {
    rl.BeginShaderMode(ent.spaceShader)
    rl.SetShaderValueTexture(ent.spaceShader, ent.spaceShaderSchemeLoc, ent.colorschemeTex)
    rl.SetShaderValue(ent.spaceShader, ent.spaceShaderTimeLoc, []float32{float32(rl.GetTime())}, rl.ShaderUniformFloat)
    rl.DrawRectangleV(rl.Vector2Zero(), rl.NewVector2(render.Rs.RenderResolution.X, render.Rs.RenderResolution.X), rl.Black)
    rl.EndShaderMode()
}
