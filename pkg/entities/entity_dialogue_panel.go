package entities

import (
	"cowboy-gorl/pkg/animation"
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/input"
	"cowboy-gorl/pkg/util"

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
    currentDialogueNode *DialogueNode
    currentText string
    currentDrawnText string

    hoveredButton *DialogueResponseButtonEntity2D

    showResponses bool

    dialogueHitboxArea rl.Rectangle

    dialogueManager *DialogueManagerEntity

    responseButtons []*DialogueResponseButtonEntity2D
}

func NewDialoguePanelEntity2D(dialogueManager *DialogueManagerEntity) *DialoguePanelEntity2D {
	new_ent := &DialoguePanelEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0, Scale: rl.Vector2One()}},

        bgSprite: rl.LoadTexture("sprites/dialogue_background.png"),
        pagingArrowSprite: rl.LoadTexture("sprites/dialogue/paging_arrow.png"),

        textShowTimer: *util.NewTimer(0.03),

        textFont: rl.LoadFont("fonts/m2x6.fnt"),
        currentText: "",
        currentDrawnText: "",

        dialogueHitboxArea: rl.NewRectangle(0, 112, 170, 74),

        dialogueManager: dialogueManager,
	}
    new_ent.pagingArrowAnim = animation.CreateAnimation[int32](1.0)
    new_ent.pagingArrowAnim.AddKeyframe(&new_ent.pagingArrowCurrFrame, 0, 0)
    new_ent.pagingArrowAnim.AddKeyframe(&new_ent.pagingArrowCurrFrame, 0.5, 1)

	return new_ent
}

func (ent *DialoguePanelEntity2D) setCurrentDialogueNode(node *DialogueNode) {
    ent.textShowTimer.ResetTime()
    ent.currentDialogueNode = node
    ent.currentText = node.Text
    ent.currentDrawnText = ""

    for i := 0; i < len(ent.currentDialogueNode.Responses); i++ {
        ent.responseButtons[i].SetData(node, int32(i))
    }
}

func (ent *DialoguePanelEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()
    ent.pagingArrowAnim.Play(true, false)
    // Minimalist
    // Modern
    // Retro 1, 2
    audio.RegisterSound("new_letter_click", "audio/sounds/ui_soundpack/Minimalist3.ogg")

    for i := 0; i < 3; i++ {
        pos := rl.NewVector2(12, float32(122 + (18*i)))
        newbtn := NewDialogueResponseButtonEntity2D(pos, ent.textFont, ent.dialogueManager, func() {
            for _, btn := range ent.responseButtons {
                gem.RemoveEntity(btn)
            }
            ent.showResponses = false
        })
        ent.responseButtons = append(ent.responseButtons, newbtn)
    }

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
        // show the next letter
        if len(ent.currentDrawnText) % 2 != 1 {
            audio.PlaySound("new_letter_click")
        }
        ent.currentDrawnText = ent.currentText[:len(ent.currentDrawnText)+1]
    }

    if input.TriggeredInArea(input.ActionClickDown, ent.dialogueHitboxArea) {
        if !ent.showResponses {
            if (ent.currentText != ent.currentDrawnText) {
                // skip text drawing animation
                ent.currentDrawnText = ent.currentText
            } else if (ent.currentText == ent.currentDrawnText) {
                // TODO: play small sound
                for i, btn := range ent.responseButtons {

                    // only enable as many buttons as there are responses
                    if i == len(ent.currentDialogueNode.Responses) {
                        break
                    }

                    gem.AddEntity(ent, btn)
                }
                ent.showResponses = true
            }

        }
    }
}

func (ent *DialoguePanelEntity2D) DrawGUI() {
    rl.DrawTextureV(ent.bgSprite, rl.NewVector2(0, 16), rl.White)
    if !ent.showResponses {
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
}

func (ent *DialoguePanelEntity2D) GetDrawIndex() int32 {
	return 1000
}
