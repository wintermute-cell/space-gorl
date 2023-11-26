package weapons

import (
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LaserProjectile struct {
    CurrentPosition rl.Vector2
    Direction rl.Vector2
    DistanceTraveled float32
    Progress        float32     // Progress towards the target (0 to 1)

    Config ProjectileConfig
}

func (prj *LaserProjectile) Start() {

}

func (prj *LaserProjectile) Stop() {

}

func (prj *LaserProjectile) IsActiveIfContinuous() bool {
    return true
}

func (prj *LaserProjectile) Update(_, _ rl.Vector2) bool {
    // Calculate the next ray length (step)
    step := prj.Config.speed * physics.GetTimestep()
    if prj.DistanceTraveled+step > prj.Config.maxRange {
        step = prj.Config.maxRange - prj.DistanceTraveled
    }

    // Calculate perpendicular offset for side raycasts
    perpendicular := rl.Vector2{X: -prj.Direction.Y, Y: prj.Direction.X}
    offsetDistance := prj.Config.size / 2 // Half size on each side
    offset := rl.Vector2Scale(perpendicular, offsetDistance)

    // Perform raycasts
    raycastPositions := []rl.Vector2{
        prj.CurrentPosition,
        rl.Vector2Add(prj.CurrentPosition, offset),
        rl.Vector2Subtract(prj.CurrentPosition, offset),
    }
    hitOccurred := false

    for _, pos := range raycastPositions {
        if hitOccurred {
            break
        }
        hits := physics.Raycast(pos, prj.Direction, step, physics.CollisionCategoryAll)

        for _, hit := range hits {
            // If hit occurs, process the hit and break
            hitOccurred = true
            prj.DistanceTraveled += rl.Vector2Distance(pos, hit.HitCollider.GetPosition())
            hit.HitCollider.GetCallbacks()[physics.CollisionCategoryBullet](float32(2.0))
            prj.CurrentPosition = hit.HitCollider.GetPosition()
            // Process hit (e.g., apply damage)
            // ...
            break
        }
    }

    if !hitOccurred {
        // Move the Projectile forward
        prj.DistanceTraveled += step
        prj.CurrentPosition = rl.Vector2Add(prj.CurrentPosition, rl.Vector2Scale(prj.Direction, step))
    }

    // Remove Projectile if it has reached max range or hit something
    return prj.DistanceTraveled >= prj.Config.maxRange || hitOccurred 
}

func (prj *LaserProjectile) Draw() {
    rl.DrawTexturePro(
        prj.Config.sprite,
        rl.NewRectangle(0, 0, float32(prj.Config.sprite.Width), float32(prj.Config.sprite.Height)),
        rl.NewRectangle(prj.CurrentPosition.X, prj.CurrentPosition.Y, prj.Config.size, prj.Config.size),
        rl.NewVector2(prj.Config.size/2, prj.Config.size/2),
        util.Vector2Angle(prj.Direction)+90, rl.Yellow)
}

func NewLaserProjectile(position, direction rl.Vector2, projectileConfig ProjectileConfig) *LaserProjectile {
    return &LaserProjectile{
        CurrentPosition:  position,
        Direction:        direction,
        Config: ProjectileConfig{
            speed:            projectileConfig.speed,
            size:             projectileConfig.size,
            maxRange:         projectileConfig.maxRange,
            sprite:           projectileConfig.sprite,
        },
        DistanceTraveled: 0,
    }
}

