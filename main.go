package main

import (
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/collision"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/gui"
	"cowboy-gorl/pkg/lighting"
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/messaging"
	mousecursor "cowboy-gorl/pkg/mouse_cursor"
	"cowboy-gorl/pkg/physics"

	//"cowboy-gorl/pkg/profiling"
	"cowboy-gorl/pkg/render"
	"cowboy-gorl/pkg/scenes"
	"cowboy-gorl/pkg/settings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// PRE-INIT
	settings_path := "settings.json"
	err := settings.LoadSettings(settings_path)
	if err != nil {
		settings.FallbackSettings()
	}

	logging.Init(settings.CurrentSettings().LogPath)
	logging.Info("Logging initialized")
	if err == nil {
		logging.Info("Settings loaded successfully.")
	} else {
		logging.Warning("Settings loading unsuccessful, using fallback.")
	}

    //profiling.RunPProf("")

	// INITIALIZATION
	// raylib window
	rl.InitWindow(
		int32(settings.CurrentSettings().ScreenWidth),
		int32(settings.CurrentSettings().ScreenHeight), "cowboy-gorl window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(int32(settings.CurrentSettings().TargetFps))

	// rendering
	render.Init(
		settings.CurrentSettings().RenderWidth,
		settings.CurrentSettings().RenderHeight)
	logging.Info("Custom rendering initialized.")

	// initialize audio
	audio.InitAudioV2()
	defer audio.DeinitAudioV2()

	// collision
	collision.InitCollision()
	defer collision.DeinitCollision()

    // physics
    physics.InitPhysics((1.0/60.0), rl.Vector2Zero(), (1.0/32.0))
    defer physics.DeinitPhysics()

    // gem (must come after physics)
    gem.InitGem()

	// lighting
	lighting.InitLighting()
	defer lighting.DeinitLighting()

	// animtion (premades need init and update)
	//animation.InitPremades(render.Rs.CurrentStage.Camera, render.GetWSCameraOffset())

	// register audio tracks
	//audio.RegisterMusic("aza-tumbleweeds", "audio/music/azakaela/azaFMP2_field7_Tumbleweeds.ogg")
	//audio.RegisterMusic("aza-outwest", "audio/music/azakaela/azaFMP2_scene1_OutWest.ogg")
	//audio.RegisterMusic("aza-frontier", "audio/music/azakaela/azaFMP2_town_Frontier.ogg")
	//audio.CreatePlaylist("main-menu", []string{"aza-tumbleweeds", "aza-outwest", "aza-frontier"})
	audio.SetGlobalVolume(0.9)
	audio.SetMusicVolume(0.7)
	audio.SetSFXVolume(0.9)

	// gui
	gui.InitBackend()

    // messaging
    messaging.InitMessageSystem()

	// cursor
	mousecursor.Init()
	rl.HideCursor()

	// raygui
	rg.SetStyle(rg.DEFAULT, rg.TEXT_COLOR_NORMAL, 0x000000)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 24)

    
    // use 1x1 white texture to draw shapes.
    // this fixes UV coordinates not working properly for the rectangle.
    // (see: https://github.com/raysan5/raylib/issues/1730)
    rl.SetShapesTexture(rl.LoadTexture("sprites/1x1white.png"), rl.NewRectangle(0, 0, 1, 1))

	// scenes
	scenes.Sm.RegisterScene("dev", &scenes.DevScene{})
	scenes.Sm.RegisterScene("gameover_menu", &scenes.GameoverMenuScene{})

	scenes.Sm.EnableScene("dev")

	//rl.DisableCursor()

	// GAME LOOP
	for !rl.WindowShouldClose() {
		//animation.UpdatePremades()
        render.UpdateEffects()

		rl.ClearBackground(rl.Black) // clearing the whole background, even behind the main rendertex
		rl.BeginDrawing()

		// begin drawing the world
		render.BeginCustomRenderWorldspace()
		rl.ClearBackground(rl.Black) // clear the main rendertex

		// Draw all registered Scenes
		gem.UpdateEntities()
		gem.DrawEntities()
		scenes.Sm.DrawScenes()
		collision.Update()
        physics.Update()

		// lighting
		lighting.DrawLight()

        //physics.DrawColliders(true, false, false)

		// begin drawing the gui
		render.BeginCustomRenderScreenspace()

		rl.ClearBackground(rl.Blank) // clear the main rendertex

		scenes.Sm.DrawScenesGUI()
        gem.DrawEntitiesGUI()

		render.EndCustomRender()
		//mousecursor.Draw()

		// Draw Debug Info
		rl.DrawFPS(10, 10)

		rl.EndDrawing()

		audio.Update()
	}

	scenes.Sm.DisableAllScenes()
    gem.RemoveEntity(gem.Root())
}
