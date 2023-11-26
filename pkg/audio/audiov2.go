package audio

import (
	"cowboy-gorl/pkg/logging"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// soundV2 represents a sound with its Wave data and a pool of sound instances.
type soundV2 struct {
    Wave rl.Wave
    Pool []*rl.Sound
}

type audioSystem struct {
	// volume
	isMute       bool
	globalVolume float32
	musicVolume  float32
	soundVolume    float32

	// tracks
	musicRegistry    map[string]rl.Music
	soundRegistry   map[string]*soundV2

    // playlists
	musicPlaylists map[string]playlist
}

var as audioSystem

// InitAudioV2 initializes the audio system, sets up the device, and prepares
// the sound registry.
func InitAudioV2() {
    rl.InitAudioDevice()
    if !rl.IsAudioDeviceReady() {
        logging.Fatal("Failed to InitAudioDevice!")
    }
    as = audioSystem{
        isMute: false,
        globalVolume: 0.5,
        musicVolume: 0.8,
        soundVolume: 0.8,

        musicRegistry: make(map[string]rl.Music),
        soundRegistry: make(map[string]*soundV2),

        musicPlaylists: make(map[string]playlist),
    }
}

// DeinitAudioV2 deinitializes the audio system, unloads all sounds and closes
// the audio device.
func DeinitAudioV2() {
    // Deinitialize all sound pools and unload waves
    for _, sound := range as.soundRegistry {
        for _, s := range sound.Pool {
            rl.UnloadSound(*s)
        }
        rl.UnloadWave(sound.Wave)
    }
    rl.CloseAudioDevice()
}


func RegisterSoundV2(path string) {
    if _, ok := as.soundRegistry[path]; ok {
        logging.Warning("Sound already registered: %v", path)
        return
    }

    wave := rl.LoadWave(path)
    if wave.FrameCount == 0 {
        logging.Error("Failed to load wave for path: %v", path)
        return
    }

    as.soundRegistry[path] = &soundV2{
        Wave: wave,
        Pool: make([]*rl.Sound, 0),
    }
}

/*
PlaySoundV2 plays a registered sound with default parameters.

- path: The file path of the sound to play. The sound must have been previously registered using RegisterSoundV2. This function is a wrapper around PlaySoundExV2, using default values for volume (1.0), pitch (1.0), and pan (0.5).
*/
func PlaySoundV2(path string) {
    PlaySoundEx(path, 1.0, 1.0, 0.5) // Default values for volume, pitch, and pan
}

/*
PlaySoundExV2 plays a sound with extended parameters.

- path: The file path of the sound to play. The sound must have been previously registered using RegisterSoundV2.

- volume: The volume at which to play the sound, ranging from 0.0 (silent) to 1.0 (full volume).

- pitch: The pitch at which to play the sound. A value of 1.0 plays the sound at normal pitch, values greater than 1.0 increase the pitch, and values less than 1.0 decrease it.

- pan: The pan at which to play the sound, ranging from 0.0 (left channel only) to 1.0 (right channel only), with 0.5 being centered.

- pitchVariance: The range of random pitch variation to apply. For example, a variance of 0.1 could randomly alter the pitch by up to ±0.1 around the base pitch value.
*/
func PlaySoundExV2(path string, volume, pitch, pan, pitchVariance float32) {
    sound, ok := as.soundRegistry[path]
    if !ok {
        logging.Warning("Attempted to play sound that is not registered: %v", path)
        logging.Warning("Registering missing sound: %v", path)
        RegisterSoundV2(path)
        sound, ok = as.soundRegistry[path]
        if !ok {
            logging.Warning("Failed to register missing sound: %v", path)
            return
        }
    }

    var soundInstance *rl.Sound
    for _, s := range sound.Pool {
        if !rl.IsSoundPlaying(*s) {
            soundInstance = s
            break
        }
    }

    if soundInstance == nil {
        logging.Info("Creating new sound instance for sound: %v", path)
        newSound := rl.LoadSoundFromWave(sound.Wave)
        sound.Pool = append(sound.Pool, &newSound)
        soundInstance = &newSound
    }

    rand.Seed(time.Now().UnixNano())
    randomVariance := pitchVariance * (2*rand.Float32() - 1) // Generate random number in range [-pitchVariance, pitchVariance]

    // Set sound properties
    rl.SetSoundPitch(*soundInstance, pitch+randomVariance)
    rl.SetSoundPan(*soundInstance, 1 - pan) // raylib does pan from 1 (left) to 0 (right) for some reason lol
    rl.SetSoundVolume(*soundInstance, volume) // Assuming global volume is handled elsewhere

    rl.PlaySound(*soundInstance)
}

