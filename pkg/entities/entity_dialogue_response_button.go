package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/input"
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DialogueResponseButton Entity
type DialogueResponseButtonEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    font rl.Font
    bgSprite rl.Texture2D
    bgSpriteNormal rl.Texture2D
    bgSpritePressed rl.Texture2D
    response DialogueResponse
    currentNode *DialogueNode
    dialogueManager *DialogueManagerEntity

    isHovered bool
    drawnText string
    currentTextScroll int32
    textScrollTimer *util.Timer

     // indicates if the mouse was clicked (button down) on this button
    primed bool

    isClickable bool
    isClickableTimer *util.Timer

    finishResponse func()

    hitbox rl.Rectangle
}

func NewDialogueResponseButtonEntity2D(position rl.Vector2, font rl.Font, dialogueManager *DialogueManagerEntity, afterButtonCallback func()) *DialogueResponseButtonEntity2D {
	new_ent := &DialogueResponseButtonEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

        bgSpriteNormal: rl.LoadTexture("sprites/dialogue_button_normal.png"),
        bgSpritePressed: rl.LoadTexture("sprites/dialogue_button_pressed.png"),

        dialogueManager: dialogueManager,

        isClickableTimer: util.NewTimer(0.3),
        textScrollTimer: util.NewTimer(0.1),

        finishResponse: afterButtonCallback,

        font: font,

        hitbox: rl.NewRectangle(position.X, position.Y, 104, 18),
	}
    new_ent.bgSprite = new_ent.bgSpriteNormal
	return new_ent
}

func (ent *DialogueResponseButtonEntity2D) SetData(node *DialogueNode, responseIdx int32) {
    ent.currentNode = node
    ent.response = node.Responses[responseIdx]
    ent.drawnText = ent.response.Text[0:util.Min(20, len(ent.response.Text))]
    if ent.drawnText != ent.response.Text {
        ent.drawnText += "..."
    }
    ent.currentTextScroll = 0
}

func (ent *DialogueResponseButtonEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()
    ent.isClickableTimer.ResetTime()
    ent.isClickable = false

	// Initialization logic for the entity
	// ...
}

func (ent *DialogueResponseButtonEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *DialogueResponseButtonEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    if rl.CheckCollisionPointRec(rl.GetMousePosition(), ent.hitbox) {
        ent.isHovered = true
        ent.dialogueManager.dialoguePanel.hoveredButton = ent
        ent.drawnText = ent.response.Text
    } else {
        if ent.dialogueManager.dialoguePanel.hoveredButton == ent {
            ent.dialogueManager.dialoguePanel.hoveredButton = nil
        }
        ent.isHovered = false
        ent.drawnText = ent.response.Text[0:util.Min(20, len(ent.response.Text))]
        if ent.drawnText != ent.response.Text {
            ent.drawnText += "..."
        }
    }

    if ent.isClickableTimer.Check() {
        ent.isClickable = true
    }
    if !ent.isClickable {
        return
    }

    if input.TriggeredInArea(input.ActionClickDown, ent.hitbox) {
        ent.bgSprite = ent.bgSpritePressed
        ent.primed = true
    }
    if input.TriggeredInArea(input.ActionClickUp, ent.hitbox) && ent.primed {
        ent.currentNode.WasUsed = true
        ent.dialogueManager.gameState.IceMeterValue -= ent.response.Value
        render.CameraShake(ent.response.Value/30)
        *ent.dialogueManager.gameState.PlayerEmotionalProfile = ent.dialogueManager.gameState.PlayerEmotionalProfile.Add(ent.response.Influence)
        //ent.dialogueManager.gameState = ...
        ent.finishResponse()
        ent.dialogueManager.beginDialogueMode()
    }
    if input.Triggered(input.ActionClickUp) {
        ent.primed = false
        ent.bgSprite = ent.bgSpriteNormal
    }
}

func (ent *DialogueResponseButtonEntity2D) DrawGUI() {
    textOffs := 0
    if ent.bgSprite == ent.bgSpritePressed {
        textOffs += 1
    }
    yscale := float32(1)
    yoffs := float32(0)
    yoffsextra := float32(0)
    if ent.isHovered {
        yscale = 2
        yoffs = float32(ent.bgSpritePressed.Height/2)
        yoffsextra = 1
    }

    tint := rl.White
    hb := ent.dialogueManager.dialoguePanel.hoveredButton 
    if hb != nil && hb != ent {
        tint = rl.Gray
    }
    rl.DrawTexturePro(
        ent.bgSprite,
        rl.NewRectangle(0, 0, float32(ent.bgSprite.Width), float32(ent.bgSprite.Height)),
        rl.NewRectangle(ent.Transform.Position.X, ent.Transform.Position.Y-yoffs, float32(ent.bgSprite.Width), float32(ent.bgSprite.Height)*yscale),
        rl.Vector2Zero(), 0, tint)
    util.DrawTextBoxed(
        ent.font,
        ent.drawnText,
        rl.NewRectangle(ent.GetPosition().X+3, ent.GetPosition().Y + float32(textOffs)-yoffs-yoffsextra, 100, 100),
        15, 0, -5, true, tint)
}

func (ent *DialogueResponseButtonEntity2D) GetDrawIndex() int32 {
    idx := int32(2000)
    if ent.isHovered {
        idx += 5
    }
	return idx
}
