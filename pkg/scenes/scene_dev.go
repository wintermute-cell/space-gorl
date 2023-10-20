package scenes

import (
	"cowboy-gorl/pkg/ai/navigation"
	"cowboy-gorl/pkg/entities"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/gui"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/profiling"
	"cowboy-gorl/pkg/render"

	//"cowboy-gorl/pkg/lighting"
	"cowboy-gorl/pkg/logging"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*DevScene)(nil)

// Dev Scene
type DevScene struct {
	entity_manager *proto.EntityManager
	scn_root_ent   *proto.BaseEntity
	pathable_world navigation.PathableWorld
	start_tile     navigation.Pathable
	end_tile       navigation.Pathable
	g              *gui.Gui
	boolmap        [][]bool
	navmap         *navigation.PathableWorld
}

func (scn *DevScene) Init() {
	scn.entity_manager = proto.NewEntityManager()
	scn.g = gui.NewGui()
    map_size := rl.NewVector2(900, 800)
    render.SetCameraClampBounds(rl.NewRectangle(0, 0, map_size.X, map_size.Y))
	scn.scn_root_ent = &proto.BaseEntity{Name: "DevSceneRoot"}
	gem.AddEntity(gem.Root(), scn.scn_root_ent)
	//lighting.Enable()

	p := entities.NewPlayerEntity2D(
		rl.NewVector2(512, 512),
		0.0,
		rl.Vector2One(),
		scn.scn_root_ent,
        func ()  {
			Sm.EnableScene("gameover_menu")
        })

    mag_counter := gui.NewLabel("", rl.NewVector2(20, 80), "color:255,255,255,255")
    mag_counter.WatchInt32(p.GetCurrentMagRef(), "Rounds: %v")
    scn.g.AddWidget(mag_counter)

	env := entities.NewEnvironmentEntity()
	gem.AddEntity(scn.scn_root_ent, p)
	gem.AddEntity(scn.scn_root_ent, env)

	res := int32(8)
	//scn.boolmap = collision.GenerateNavmap([]string{"navmap"}, rl.NewRectangle(0, 0, map_size.X, map_size.Y), res)
    profiling.TimedCall(func() {
        scn.boolmap = physics.GenerateMatrixMap(rl.NewRectangle(0, 0, map_size.X, map_size.Y), res, physics.CollisionCategoryEnvironment)
    }, "GenerateMatrixMap")

    profiling.TimedCall(func() {
        scn.navmap = navigation.NewPathableWorld(rl.NewRectangle(0, 0, map_size.X/float32(res), map_size.Y/float32(res)), res)
    }, "NewPathableWorld")

	for i := 0; i < len(scn.boolmap); i++ {
		column := scn.boolmap[i]
		for j := 0; j < len(column); j++ {
			if scn.boolmap[i][j] {
				scn.navmap.SetCost(rl.NewVector2(float32(i)*float32(res), float32(j)*float32(res)), -1)
			}
		}
	}

	spawner := entities.NewEnemySpawnerEntity(
		[]rl.Vector2{
			{X: 859, Y: 676},
			{X: 113, Y: 759},
			{X: 818, Y: 39},
			{X: 36, Y: 93},
		},
		scn.scn_root_ent,
		p,
		scn.navmap)

	gem.AddEntity(scn.scn_root_ent, spawner)

	logging.Info("DevScene initialized.")
}

func (scn *DevScene) Deinit() {
    gem.RemoveEntity(scn.scn_root_ent)
	scn.entity_manager.DisableAllEntities()
	//lighting.Disable()
    render.SetCameraClampBounds(rl.Rectangle{})
	logging.Info("DevScene de-initialized.")
}

func (scn *DevScene) DrawGUI() {
	scn.g.Draw()
}

func (scn *DevScene) Draw() {
	scn.entity_manager.UpdateEntities()
    //scn.draw_navmap()
}


func (scn *DevScene) draw_navmap() {
	tile_size := float32(8) // this must match the resolution used for GenerateNavmap

	half_tile_size := tile_size / 2.0

	for i := 0; i < len(scn.boolmap); i++ {
		column := scn.boolmap[i]
		for j := 0; j < len(column); j++ {
			if scn.boolmap[i][j] {
				rl.DrawCircle(int32(float32(i)*tile_size + half_tile_size), int32(float32(j)*tile_size + half_tile_size), half_tile_size, rl.Red)
			} else {
				//rl.DrawCircle(int32(float32(i)*tile_size + half_tile_size), int32(float32(j)*tile_size + half_tile_size), half_tile_size, rl.Green)
			}
		}
	}
}
