package weapons

import (
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type LaserBeamProjectile struct {
    CurrentPosition rl.Vector2
    Target TurretTarget
    DistanceTraveled float32
    Progress        float32     // Progress towards the target (0 to 1)

    isActive bool

    Config ProjectileConfig

    hitColliders []physics.RaycastHit
}

func (prj *LaserBeamProjectile) IsContinuous() bool {
    return true
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

func (prj *LaserBeamProjectile) Update() bool {
    if !prj.isActive {
        // begin phasing out logic / animation
        // return true at some point to remove projectile
        return true
    }


    direction := util.Vector2NormalizeSafe(rl.Vector2Subtract(prj.Target.GetPosition(), prj.CurrentPosition)) 

    prj.hitColliders = physics.Raycast(prj.CurrentPosition, direction, prj.Config.maxRange, physics.CollisionCategoryAll)
    if len(prj.hitColliders) > 0 {
        sort.Slice(prj.hitColliders, func(i, j int) bool {
            idist := rl.Vector2Distance(prj.hitColliders[i].HitCollider.GetPosition(), prj.CurrentPosition)
            jdist := rl.Vector2Distance(prj.hitColliders[j].HitCollider.GetPosition(), prj.CurrentPosition)
            return idist < jdist
        })
        prj.hitColliders[0].HitCollider.GetCallbacks()[physics.CollisionCategoryBullet]("apply-damage", float32(180)*rl.GetFrameTime())
    }

    return false
}

func (prj *LaserBeamProjectile) Draw() {
    if prj.isActive {
        rl.DrawLineV(prj.CurrentPosition, prj.Target.GetPosition(), rl.Blue)
    }
    //rl.DrawTexturePro(
    //    prj.Config.sprite,
    //    rl.NewRectangle(0, 0, float32(prj.Config.sprite.Width), float32(prj.Config.sprite.Height)),
    //    rl.NewRectangle(prj.CurrentPosition.X, prj.CurrentPosition.Y, prj.Config.size, prj.Config.size),
    //    rl.NewVector2(prj.Config.size/2, prj.Config.size/2),
    //    util.Vector2Angle(prj.Direction)+90, rl.Blue)
}

func NewLaserBeamProjectile(position rl.Vector2, target TurretTarget, projectileConfig ProjectileConfig) Projectile {
    return &LaserBeamProjectile{
        CurrentPosition:  position,
        Target:        target,
        Config: ProjectileConfig{
            speed:            projectileConfig.speed,
            size:             projectileConfig.size,
            maxRange:         projectileConfig.maxRange,
            sprite:           projectileConfig.sprite,
        },
    }
}
