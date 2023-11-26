package weapons

import (
	"cowboy-gorl/pkg/physics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LaserBeamProjectile struct {
    CurrentPosition rl.Vector2
    Direction rl.Vector2
    DistanceTraveled float32
    Progress        float32     // Progress towards the target (0 to 1)

    isActive bool

    Config ProjectileConfig
}

func (prj *LaserBeamProjectile) IsActiveIfContinuous() bool {
    return prj.isActive
}

func (prj *LaserBeamProjectile) Start() {
    // Start the laser
    prj.isActive = true
    // other start logic...
}

func (prj *LaserBeamProjectile) Stop() {
    // Stop the laser
    prj.isActive = false
    // other stop logic...
}

func (prj *LaserBeamProjectile) Update(position, direction rl.Vector2) bool {
    if !prj.isActive {
        // begin phasing out logic / animation
        // return true at some point to remove projectile
        return true
    }

    prj.CurrentPosition = position
    prj.Direction = direction

    hits := physics.Raycast(prj.CurrentPosition, prj.Direction, prj.Config.maxRange, physics.CollisionCategoryAll)

    for _, hit := range hits {
        // If hit occurs, process the hit and break
        // Process hit (e.g., apply damage)
        // ...
        _ = hit
        break
    }

    return false
}

func (prj *LaserBeamProjectile) Draw() {
    endPoint := rl.Vector2Add(prj.CurrentPosition, rl.Vector2Scale(prj.Direction, prj.Config.maxRange))
    rl.DrawLineV(prj.CurrentPosition, endPoint, rl.Blue)
    //rl.DrawTexturePro(
    //    prj.Config.sprite,
    //    rl.NewRectangle(0, 0, float32(prj.Config.sprite.Width), float32(prj.Config.sprite.Height)),
    //    rl.NewRectangle(prj.CurrentPosition.X, prj.CurrentPosition.Y, prj.Config.size, prj.Config.size),
    //    rl.NewVector2(prj.Config.size/2, prj.Config.size/2),
    //    util.Vector2Angle(prj.Direction)+90, rl.Blue)
}

func NewLaserBeamProjectile(position, direction rl.Vector2, projectileConfig ProjectileConfig) *LaserBeamProjectile {
    return &LaserBeamProjectile{
        CurrentPosition:  position,
        Direction:        direction,
        Config: ProjectileConfig{
            speed:            projectileConfig.speed,
            size:             projectileConfig.size,
            maxRange:         projectileConfig.maxRange,
            sprite:           projectileConfig.sprite,
            isContinuous: true,
        },
        DistanceTraveled: 0,
        
    }
}



