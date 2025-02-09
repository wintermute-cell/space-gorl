package weapons

import (
	"cowboy-gorl/pkg/audio"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BallisticProjectile struct {
	CurrentPosition  rl.Vector2
	Direction        rl.Vector2
	DistanceTraveled float32
	Progress         float32 // Progress towards the target (0 to 1)

	Config ProjectileConfig
}

func (prj *BallisticProjectile) Start() {
	audio.PlaySoundExV2(prj.Config.fireSound, 0.4, 1.0, 0.5, 0.1)
}

func (prj *BallisticProjectile) Stop() {

}

func (prj *BallisticProjectile) IsContinuous() bool {
	return false
}

func (prj *BallisticProjectile) Update() bool {
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
			hit.HitCollider.GetCallbacks()[physics.CollisionCategoryBullet]("apply-damage", float32(2))
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

func (prj *BallisticProjectile) Draw() {
	rl.DrawTexturePro(
		prj.Config.sprite,
		rl.NewRectangle(0, 0, float32(prj.Config.sprite.Width), float32(prj.Config.sprite.Height)),
		rl.NewRectangle(prj.CurrentPosition.X, prj.CurrentPosition.Y, prj.Config.size, prj.Config.size),
		rl.NewVector2(prj.Config.size/2, prj.Config.size/2),
		util.Vector2Angle(prj.Direction)+90, rl.Yellow)
}

func NewBallisticProjectile(position rl.Vector2, target TurretTarget, projectileConfig ProjectileConfig) Projectile {
	direction := rl.Vector2Subtract(target.GetPosition(), position)
	direction = util.RotatePointAroundOrigin(direction, position, (rand.Float32()*2-1)*projectileConfig.accuracy)
	direction = rl.Vector2Normalize(direction)

	return &BallisticProjectile{
		CurrentPosition: position,
		Direction:       direction,
		Config: ProjectileConfig{
			speed:     projectileConfig.speed,
			size:      projectileConfig.size,
			maxRange:  projectileConfig.maxRange,
			sprite:    projectileConfig.sprite,
			fireSound: projectileConfig.fireSound,
		},
		DistanceTraveled: 0,
	}
}
