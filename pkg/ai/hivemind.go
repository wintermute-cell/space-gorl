package ai

import (
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Hivemind struct {
	controllers []AiController
}

func NewHivemind() *Hivemind {
	return &Hivemind{}
}

func (h *Hivemind) AddController(controller AiController) {
	h.controllers = append(h.controllers, controller)
}

func (h *Hivemind) RemoveController(controller AiController) {
	idx := util.SliceIndex(h.controllers, controller)
	if idx == -1 {
		logging.Error("Tried to remove controller %p, but controller was not found in slice: %v", controller, h.controllers)
	}
	h.controllers = util.SliceDelete(h.controllers, idx, idx+1)
}

func (h *Hivemind) Update() {
	for _, controller := range h.controllers {
		// Calculate group steering behaviors for each controller

		separation := h.calculateSeparation(controller)
		alignment := h.calculateAlignment(controller)
		cohesion := h.calculateCohesion(controller)

		// Average the forces (in a real implementation, you might weight them differently)
		averageForce := rl.Vector2{
			X: (separation.X + alignment.X + cohesion.X) / 3,
			Y: (separation.Y + alignment.Y + cohesion.Y) / 3,
		}
		averageForce = rl.Vector2{
			X: (separation.X + cohesion.X/10000) / 2,
			Y: (separation.Y + cohesion.Y/10000) / 2,
		}

		// Apply this force to the controller's steering behavior
		controller.SetSteeringForce(averageForce)
	}
}

func (h *Hivemind) calculateSeparation(agent AiController) rl.Vector2 {
	const desiredSeparation = 50.0 // This value represents the minimum desired distance between agents
	var steer rl.Vector2
	count := 0

	for _, other := range h.controllers {
		distance := rl.Vector2Distance(agent.GetControllable().GetPosition(), other.GetControllable().GetPosition())
		if distance > 0 && distance < desiredSeparation {
			diff := rl.Vector2Subtract(agent.GetControllable().GetPosition(), other.GetControllable().GetPosition())
			diff = rl.Vector2Normalize(diff)
			diff = rl.Vector2Scale(diff, 1.0/distance) // Give more weight to closer agents
			steer = rl.Vector2Add(steer, diff)
			count++
		}
	}

	if count > 0 {
		steer = rl.Vector2Scale(steer, 1.0/float32(count))
	}

	return steer
}

func (h *Hivemind) calculateAlignment(agent AiController) rl.Vector2 {
	const neighborDist = 100.0 // Maximum distance to consider for alignment
	var sum rl.Vector2
	count := 0

	for _, other := range h.controllers {
		distance := rl.Vector2Distance(agent.GetControllable().GetPosition(), other.GetControllable().GetPosition())
		if distance > 0 && distance < neighborDist {
			sum = rl.Vector2Add(sum, other.GetSteeringForce())
			count++
		}
	}

	if count > 0 {
		sum = rl.Vector2Scale(sum, 1.0/float32(count))
		return sum
	} else {
		return rl.Vector2{}
	}
}

func (h *Hivemind) calculateCohesion(agent AiController) rl.Vector2 {
	const neighborDist = 100.0 // Maximum distance to consider for cohesion
	var sum rl.Vector2
	count := 0

	for _, other := range h.controllers {
		distance := rl.Vector2Distance(agent.GetControllable().GetPosition(), other.GetControllable().GetPosition())
		if distance > 0 && distance < neighborDist {
			sum = rl.Vector2Add(sum, other.GetControllable().GetPosition())
			count++
		}
	}

	if count > 0 {
		sum = rl.Vector2Scale(sum, 1.0/float32(count)) // Get average position
		desired := rl.Vector2Subtract(sum, agent.GetControllable().GetPosition())
		return desired
	} else {
		return rl.Vector2{}
	}
}
