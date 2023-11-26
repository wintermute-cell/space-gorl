package physics

import (
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/util"
	"sort"

	"github.com/ByteArena/box2d"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RaycastHit struct {
	HitCollider        *Collider
	IntersectionPoint  rl.Vector2
    HitNormal rl.Vector2
}

func createInternalRaycastCallback(results *[]RaycastHit, filter CollisionCategory) box2d.B2RaycastCallback {
    return func(fixture *box2d.B2Fixture, point, normal box2d.B2Vec2, fraction float64) float64 {
        // check if the filter mask contains the category of the collider
        if (fixture.GetFilterData().CategoryBits & uint16(filter)) > 0 {
            // Create a RaycastHit for the hit fixture
            hit := RaycastHit{
                HitCollider: fixture.GetBody().GetUserData().(*Collider),
                IntersectionPoint: simulationToPixelScaleV(rl.Vector2{X: float32(point.X), Y: float32(point.Y)}),
                HitNormal: simulationToPixelScaleV(rl.Vector2{X: float32(normal.X), Y: float32(normal.Y)}),
            }
            *results = append(*results, hit)
            
            return 1  // continue the raycast to get all fixtures in its path
        } else {
            return -1
        }
    }
}


// Raycast casts a ray from origin to direction, returning a list of all
// colliders that were hit.
func Raycast(origin, direction rl.Vector2, length float32, categoriesToHit CollisionCategory) []RaycastHit {
    if length == 0 {
        logging.Warning("Attempted zero length raycast.")
        return []RaycastHit{}
    }
    // translate input values to simulation scale
    oOrigin := origin
    origin = pixelToSimulationScaleV(origin)
    length = pixelToSimulationScale(length)

    // calculate the endpoint
	normalized_direction := util.Vector2NormalizeSafe(direction)
    max_range := length
    endpoint := rl.Vector2Add(origin, rl.Vector2Scale(normalized_direction, max_range))

    // translate to box2d data types
    b2origin := box2d.MakeB2Vec2(float64(origin.X), float64(origin.Y))
    b2endpoint := box2d.MakeB2Vec2(float64(endpoint.X), float64(endpoint.Y))

    // do the raycast
    var results []RaycastHit
    callback := createInternalRaycastCallback(&results, categoriesToHit)
    State.physicsWorld.RayCast(callback, b2origin, b2endpoint)

    // sort results by distance to origin
    sort.Slice(results, func(i, j int) bool {
        return rl.Vector2Distance(oOrigin, results[i].IntersectionPoint) < rl.Vector2Distance(oOrigin, results[j].IntersectionPoint)
    })
    
    return results
}

// Circlecast checks for colliders within a specified radius around a position,
// returning a list of all colliders that were hit within the specified categories.
func Circlecast(position rl.Vector2, radius float32, categoriesToHit CollisionCategory) []RaycastHit {
    if radius <= 0 {
        logging.Warning("Attempted zero or negative radius circlecast.")
        return []RaycastHit{}
    }

    // Translate input values to simulation scale
    simulationPosition := pixelToSimulationScaleV(position)
    simulationRadius := pixelToSimulationScale(radius)

    // Create a circular shape for the circlecast
    circleShape := box2d.NewB2CircleShape()
    circleShape.M_radius = float64(simulationRadius)

    // Create a transform for the circle shape to represent its position
    circleTransform := box2d.B2Transform{}
    circleTransform.P = box2d.MakeB2Vec2(float64(simulationPosition.X), float64(simulationPosition.Y))

    var results []RaycastHit

    // Iterate over all bodies in the physics world
    for b := State.physicsWorld.GetBodyList(); b != nil; b = b.GetNext() {
        for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
            // Check if the fixture's category matches one of the desired categories
            if (f.GetFilterData().CategoryBits & uint16(categoriesToHit)) > 0 {
                fixtureShape := f.GetShape()

                // Check for overlap between the circle shape and the fixture's shape
                if box2d.B2TestOverlapShapes(circleShape, 0, fixtureShape, 0, circleTransform, f.GetBody().GetTransform()) {
                    // Create a RaycastHit for the intersecting fixture
                    hit := RaycastHit{
                        HitCollider: f.GetBody().GetUserData().(*Collider),
                        // IntersectionPoint and HitNormal calculation might be complex in this case
                    }
                    results = append(results, hit)
                }
            }
        }
    }

    // Optionally, sort results by distance to the circle's center
    sort.Slice(results, func(i, j int) bool {
        return rl.Vector2Distance(position, results[i].IntersectionPoint) < rl.Vector2Distance(position, results[j].IntersectionPoint)
    })

    return results
}
