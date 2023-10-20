package ai

import (
	"cowboy-gorl/pkg/ai/navigation"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var _ AiController = (*AiControllerEnemyGrunt)(nil)

type AiControllerEnemyGrunt struct {
	controllable     AiControllable
	navmap           *navigation.PathableWorld
	pathing_cooldown util.Timer

	current_path     []rl.Vector2
	current_path_pos int32

	current_target *AiTarget

    steering_force rl.Vector2
}

func NewAiControllerEnemyGrunt(controllable AiControllable, navmap *navigation.PathableWorld) *AiControllerEnemyGrunt {
	new_ai := &AiControllerEnemyGrunt{
		controllable:     controllable,
		navmap:           navmap,
		pathing_cooldown: *util.NewTimer(1),
		current_path:     []rl.Vector2{},
		current_target:   &AiTarget{},
	}
	return new_ai
}

// Update is called every frame
func (ai *AiControllerEnemyGrunt) Update() {

    ai.controllable.SetAiSteeringForce(ai.steering_force)

    if ai.controllable.CanSeePlayer() {
        ai.current_target = &AiTarget{Position: ai.controllable.GetPlayerPosition()}
        ai.controllable.SetAiTarget(ai.current_target)
    } else {
        if ai.pathing_cooldown.Check() && ai.navmap.GetTile(ai.controllable.GetPlayerPosition()).Cost != -1 {
            self_tile := ai.navmap.GetTile(ai.controllable.GetPosition())
            player_tile := ai.navmap.GetTile(ai.controllable.GetPlayerPosition())
            new_path, _, _ := navigation.FindPath(self_tile, player_tile, navigation.NavigationMethodAstar)
            if new_path != nil {
                ai.current_path = []rl.Vector2{}
                ai.current_path_pos = 0
                cut := false // have we pruned the first part?
                for i, p := range new_path {
                    point := rl.Vector2Scale(p.GetPosition(), float32(ai.navmap.Resolution))
                    if len(ai.current_path) > 0 && !cut {
                        // we prune the first part of the path, until the
                        // last point the entity can see. since there are no
                        // obstructions to this point, the entity can go there
                        // directly. This avoids unnatural grid-like movement
                        // during this part.
                        if !ai.controllable.CanSeePoint(point) {
                            ai.current_path = ai.current_path[i-1:]
                            cut = true
                        }
                    }
                    ai.current_path = append([]rl.Vector2{point}, ai.current_path...)
                }
            }
            // removing the first path element fixes the back-rocking
            if len(ai.current_path) > 1 {
                ai.current_path = util.SliceDelete(ai.current_path, 0, 1)
            }
        }

        // if the current path step has been reached, or if we're at the start of the path...
        if ai.current_target.IsCompleted || ai.current_path_pos == 0 {
            // if we are not at the end of the path...
            if ai.current_path_pos < int32(len(ai.current_path)) {
                // set the next path step as a target
                ai.current_target = &AiTarget{Position: ai.current_path[ai.current_path_pos]}
                ai.controllable.SetAiTarget(ai.current_target)
                ai.current_path_pos += 1
            }
        }
    }


    //for _, p := range ai.current_path {
    //    rl.DrawCircleV(p, 4, rl.Blue)
    //}
}

func (ai *AiControllerEnemyGrunt) SetSteeringForce(new_force rl.Vector2) {
    new_force = rl.Vector2Scale(new_force, rl.GetFrameTime() * 3400)
        adjustedPosition := rl.Vector2Add(ai.controllable.GetPosition(), ai.steering_force)
        if ai.navmap.GetTile(adjustedPosition).Cost >= 0 {
            ai.steering_force = new_force
    }
}

func (ai *AiControllerEnemyGrunt) GetSteeringForce() rl.Vector2 {
    return ai.steering_force
}

func (ai *AiControllerEnemyGrunt) GetControllable() AiControllable {
    return ai.controllable
}
