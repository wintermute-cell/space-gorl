package ai

import (
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/util"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mroth/weightedrand/v2"
)

type AiIKTrollState int32

const (
	IDLE AiIKTrollState = iota
	AiTrollStateRoaming
	AiTrollStateStalking
	AiTrollStateCircling
	AiTrollStateRetreating
)

type AiControllerIKTroll struct {
	pathmap         rl.Image
	controllable    AiControllable
	State           AiIKTrollState
	action_timer    *util.Timer
	pause_timer     *util.Timer
	current_target  *AiTarget
	detection_range float32
	critical_range  float32
	circlingAngle   float64 // Add this field to store the current angle
}

func NewAiControllerIKTroll(pathmap rl.Image, controllable AiControllable, detection_range, critical_range float32) *AiControllerIKTroll {
	new_ai := &AiControllerIKTroll{
		pathmap:         pathmap,
		controllable:    controllable,
		action_timer:    util.NewTimer(10),
		current_target:  &AiTarget{IsCompleted: true},
		detection_range: detection_range,
		critical_range:  critical_range,
	}
	return new_ai
}

// makeTargetOptions returns a list of possible target positions
func (ai *AiControllerIKTroll) makeTargetOptions() rl.Vector2 {
	choices := []weightedrand.Choice[rl.Vector2, int32]{}
	for i := 0; i < 3; i++ {
		new_target_pos := rl.Vector2Zero()
		new_target_pos.X += util.Clamp(rand.Float32()*float32(ai.pathmap.Width), 80, float32(ai.pathmap.Width)-80)
		new_target_pos.Y += util.Clamp(rand.Float32()*float32(ai.pathmap.Height), 80, float32(ai.pathmap.Height)-80)
		red := rl.GetImageColor(ai.pathmap, int32(new_target_pos.X), int32(new_target_pos.Y)).R
		w := int32(0)
		if red > 200 {
			w = 2
		} else {
			w = 1
		}
		choices = append(choices, weightedrand.NewChoice(new_target_pos, w))
	}

	chooser, _ := weightedrand.NewChooser[rl.Vector2, int32](choices...)
	pick := chooser.Pick()
	return pick
}

// Update is called every frame
func (ai *AiControllerIKTroll) Update() {
	playerPos := ai.controllable.GetPlayerPosition()
	trollPos := ai.controllable.GetPosition()
	distanceToPlayer := rl.Vector2Distance(trollPos, playerPos)

	switch ai.State {
	case IDLE:
		if ai.action_timer.Check() {
			ai.State = AiTrollStateRoaming
		}

	case AiTrollStateRoaming:
		if ai.current_target.IsCompleted {
			if ai.controllable.CanSeePlayer() && distanceToPlayer < ai.detection_range {
				ai.State = AiTrollStateStalking
			} else {
				new_target_pos := ai.makeTargetOptions()
				ai.current_target = &AiTarget{Position: new_target_pos}
				ai.controllable.SetAiTarget(ai.current_target)
			}
		}

	case AiTrollStateStalking:
		if distanceToPlayer < ai.critical_range {
			ai.State = AiTrollStateCircling
		} else {
			stalkingPoint := point_towards(trollPos, playerPos, ai.critical_range+25)
			ai.current_target = &AiTarget{Position: stalkingPoint}
			ai.controllable.SetAiTarget(ai.current_target)
		}

	case AiTrollStateCircling:
		// The troll will try to move in a circular path around the player
		ai.CirclingState()
	case AiTrollStateRetreating:
		retreatingPoint := point_away(trollPos, playerPos, ai.detection_range+100)
		ai.current_target = &AiTarget{Position: retreatingPoint}
		ai.controllable.SetAiTarget(ai.current_target)

		if ai.current_target.IsCompleted {
			ai.State = IDLE
		}
	}
}

func (ai *AiControllerIKTroll) CirclingState() {
	playerPos := ai.controllable.GetPlayerPosition()
	trollPos := ai.controllable.GetPosition()

	// If the troll is still far from its target, don't change target.
	logging.Debug("%v", rl.Vector2Distance(trollPos, ai.current_target.Position))
	if rl.Vector2Distance(trollPos, ai.current_target.Position) > 20 {
		return
	}

	// Check for woodland preference.
	redValue := rl.GetImageColor(ai.pathmap, int32(trollPos.X), int32(trollPos.Y)).R
	woodlandPreference := redValue < 100 // Assuming black is close to 0
	circlingRadius := ai.critical_range - 25
	if woodlandPreference {
		circlingRadius += 50 // Increase radius to find more woodland regions.
	}

	// Use the initial red value to determine the direction of circling.
	// This is just a basic mechanism, and can be expanded upon for more complex behaviors.
	// Increase angleStep for more pronounced movement.
	angleStep := math.Pi / 30

	// Add angleStep to the current circling angle.
	ai.circlingAngle += angleStep

	// Ensure the angle is wrapped between 0 and 2π.
	ai.circlingAngle = math.Mod(ai.circlingAngle, 2*math.Pi)

	nextPosition := rl.Vector2{
		X: playerPos.X + circlingRadius*float32(math.Cos(ai.circlingAngle)),
		Y: playerPos.Y + circlingRadius*float32(math.Sin(ai.circlingAngle)),
	}

	// Occasionally pause for a more natural behavior.
	if ai.pause_timer == nil || (ai.pause_timer.Check() && rand.Float32() < 0.1) { // 10% chance to pause.
		ai.current_target = &AiTarget{IsCompleted: true}
		ai.pause_timer = util.NewTimer(rand.Float32() * 5) // Pause for up to 5 seconds.
	} else {
		ai.current_target = &AiTarget{Position: nextPosition}
		ai.controllable.SetAiTarget(ai.current_target)
	}
}

func point_towards(from, to rl.Vector2, distance float32) rl.Vector2 {
	dir := util.Vector2NormalizeSafe(rl.Vector2Subtract(to, from))
	return rl.Vector2Add(from, rl.Vector2Scale(dir, distance))
}

func point_away(from, to rl.Vector2, distance float32) rl.Vector2 {
	dir := util.Vector2NormalizeSafe(rl.Vector2Subtract(from, to))
	return rl.Vector2Add(from, rl.Vector2Scale(dir, -distance))
}

// point_circular returns a point on a circle with the given radius, centered at
// the given center point
func point_circular(from, center rl.Vector2, radius float32) rl.Vector2 {
	dir := util.Vector2NormalizeSafe(rl.Vector2Subtract(from, center))
	return rl.Vector2Add(center, rl.Vector2Scale(dir, radius))
}
