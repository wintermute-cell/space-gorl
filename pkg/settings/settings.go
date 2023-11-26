package settings

import (
	"encoding/json"
	"os"
)

type GameSettings struct {
	// Display
	ScreenWidth     int  `json:"screenWidth"`     // 1280
	ScreenHeight    int  `json:"screenHeight"`    // 720
	RenderWidth     int  `json:"renderWidth"`     // 320
	RenderHeight    int  `json:"renderHeight"`    // 180
	TargetFps       int  `json:"targetFps"`       // 144
	Fullscreen      bool `json:"fullscreen"`      // false
	EnableCrtEffect bool `json:"enableCrtEffect"` // true
	// Gameplay
	MouseSensitivity float32 `json:"mouseSensitivity"` // 1.0
	// Audio
	SoundVolume float32 `json:"soundVolume"` // 0.5
	// Logging
	LogPath string `json:"logPath"` // logs/
	// Controls
	EnableGamepad bool `json:"enableGamepad"` // false
}

var (
	settings *GameSettings
)

// Get the current settings
func CurrentSettings() *GameSettings {
	return settings
}

func FallbackSettings() {
	settings = &GameSettings{
		ScreenWidth:  1280,
		ScreenHeight:  720,
		RenderWidth:  320,
		RenderHeight:  180,
		TargetFps:  144,
		Fullscreen:  false,
		EnableCrtEffect:  true,
		MouseSensitivity:  1.0,
		SoundVolume:  0.5,
		LogPath:  "logs/",
		EnableGamepad:  false,
	}
}

func LoadSettings(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	settings = new(GameSettings)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(settings)
	if err != nil {
		return err
	}

	return nil
}
