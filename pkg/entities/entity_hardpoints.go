package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"cowboy-gorl/pkg/weapons"
	"math"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Hardpoints Entity
type HardpointsEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
	turrets []*weapons.Turret

	currentTargets []physics.RaycastHit
    hardpointPositions []rl.Vector2
}

func NewHardpointsEntity2D() *HardpointsEntity2D {
	new_ent := &HardpointsEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: rl.Vector2Zero(), Rotation: 0, Scale: rl.Vector2One()}},

		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *HardpointsEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
	ent.turrets = append(ent.turrets,
		weapons.NewTurret(ent.GetPosition(), *weapons.NewTurretConfig(0.05),
			*weapons.NewProjectileConfig(200, 3, 100, 4.0,
				"audio/sounds/weapons/chaingun1.ogg",
				rl.LoadTexture("sprites/projectiles/ballistic1.png"),
			),
			1,
		),
	)
}

func (ent *HardpointsEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *HardpointsEntity2D) FixedUpdate() {
	for _, trrt := range ent.turrets {
		trrt.FixedUpdate(ent.GetPosition())
	}
}
func (ent *HardpointsEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

    // distribute turrets
	ent.hardpointPositions = util.DistributePointsOnCircleSection(195, 345, len(ent.turrets), 30, ent.GetPosition())
    for i, trrt := range ent.turrets {
        trrt.SetPosition(ent.hardpointPositions[i])
    }

    // Get targets within range
    ent.currentTargets = physics.Circlecast(ent.GetPosition(), 100, physics.CollisionCategoryAll)

// Create sectors based on the number of turrets
sectors := createSectors(195, 345, int32(len(ent.turrets)))

// Initialize a map to keep track of targets not assigned to any sector
unassignedTargets := make(map[*physics.RaycastHit]bool)
for _, target := range ent.currentTargets {
    unassignedTargets[&target] = true
}

// Categorize targets into sectors and mark assigned targets
for _, target := range ent.currentTargets {
    direction := util.Vector2NormalizeSafe(rl.Vector2Subtract(target.HitCollider.GetPosition(), ent.GetPosition()))
    angle := util.Vector2Angle(direction)
    isAssigned := false
    for i := range sectors {
        if sectors[i].isWithinSector(angle) {
            sectors[i].Targets = append(sectors[i].Targets, target)
            isAssigned = true
            break
        }
    }
    if isAssigned {
        delete(unassignedTargets, &target)
    }
}

// Assign targets to turrets
for i, trrt := range ent.turrets {
    var selectedTarget *physics.RaycastHit
    if len(sectors[i].Targets) > 0 {
        // Prefer target within the sector
        selectedTarget = ent.selectMostUrgentTarget(sectors[i].Targets)
    } else if len(unassignedTargets) > 0 {
        // If no target in sector, choose from unassigned targets
        for t := range unassignedTargets {
            selectedTarget = t
            break // Select the first unassigned target (or use another selection logic)
        }
    }

    if selectedTarget != nil {
        enemy, ok := selectedTarget.HitCollider.GetOwner().(*EnemyEntity2D)
        if !ok {
            logging.Error("Hardpoints tried to target something that is not an EnemyEntity2D!")
            return
        }
        trrt.SetTarget(enemy)
    }
}


	//if input.Triggered(input.ActionClickHeld) {
	//    for _, trrt := range ent.turrets {
	//        trrt.Shoot(input.GetCursorPosition())
	//    }
	//}

	//if input.Triggered(input.ActionClickUp) {
	//    for _, trrt := range ent.turrets {
	//        trrt.StopShooting()
	//    }
	//}
}

func (ent *HardpointsEntity2D) Draw() {
	// HARDPOINTS
	// Get the current time
	//rl.DrawCircleV(ent.GetPosition(), 100, rl.NewColor(255, 0, 0, 60))

	currentTime := rl.GetTime()

	for i, point := range ent.hardpointPositions {
		// Calculate a unique phase offset for each point
		phaseOffset := float64(i) * math.Pi / 2.0 // Half Pi offset for each point

		// Calculate the vertical bounce amount (adjust amplitude and frequency as desired)
		bounceAmplitude := float32(3) // Max vertical movement amount
		bounceFrequency := float64(2) // Speed of the bounce
		bounce := bounceAmplitude * float32(math.Sin(currentTime*bounceFrequency+phaseOffset))

		// Apply the bounce to the y-coordinate
		bouncingPoint := rl.Vector2{X: point.X, Y: point.Y + bounce}

		// Draw the point
		rl.DrawCircleV(bouncingPoint, 2, rl.Purple)
	}

	for _, trrt := range ent.turrets {
		trrt.Draw()
	}

}

type Sector struct {
	StartAngle float32
	EndAngle   float32
	Targets    []physics.RaycastHit // Targets within this sector
}

// CreateSectorsWithinCircleSection returns sectors distributed on a section of a 2D circle.
// startAngle and endAngle are in degrees, sectorCount is the number of sectors to create.
func createSectors(startAngle, endAngle float32, sectorCount int32) []Sector {
	if sectorCount <= 0 {
		return nil
	}

	// Convert angles from degrees to radians
	startRad := startAngle * math.Pi / 180.0
	endRad := endAngle * math.Pi / 180.0

	// Calculate angle increment for each sector
	angleIncrement := (endRad - startRad) / float32(sectorCount)

	sectors := make([]Sector, sectorCount)
	for i := int32(0); i < sectorCount; i++ {
		// Calculate the start and end angle for this sector
		sectorStart := startRad + angleIncrement*float32(i)
		sectorEnd := sectorStart + angleIncrement

		sectors[i] = Sector{
			StartAngle: sectorStart * 180.0 / math.Pi, // Convert back to degrees
			EndAngle:   sectorEnd * 180.0 / math.Pi,   // Convert back to degrees
		}
	}
	return sectors
}

func (s Sector) isWithinSector(angle float32) bool {
	return angle >= s.StartAngle && angle < s.EndAngle
}

func (ent *HardpointsEntity2D) selectMostUrgentTarget(targets []physics.RaycastHit) *physics.RaycastHit {
	sort.Slice(targets, func(i, j int) bool {
		idist := rl.Vector2Distance(targets[i].HitCollider.GetPosition(), ent.GetPosition())
		jdist := rl.Vector2Distance(targets[j].HitCollider.GetPosition(), ent.GetPosition())
		return idist < jdist
	})
	return &targets[0]
}
