package scenes

import (
	"cowboy-gorl/pkg/ai/navigation"
	"cowboy-gorl/pkg/entities"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/gui"
	"cowboy-gorl/pkg/render"

	//"cowboy-gorl/pkg/lighting"
	"cowboy-gorl/pkg/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*DevScene)(nil)

// Dev Scene
type DevScene struct {
	scn_root_ent   *proto.BaseEntity
	pathable_world navigation.PathableWorld
	start_tile     navigation.Pathable
	end_tile       navigation.Pathable
	g              *gui.Gui
	boolmap        [][]bool
	navmap         *navigation.PathableWorld
}

func (scn *DevScene) Init() {

	scn.g = gui.NewGui()
	scn.scn_root_ent = &proto.BaseEntity{Name: "DevSceneRoot"}
	gem.AddEntity(gem.Root(), scn.scn_root_ent)
	//lighting.Enable()
    render.SetCameraTarget(
        rl.Vector2Scale(render.Rs.RenderResolution, 0.5),
        )

    cursor := entities.NewCursorEntity2D()
    gem.AddEntity(scn.scn_root_ent, cursor)
    
    spaceBackground := entities.NewSpaceBackgroundEntity2D()
    gem.AddEntity(scn.scn_root_ent, spaceBackground)

    pyramid := entities.NewPyramidEntity2D()
    gem.AddEntity(scn.scn_root_ent, pyramid)

    waveManager := entities.NewWaveManagerEntity()
    gem.AddEntity(scn.scn_root_ent, waveManager)

	logging.Info("DevScene initialized.")
}

func (scn *DevScene) Deinit() {
    gem.RemoveEntity(scn.scn_root_ent)
	//lighting.Disable()
	logging.Info("DevScene de-initialized.")
}

func (scn *DevScene) DrawGUI() {
	scn.g.Draw()
}

func (scn *DevScene) Draw() {
}
