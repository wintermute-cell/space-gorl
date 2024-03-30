package weapons

import (
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TurretTarget interface {
	GetPosition() rl.Vector2
    IsAlive() bool
}

type Projectile interface {
	Start()
	Update() bool
	Stop()
	Draw()
	IsContinuous() bool
}

type ProjectileFactoryId int32

var projectileFactoryMap map[ProjectileFactoryId]ProjectileFactory = map[ProjectileFactoryId]ProjectileFactory{
	0: NewBallisticProjectile,
    1: NewLaserBeamProjectile,
}

// this is not ok
type ProjectileFactory func(position rl.Vector2, target TurretTarget, projectileConfig ProjectileConfig) Projectile

type ProjectileConfig struct {
	speed    float32
	size     float32
	maxRange float32
	accuracy float32

	// resources
	fireSound string
	sprite    rl.Texture2D
}

func NewProjectileConfig(speed, size, maxRange, accuracy float32, fireSound string, sprite rl.Texture2D) *ProjectileConfig {
	return &ProjectileConfig{
		speed:     speed,
		size:      size,
		maxRange:  maxRange,
		accuracy:  accuracy,
		fireSound: fireSound,
		sprite:    sprite,
	}
}

type TurretConfig struct {
	secondsBetweenShots float32
}

func NewTurretConfig(secondsBetweenShots float32) *TurretConfig {
	return &TurretConfig{
		secondsBetweenShots: secondsBetweenShots,
	}
}

type Turret struct {
	turretConfig     TurretConfig
	projectileConfig ProjectileConfig

	projectilefFactory ProjectileFactory
	position           rl.Vector2
	rofTimer           *util.Timer

	currentTarget     TurretTarget
	activeProjectiles []Projectile
}

func NewTurret(position rl.Vector2, turretConfig TurretConfig, projectileConfig ProjectileConfig, projectileFactoryId ProjectileFactoryId) *Turret {
	nt := &Turret{
		turretConfig:     turretConfig,
		projectileConfig: projectileConfig,

		projectilefFactory: projectileFactoryMap[projectileFactoryId],
		position:           position,
		rofTimer:           util.NewTimer(turretConfig.secondsBetweenShots),
	}
	return nt
}

func (t *Turret) SetTarget(target TurretTarget) {
    if target != t.currentTarget {
        t.stopShooting()
    }
	t.currentTarget = target
}

func (t *Turret) shootAt(target TurretTarget) {
	isContinuousActive := len(t.activeProjectiles) > 0 && t.activeProjectiles[0].IsContinuous()

	// For continuous projectile, don't create a new one if it's already active
	if isContinuousActive {
		return
	}

	// For non-continuous projectile, check rate-of-fire timer
	if !isContinuousActive && !t.rofTimer.Check() {
		return
	}

	prj := t.projectilefFactory(t.position, target, t.projectileConfig)
	prj.Start()
	t.activeProjectiles = append(t.activeProjectiles, prj)
}

func (t *Turret) stopShooting() {
	// Stop all activeProjectiles. This does not remove them, but allows them
	// to initiate their deletion (fading out for example). Removal takes place
	// in FixedUpdate.
	for _, projectile := range t.activeProjectiles {
		projectile.Stop()
	}
}

func (t *Turret) FixedUpdate(newTurretPosition rl.Vector2) {
    if t.currentTarget != nil && !t.currentTarget.IsAlive() {
        t.currentTarget = nil
    }

	if t.currentTarget != nil {
		t.shootAt(t.currentTarget)
	} else {
		t.stopShooting()
	}

	// update projectiles and delete them if they request to be deleted
	for i := 0; i < len(t.activeProjectiles); {
		projectile := t.activeProjectiles[i]
		if projectile.Update() {
			// remove projectile
			t.activeProjectiles = append(t.activeProjectiles[:i], t.activeProjectiles[i+1:]...)
		} else {
			i++
		}
	}
}

func (t *Turret) Draw() {
	for _, projectile := range t.activeProjectiles {
		projectile.Draw()
	}
}

func (t *Turret) SetPosition(position rl.Vector2) {
    t.position = position
}

//type Projectile interface { Start()
//    Update(currentPosition, currentDirection rl.Vector2) bool
//    Stop()
//    Draw()
//    IsActiveIfContinuous() bool
//    //Collide(target Target) bool
//    // other necessary methods...
//}
//
//type ProjectileType int32
//const (
//    ProjectileTypeBallistic ProjectileType = iota
//    ProjectileTypeLaser
//    ProjectileTypeLaserBeam
//)
//
//type ShootingType int32
//const (
//    ShootingTypeSingle ShootingType = iota
//    ShootingTypeContinuous
//)
//
//type ProjectileConfig struct {
//    speed float32
//    size float32
//    maxRange float32
//
//    // resources
//    fireSound string
//    sprite rl.Texture2D
//
//    projectileType ProjectileType
//    isContinuous bool
//}
//
//func NewProjectileConfig(speed, size, maxRange float32, fireSound string, sprite rl.Texture2D, projectileType ProjectileType, isContinuous bool) ProjectileConfig {
//    return ProjectileConfig{
//        speed: speed,
//        size: size,
//        maxRange: maxRange,
//        fireSound: fireSound,
//        sprite: sprite,
//        projectileType: projectileType,
//        isContinuous: isContinuous,
//    }
//}
//
//type Turret struct {
//    Position       rl.Vector2
//    ProjectileConfig     ProjectileConfig
//    ActiveProjectiles []Projectile
//
//    Accuracy float32
//    rofTimer *util.Timer
//}
//
//func NewTurret(position rl.Vector2, projectileConfig ProjectileConfig, accuracy float32, shotCooldown float32) *Turret {
//    return &Turret{
//        Position: position,
//        ProjectileConfig: projectileConfig,
//        Accuracy: accuracy,
//        rofTimer: util.NewTimer(shotCooldown),
//    }
//}
//
//func (t *Turret) Shoot(target rl.Vector2) {
//    direction := rl.Vector2Subtract(target, t.Position)
//    direction = util.RotatePointAroundOrigin(direction, t.Position, (rand.Float32()*2-1)*t.Accuracy)
//    direction = rl.Vector2Normalize(direction)
//
//    var newProjectile Projectile
//    if t.ProjectileConfig.isContinuous {
//        if (len(t.ActiveProjectiles) == 0) || (len(t.ActiveProjectiles) > 0 && !t.ActiveProjectiles[0].IsActiveIfContinuous()) {
//            switch t.ProjectileConfig.projectileType {
//            case ProjectileTypeLaserBeam:
//                newProjectile = NewLaserBeamProjectile(t.Position, direction, t.ProjectileConfig)
//            }
//        } else {
//            return
//        }
//    } else {
//        if !t.rofTimer.Check() {
//            return
//        }
//
//        // Initialize the projectile
//        switch t.ProjectileConfig.projectileType {
//        case ProjectileTypeBallistic:
//            newProjectile = NewBallisticProjectile(t.Position, direction, t.ProjectileConfig)
//        case ProjectileTypeLaser:
//            newProjectile = NewLaserProjectile(t.Position, direction, t.ProjectileConfig)
//        }
//    }
//
//    if newProjectile == nil {
//        logging.Error("Tried to shoot projectile but result was nil. Most likely invalid configuration")
//        return
//    }
//
//    newProjectile.Start()
//    t.ActiveProjectiles = append(t.ActiveProjectiles, newProjectile)
//}
//
//func (t *Turret) StopShooting() {
//    for _, projectile := range t.ActiveProjectiles {
//        projectile.Stop()
//    }
//}
//
//func (t *Turret) FixedUpdate(position, target rl.Vector2) {
//    direction := rl.Vector2Subtract(target, position)
//    direction = util.RotatePointAroundOrigin(direction, position, (rand.Float32()*2-1)*t.Accuracy)
//    direction = rl.Vector2Normalize(direction)
//    for i := 0; i < len(t.ActiveProjectiles); {
//        projectile := t.ActiveProjectiles[i]
//        if projectile.Update(position, direction) {
//            // remove projectile
//            t.ActiveProjectiles = append(t.ActiveProjectiles[:i], t.ActiveProjectiles[i+1:]...)
//        } else {
//            i++
//        }
//    }
//}
//
//func (t *Turret) Draw() {
//    for _, projectile := range t.ActiveProjectiles {
//        projectile.Draw()
//    }
//}
