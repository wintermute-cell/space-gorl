package entities

import (
	"cowboy-gorl/pkg/animation"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DialoguePanel Entity
type DialoguePanelEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    bgSprite rl.Texture2D
    pagingArrowSprite rl.Texture2D

    pagingArrowAnim *animation.Animation[int32]
    pagingArrowCurrFrame int32

    textShowTimer util.Timer

    textFont rl.Font
    currentText string
    currentDrawnText string
}

func NewDialoguePanelEntity2D() *DialoguePanelEntity2D {
	new_ent := &DialoguePanelEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0, Scale: rl.Vector2One()}},

        bgSprite: rl.LoadTexture("sprites/dialogue_background.png"),
        pagingArrowSprite: rl.LoadTexture("sprites/dialogue/paging_arrow.png"),

        textShowTimer: *util.NewTimer(0.05),

        textFont: rl.LoadFont("fonts/m2x6.fnt"),
        currentText: "Brrr! It's like a snowman convention in here! Do you think they'd mind an extra member?",
        currentDrawnText: "",
	}
    new_ent.pagingArrowAnim = animation.CreateAnimation[int32](1.0)
    new_ent.pagingArrowAnim.AddKeyframe(&new_ent.pagingArrowCurrFrame, 0, 0)
    new_ent.pagingArrowAnim.AddKeyframe(&new_ent.pagingArrowCurrFrame, 0.5, 1)

	return new_ent
}

func (ent *DialoguePanelEntity2D) setCurrentText(text string) {
    ent.textShowTimer.ResetTime()
    ent.currentText = text
}

func (ent *DialoguePanelEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()
    ent.pagingArrowAnim.Play(true, false)

	// Initialization logic for the entity
	// ...
}

func (ent *DialoguePanelEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *DialoguePanelEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()
    ent.pagingArrowAnim.Update()

    if (ent.currentText != ent.currentDrawnText) && ent.textShowTimer.Check() {
        ent.currentDrawnText = ent.currentText[:len(ent.currentDrawnText)+1]
    }
}

func (ent *DialoguePanelEntity2D) DrawGUI() {
    rl.DrawTextureV(ent.bgSprite, rl.NewVector2(0, 16), rl.White)
    util.DrawTextBoxed(
        ent.textFont,
        ent.currentDrawnText,
        rl.NewRectangle(13, 122, 100, 100),
        15, 0, -4, true, rl.White)

    if ent.currentDrawnText == ent.currentText {
        rl.DrawTexturePro(
            ent.pagingArrowSprite,
            rl.NewRectangle(float32(ent.pagingArrowCurrFrame)*8, 0, 8, 8),
            rl.NewRectangle(59, 155+16, 8, 8),
            rl.Vector2Zero(), 0, rl.White)
    }
}

func (ent *DialoguePanelEntity2D) GetDrawIndex() int32 {
	return math.MaxInt32 - 10
}
