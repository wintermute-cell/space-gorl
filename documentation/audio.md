<!-- LTeX: language=en-US -->
# Audio

The audio module manages all game audio, both playing music and playing
one-shot sounds.

## Basic Usage
Before a Sound or a Music can be played, they need to be *registered*. You do
that by using the corresponding function, for example:
```go
audio.RegisterMusic("my-music-name", "path/to/my/music.ogg")
```

then you may play the music track directly:
```go
audio.PlayMusicNow("my-music-name")
```

or group multiple music tracks into a playlist, and play that:
```go
audio.CreatePlaylist("my-playlist", []string{"my-music-name", "my-other-music"})
audio.SetCurrentPlaylist("my-playlist")
```

The process for playing Sfx/Sounds is similar, though there are no playlists.
For details, please search the [API](#api) section below.

### Example

```go
audio.RegisterMusic("aza-tumbleweeds", "audio/music/azakaela/azaFMP2_field7_Tumbleweeds.ogg")
audio.RegisterMusic("aza-outwest", "audio/music/azakaela/azaFMP2_scene1_OutWest.ogg")
audio.RegisterMusic("aza-frontier", "audio/music/azakaela/azaFMP2_town_Frontier.ogg")
audio.CreatePlaylist("main-menu", []string{"aza-tumbleweeds", "aza-outwest", "aza-frontier"})
audio.SetCurrentPlaylist("main-menu")
// starts playing random songs from the "main-menu" playlist
```

## API

### Loading (Music-/Sfx-)Tracks & Playlists

#### Registering Music:

```go
func RegisterMusic(name, path string) 
```

- **Description**: Load a music file from the specified path and register it with a name.
- **Parameters**:
  - `name`: Name with which the music track will be registered.
  - `path`: Path to the music file.
  
#### Registering Sound:

```go
func RegisterSound(name, path string)
```

- **Description**: Load a sound file from the specified path and register it with a name.
- **Parameters**:
  - `name`: Name with which the sound track will be registered.
  - `path`: Path to the sound file.
  
#### Creating Playlists:

```go
func CreatePlaylist(name string, p []string)
```

- **Description**: Create a music playlist with the specified name and track list.
- **Parameters**:
  - `name`: Name for the new playlist.
  - `p`: List of track names to include in the playlist.

### Configuration

#### Mute & Unmute:

```go
func Mute()
func Unmute()
```

- **Description**: Mute or Unmute the audio.
  
#### Toggle Mute:

```go
func ToggleMute()
```

- **Description**: Toggle the mute state.

#### Set Volumes:

```go
func SetGlobalVolume(new_volume float32)
func SetMusicVolume(new_volume float32)
func SetSFXVolume(new_volume float32)
```

- **Description**: Set the volume levels for global, music, and sound effects.
- **Parameters**:
  - `new_volume`: Desired volume level, range from 0.0 to 1.0.
  
#### Set Music Fade:

```go
func SetMusicFade(fade_secs float32)
```

- **Description**: Set the fade time for music tracks.
- **Parameters**:
  - `fade_secs`: Fade duration in seconds.

### Playback

#### Play Sound:

```go
func PlaySound(name string)
```

- **Description**: Play a registered sound.
- **Parameters**:
  - `name`: Name of the registered sound.

#### Play Sound with Extended Parameters:

```go
func PlaySoundEx(name string, volume, pitch, pan float32)
```

- **Description**: Play a registered sound with specified volume, pitch, and pan settings.
- **Parameters**:
  - `name`: Name of the registered sound.
  - `volume`: Volume level for the sound.
  - `pitch`: Pitch adjustment for the sound. (Default is 1.0)
  - `pan`: Pan position for the sound. (Default is 0.5)

#### Play Music Instantly:

```go
func PlayMusicNow(name string)
```

- **Description**: Instantly start playing a registered music track.
- **Parameters**:
  - `name`: Name of the registered music track.
  
#### Play Music with Fade:

```go
func PlayMusicNowFade(name string)
```

- **Description**: Play a registered music track and fade out the currently playing track.
- **Parameters**:
  - `name`: Name of the registered music track.

#### Set Current Playlist:

```go
func SetCurrentPlaylist(name string, fade_current bool)
```

- **Description**: Set the current playlist and fade out any playing music.
- **Parameters**:
  - `name`: Name of the playlist.
  - `fade_current`: Determines if the currently playing track should be faded out now.

