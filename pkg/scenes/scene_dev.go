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

    origin := rl.NewVector2(0, 16)
    num_frames := int32(32)


    gameStateHandler := entities.NewGameStateHandlerEntity()
    dialogueManager := entities.NewDialogueManagerEntity(gameStateHandler)
    snowpile := entities.NewSnowpileEntity2D(origin, gameStateHandler, dialogueManager)
    cursor := entities.NewCursorEntity2D(gameStateHandler, snowpile)

    bg := entities.NewBackgroundEntity2D(origin, 0, rl.Vector2One())
    gem.AddEntity(scn.scn_root_ent, bg)
    snowL := entities.NewSnowEntity2D(origin, rl.LoadTexture("sprites/snow/snowL.png"), num_frames)
    gem.AddEntity(scn.scn_root_ent, snowL)
    bgFog := entities.NewEffectLayerEntity2D(origin, rl.LoadTexture("sprites/effect_layers/bg_fog.png"), false, 155, true)
    gem.AddEntity(scn.scn_root_ent, bgFog)

    girl := entities.NewGirlEntity2D(origin, num_frames)
    gem.AddEntity(scn.scn_root_ent, girl)

    lightBelow := entities.NewEffectLayerEntity2D(origin, rl.LoadTexture("sprites/effect_layers/light_below.png"), true, 150, false)
    gem.AddEntity(scn.scn_root_ent, lightBelow)

    gem.AddEntity(scn.scn_root_ent, snowpile)

    frostVignette := entities.NewEffectLayerEntity2D(origin, rl.LoadTexture("sprites/effect_layers/frost_vignette.png"), false, 130, false)
    gem.AddEntity(scn.scn_root_ent, frostVignette)
    vignette := entities.NewEffectLayerEntity2D(origin, rl.LoadTexture("sprites/effect_layers/vignette.png"), false, 100, false)
    gem.AddEntity(scn.scn_root_ent, vignette)

    snowM := entities.NewSnowEntity2D(origin, rl.LoadTexture("sprites/snow/snowM.png"), num_frames)
    gem.AddEntity(scn.scn_root_ent, snowM)
    snowS := entities.NewSnowEntity2D(origin, rl.LoadTexture("sprites/snow/snowS.png"), num_frames)
    gem.AddEntity(scn.scn_root_ent, snowS)

    icebar := entities.NewIcebarEntity2D(rl.Vector2Zero(), gameStateHandler)
    gem.AddEntity(scn.scn_root_ent, icebar)

    // TODO: adding an entity with a lower draw index after an entity with a higher one,
    // for example adding dialogueManager after the cursor, causes drawing bugs.
    // One approach to fix would be to replace the current sorting with just some stable sorting algo.

    gem.AddEntity(scn.scn_root_ent, dialogueManager)

    gem.AddEntity(scn.scn_root_ent, cursor)

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
