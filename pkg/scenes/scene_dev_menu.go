package scenes

import (
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/gui"
	"cowboy-gorl/pkg/settings"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*DevMenuScene)(nil)

// DevMenu Scene
type DevMenuScene struct {
	// Required fields
	entity_manager *proto.EntityManager

	g *gui.Gui
}

func (scn *DevMenuScene) Init() {
	// Required initialization
	scn.entity_manager = proto.NewEntityManager()
	audio.SetCurrentPlaylist("main-menu", true)

	scn.g = gui.NewGui()
	btn := gui.NewButton("Dev Scene", rl.NewVector2(40, 4), rl.NewVector2(90, 8), func(s gui.ButtonState) {
		if s == gui.ButtonStateReleased {
			Sm.DisableAllScenesExcept([]string{"dev_menu"})
			Sm.EnableScene("dev")
		}
	}, "")

	scn.g.AddWidget(btn)

	// Initialization logic for the scene
	// ...
}

func (scn *DevMenuScene) Deinit() {
	// De-initialization logic for the scene
}

func (scn *DevMenuScene) DrawGUI() {
	scn.g.Draw()
	original_text_size := rg.GetStyle(rg.DEFAULT, rg.TEXT_SIZE)
	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, 16)

	if rg.Button(rl.NewRectangle(4, 4, 16, 16), "") {
		settings.CurrentSettings().EnableCrtEffect = !settings.CurrentSettings().EnableCrtEffect
	}

	rg.SetStyle(rg.DEFAULT, rg.TEXT_SIZE, original_text_size)
}

func (scn *DevMenuScene) Draw() {
	// Draw the scene
}
