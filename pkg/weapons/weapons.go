package weapons

import (
	"cowboy-gorl/pkg/audio"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpriteSheetInfo struct {
	Spritesheet     rl.Texture2D
	FrameDimensions rl.Vector2
	NumFrames       int32
}

type Firearm interface {
	GetShotSoundName() string
	GetSheetInfo() SpriteSheetInfo
	GetMuzzleOffset() rl.Vector2 // the position of the muzzle exit for this weapon, relative to the players origin
}

type weapons_resources struct {
	shotgun_muzzleflash_sheet rl.Texture2D
}

var wr weapons_resources

func Init() {
	wr = weapons_resources{}
	wr.shotgun_muzzleflash_sheet = rl.LoadTexture("sprites/player/weapons/muzzleflash_shotgun.png")
	audio.RegisterSound("shotgun-shot1", "audio/sounds/gunshots/shotgun1.ogg")
}

type firearm struct {
	shot_sound_name         string
	muzzle_flash_sheet_info SpriteSheetInfo
	muzzle_offset           rl.Vector2
}

func NewShotgun() Firearm {
	return &firearm{
		shot_sound_name: "shotgun-shot1",
		muzzle_flash_sheet_info: SpriteSheetInfo{
			Spritesheet:     wr.shotgun_muzzleflash_sheet,
			FrameDimensions: rl.NewVector2(12, 18),
			NumFrames:       1,
		},
		muzzle_offset: rl.NewVector2(-7, 38),
	}
}

func (f *firearm) GetShotSoundName() string {
	return f.shot_sound_name
}

func (f *firearm) GetSheetInfo() SpriteSheetInfo {
	return f.muzzle_flash_sheet_info
}

func (f *firearm) GetMuzzleOffset() rl.Vector2 {
	return f.muzzle_offset
}
