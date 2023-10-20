package scenes

import (
	"cowboy-gorl/pkg/gui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*GameoverMenuScene)(nil)

// GameoverMenu Scene
type GameoverMenuScene struct {

	g *gui.Gui
}

func (scn *GameoverMenuScene) Init() {
	scn.g = gui.NewGui()
	btn := gui.NewButton("Restart", rl.NewVector2(300, 150), rl.NewVector2(60, 20), func(s gui.ButtonState) {
		if s == gui.ButtonStateReleased {
			Sm.DisableAllScenesExcept([]string{"dev_menu"})
			Sm.EnableScene("dev")
		}
	}, "")

	scn.g.AddWidget(btn)
}

func (scn *GameoverMenuScene) Deinit() {
	// De-initialization logic for the scene
}

func (scn *GameoverMenuScene) DrawGUI() {
	// Draw the GUI for the scene
    scn.g.Draw()
}

func (scn *GameoverMenuScene) Draw() {
	// Draw the scene
}
