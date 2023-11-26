package entities

import (
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ==================
// = FORCE APPLIERS =
// ==================

func ApplyTorqueToRotate(collider *physics.Collider, target rl.Vector2, maxTorque float32, alignmentForce, cohesionForce rl.Vector2) {
    currentPosition := collider.GetPosition()

    // Modified direction to target
    directionToTarget := rl.Vector2Add(rl.Vector2Subtract(target, currentPosition), rl.Vector2Add(alignmentForce, cohesionForce))
    desiredAngle := float32(math.Atan2(float64(directionToTarget.Y), float64(directionToTarget.X)))
    currentAngle := collider.GetB2Body().GetAngle()

    // Calculate angular difference
    angleDiff := desiredAngle - float32(currentAngle)
    // Normalize angle to the range [-Pi, Pi]
    for angleDiff < -math.Pi {
        angleDiff += 2 * math.Pi
    }
    for angleDiff > math.Pi {
        angleDiff -= 2 * math.Pi
    }

    // Apply torque proportional to the angle difference, clamped by maxTorque
    torque := util.Clamp(angleDiff, -maxTorque, maxTorque)
    collider.ApplyTorque(torque)
}

func ApplyForceForAcceleration(collider *physics.Collider, acceleration float32, separationForce rl.Vector2) {
    currentAngle := collider.GetB2Body().GetAngle()
    direction := rl.Vector2{X: float32(math.Cos(float64(currentAngle))), Y: float32(math.Sin(float64(currentAngle)))}

    // Modified force direction
    modifiedDirection := rl.Vector2Add(direction, separationForce)
    force := rl.Vector2Scale(modifiedDirection, acceleration)

    collider.ApplyForceToCenter(force)
}

// =======================
// = STEERING BEHAVIOURS =
// =======================

func CalculateSeparationForce(collider *physics.Collider, nearbyEnemies []*EnemyEntity2D, separationStrength float32) rl.Vector2 {
    var force rl.Vector2
    for _, neighbor := range nearbyEnemies {
        diff := rl.Vector2Subtract(collider.GetPosition(), neighbor.collider.GetPosition())
        distance := rl.Vector2Length(diff)
        if distance < 10 && distance > 0 { // separationThreshold is a defined constant
            pushForce := rl.Vector2Scale(util.Vector2NormalizeSafe(diff), separationStrength/distance)
            force = rl.Vector2Add(force, pushForce)
        }
    }
    return force
}

func CalculateAlignmentForce(nearbyEnemies []*EnemyEntity2D, alignmentStrength float32) rl.Vector2 {
    var averageVelocity rl.Vector2
    var count int
    for _, neighbor := range nearbyEnemies {
        averageVelocity = rl.Vector2Add(averageVelocity, neighbor.collider.GetVelocity()) // Assuming Velocity field exists
        count++
    }
    if count > 0 {
        averageVelocity = rl.Vector2Scale(averageVelocity, 1/float32(count))
        return rl.Vector2Scale(util.Vector2NormalizeSafe(averageVelocity), alignmentStrength)
    }
    return rl.Vector2{}
}

func CalculateCohesionForce(collider *physics.Collider, nearbyEnemies []*EnemyEntity2D, cohesionStrength float32) rl.Vector2 {
    var centerOfMass rl.Vector2
    var count int
    for _, neighbor := range nearbyEnemies {
        centerOfMass = rl.Vector2Add(centerOfMass, neighbor.collider.GetPosition())
        count++
    }
    if count > 0 {
        centerOfMass = rl.Vector2Scale(centerOfMass, 1/float32(count))
        toCenter := rl.Vector2Subtract(centerOfMass, collider.GetPosition())
        return rl.Vector2Scale(util.Vector2NormalizeSafe(toCenter), cohesionStrength)
    }
    return rl.Vector2{}
}

