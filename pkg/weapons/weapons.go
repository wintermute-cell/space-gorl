package weapons

import (
	"cowboy-gorl/pkg/util"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Projectile interface {
    Start()
    Update(currentPosition, currentDirection rl.Vector2) bool
    Stop()
    Draw()
    IsActiveIfContinuous() bool
    //Collide(target Target) bool
    // other necessary methods...
}

type ProjectileType int32
const (
    ProjectileTypeBallistic ProjectileType = iota
    ProjectileTypeLaser
    ProjectileTypeLaserBeam
)

type ShootingType int32
const (
    ShootingTypeSingle ShootingType = iota
    ShootingTypeContinuous
)

type ProjectileConfig struct {
    speed float32
    size float32
    maxRange float32

    // resources
    fireSound string
    sprite rl.Texture2D

    projectileType ProjectileType
    isContinuous bool
}

func NewProjectileConfig(speed, size, maxRange float32, fireSound string, sprite rl.Texture2D, projectileType ProjectileType, isContinuous bool) ProjectileConfig {
    return ProjectileConfig{
        speed: speed,
        size: size,
        maxRange: maxRange,
        fireSound: fireSound,
        sprite: sprite,
        projectileType: projectileType,
        isContinuous: isContinuous,
    }
}

type Turret struct {
    Position       rl.Vector2
    ProjectileConfig     ProjectileConfig
    ActiveProjectiles []Projectile

    Accuracy float32
    rofTimer *util.Timer
}

func NewTurret(position rl.Vector2, projectileConfig ProjectileConfig, accuracy float32, shotCooldown float32) *Turret {
    return &Turret{
        Position: position,
        ProjectileConfig: projectileConfig,
        Accuracy: accuracy,
        rofTimer: util.NewTimer(shotCooldown),
    }
}

func (t *Turret) Shoot(target rl.Vector2) {
    direction := rl.Vector2Subtract(target, t.Position)
    direction = util.RotatePointAroundOrigin(direction, t.Position, (rand.Float32()*2-1)*t.Accuracy)
    direction = rl.Vector2Normalize(direction)

    var newProjectile Projectile
    if t.ProjectileConfig.isContinuous {
        if (len(t.ActiveProjectiles) == 0) || (len(t.ActiveProjectiles) > 0 && !t.ActiveProjectiles[0].IsActiveIfContinuous()) {
            switch t.ProjectileConfig.projectileType {
            case ProjectileTypeLaserBeam:
                newProjectile = NewLaserBeamProjectile(t.Position, direction, t.ProjectileConfig)
            }
        }
    } else {
        if !t.rofTimer.Check() {
            return
        }

        // Initialize the projectile
        switch t.ProjectileConfig.projectileType {
        case ProjectileTypeBallistic:
            newProjectile = NewBallisticProjectile(t.Position, direction, t.ProjectileConfig)
        case ProjectileTypeLaser:
            newProjectile = NewLaserProjectile(t.Position, direction, t.ProjectileConfig)
        }
    }

    newProjectile.Start()
    t.ActiveProjectiles = append(t.ActiveProjectiles, newProjectile)
}

func (t *Turret) StopShooting() {
    for _, projectile := range t.ActiveProjectiles {
        projectile.Stop()
    }
}

func (t *Turret) FixedUpdate(position, target rl.Vector2) {
    direction := rl.Vector2Subtract(target, position)
    direction = util.RotatePointAroundOrigin(direction, position, (rand.Float32()*2-1)*t.Accuracy)
    direction = rl.Vector2Normalize(direction)
    for i := 0; i < len(t.ActiveProjectiles); {
        projectile := t.ActiveProjectiles[i]
        if projectile.Update(position, direction) {
            // remove projectile
            t.ActiveProjectiles = append(t.ActiveProjectiles[:i], t.ActiveProjectiles[i+1:]...)
        } else {
            i++
        }
    }
}

func (t *Turret) Draw() {
    for _, projectile := range t.ActiveProjectiles {
        projectile.Draw()
    }
}
