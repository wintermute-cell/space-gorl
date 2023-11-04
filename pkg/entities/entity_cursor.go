package entities

import (
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/input"
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/util"
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CursorMode int32
const (
    CursorModeNormal CursorMode = iota
    CursorModeShovel
)

// Cursor Entity
type CursorEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite rl.Texture2D
    lastSpriteIdx int32
    sprite_index_shovel int32
    frames_in_sprite int32
    cursorMode CursorMode

    snowSplashSprite rl.Texture2D

    gameState *GameStateHandlerEntity

    movementLocked bool

    shovelAnimationTimer util.Timer
    shovelCooldown util.Timer
    isShovelDown bool

    isClickDown bool

    // sounds
    shovelStabSounds []string
    shovelThrowSounds []string
}

func NewCursorEntity2D(gameState *GameStateHandlerEntity) *CursorEntity2D {
	new_ent := &CursorEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0, Scale: rl.Vector2One()}},

		sprite: rl.LoadTexture("sprites/cursor.png"),
        snowSplashSprite: rl.LoadTexture("sprites/snow_splash.png"),
        sprite_index_shovel: 2,
        frames_in_sprite: 4,
        cursorMode: CursorModeNormal,
        gameState: gameState,
        shovelAnimationTimer: *util.NewTimer(0.44),
        shovelCooldown: *util.NewTimer(0.58),

        shovelStabSounds: []string{
            "shovel_stab1",
            "shovel_stab2",
            "shovel_stab3",
        },
        shovelThrowSounds: []string{
            "shovel_throw1",
            "shovel_throw2",
            "shovel_throw3",
        },
	}
	return new_ent
}

func (ent *CursorEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

    for _, s := range ent.shovelStabSounds {
        audio.RegisterSound(s, fmt.Sprintf("audio/sounds/shovel/%v.ogg", s))
    }
    for _, s := range ent.shovelThrowSounds {
        audio.RegisterSound(s, fmt.Sprintf("audio/sounds/shovel/%v.ogg", s))
    }

    rand.Seed(time.Now().UnixNano())

	// Initialization logic for the entity
	// ...
}

func (ent *CursorEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()
    rl.UnloadTexture(ent.sprite)
    rl.UnloadTexture(ent.snowSplashSprite)

	// De-initialization logic for the entity
	// ...
}


const lerpFactor = 0.8

func lerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

func (ent *CursorEntity2D) calcCursorMode() {
    ent.cursorMode = CursorModeNormal
    if (ent.gameState.MouseHoversSnowpile && ent.gameState.currentPlayState == PlayStateDigging) || ent.isShovelDown {
        ent.cursorMode = CursorModeShovel
    }
}

func (ent *CursorEntity2D) followMousePosition() {
    if !ent.movementLocked {
        currentPos := ent.GetPosition()
        targetPos := rl.GetMousePosition()
        dist := rl.Vector2Distance(currentPos, targetPos)

        if dist > 4 {
            // Calculate dynamic lerp factor based on distance.
            dynamicFactor := lerpFactor / (1.0 + 0.1*dist*dist/128) 

            // Calculate new smoothed position
            newX := lerp(currentPos.X, targetPos.X, dynamicFactor)
            newY := lerp(currentPos.Y, targetPos.Y, dynamicFactor)

            ent.SetPosition(rl.NewVector2(newX, newY))
        } else {
            ent.SetPosition(targetPos)
        }
    }
}

func (ent *CursorEntity2D) Update() {
	ent.BaseEntity2D.Update()

    ent.movementLocked = false
    ent.calcCursorMode()


    switch ent.cursorMode {
    case CursorModeNormal:
        if input.Triggered(input.ActionClickHeld) {
            ent.isClickDown = true
        } else {
            ent.isClickDown = false
        }
    case CursorModeShovel:
        if input.Triggered(input.ActionClickDown) && !ent.isShovelDown && ent.shovelCooldown.Check() {
            ent.shovelAnimationTimer.ResetTime()
            ent.isShovelDown = true
        }
        if ent.shovelAnimationTimer.Check() {
            ent.isShovelDown = false
        }
        if ent.isShovelDown {
            ent.movementLocked = true
        }
    }

    ent.followMousePosition()
}

func (ent *CursorEntity2D) DrawGUI() {
    sprite_idx := int32(0)
    switch ent.cursorMode {
    case CursorModeNormal:
        if ent.isClickDown {
            sprite_idx += 1
        }
        
    case CursorModeShovel:
        if ent.lastSpriteIdx == ent.sprite_index_shovel + 1  && !ent.isShovelDown{
            gem.AddEntity(ent.GetParent(), 
                NewSnowSplashEntity2D(rl.Vector2Add(ent.GetPosition(), rl.NewVector2(3, 30)), ent.snowSplashSprite))
            render.CameraShake(0.2)
            throwIdx := rand.Intn(3)
            audio.PlaySound(ent.shovelThrowSounds[throwIdx])
        }
        sprite_idx = ent.sprite_index_shovel
        if ent.isShovelDown {
            if ent.lastSpriteIdx == ent.sprite_index_shovel {
                stabIdx := rand.Intn(3)
                audio.PlaySound(ent.shovelStabSounds[stabIdx])
            }
            sprite_idx += 1
        }
    }

    ent.lastSpriteIdx = sprite_idx

    frame_width := float32(ent.sprite.Width)/float32(ent.frames_in_sprite)
    p := ent.GetPosition()
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(float32(sprite_idx)*frame_width, 0, frame_width, float32(ent.sprite.Height)),
        rl.NewRectangle(p.X, p.Y, frame_width, float32(ent.sprite.Height)),
        rl.Vector2Zero(),
        0.0,
        rl.White)
}

// GetDrawIndex returns the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *CursorEntity2D) GetDrawIndex() int32 {
	return math.MaxInt32
}
