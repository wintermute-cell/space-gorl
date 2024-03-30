package util

import (
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Rotate a point around the X-axis
func rotateX(point rl.Vector3, angle float32) rl.Vector3 {
    y := point.Y*float32(math.Cos(float64(angle))) - point.Z*float32(math.Sin(float64(angle)))
    z := point.Y*float32(math.Sin(float64(angle))) + point.Z*float32(math.Cos(float64(angle)))
    return rl.Vector3{X: point.X, Y: y, Z: z}
}

// Rotate a point around the Y-axis
func rotateY(point rl.Vector3, angle float32) rl.Vector3 {
    x := point.X*float32(math.Cos(float64(angle))) + point.Z*float32(math.Sin(float64(angle)))
    z := -point.X*float32(math.Sin(float64(angle))) + point.Z*float32(math.Cos(float64(angle)))
    return rl.Vector3{X: x, Y: point.Y, Z: z}
}


// Function to create a 3D pyramid and rotate it based on the given time
func GetRotatedPyramidPoints(currentTime float64, baseSize, height float32) []rl.Vector2 {

    // Define the points of the pyramid in 3D space
    pyramidPoints := []rl.Vector3{
        rl.NewVector3(0, height, 0),                     // Top point
        rl.NewVector3(-baseSize, 0, -baseSize),          // Base vertices
        rl.NewVector3(baseSize, 0, -baseSize),
        rl.NewVector3(baseSize, 0, baseSize),
        rl.NewVector3(-baseSize, 0, baseSize),
    }

    // Rotation angle based on current time (e.g., one full rotation per 10 seconds)
    angleY := float32(currentTime * (2 * math.Pi / 10))

    // Fixed angles for viewing from the side and slightly above
    angleX := float32(math.Pi / -1.5) // 30 degrees
    //angleZ := float32(math.Pi / 6) // 30 degrees (if you want to rotate around Z as well)

    // Rotate points around the Y-axis, then X and Z
    for i, point := range pyramidPoints {
        rotatedPoint := rotateY(point, angleY)
        rotatedPoint = rotateX(rotatedPoint, angleX)
        // rotatedPoint = rotateZ(rotatedPoint, angleZ) // Uncomment if needed
        pyramidPoints[i] = rotatedPoint
    }

    // Project 3D points into 2D
    projectedPoints := make([]rl.Vector2, len(pyramidPoints))
    for i, point := range pyramidPoints {
        // Simple orthographic projection
        projectedPoints[i] = rl.Vector2{X: point.X, Y: point.Z}
    }

    return projectedPoints
}

// DistributePointsOnCircleSection returns points distributed on a section of a 2D circle.
// startAngle and endAngle are in degrees, pointsCount is the number of points to distribute.
func DistributePointsOnCircleSection(startAngle, endAngle float64, pointsCount int, radius float32, center rl.Vector2) []rl.Vector2 {
	if pointsCount <= 0 {
		return nil
	}

	// Convert angles from degrees to radians
	startRad := startAngle * math.Pi / 180.0
	endRad := endAngle * math.Pi / 180.0

    // If only one point is to be distributed, place it in the middle of the section
    if pointsCount == 1 {
        midAngle := (startRad + endRad) / 2
        x := center.X + radius*float32(math.Cos(midAngle))
        y := center.Y + radius*float32(math.Sin(midAngle))
        return []rl.Vector2{{X: x, Y: y}}
    }

	// Calculate angle increment
	angleIncrement := (endRad - startRad) / float64(pointsCount-1)

	points := make([]rl.Vector2, pointsCount)
	for i := 0; i < pointsCount; i++ {
		// Calculate the angle for this point
		angle := startRad + angleIncrement*float64(i)

		// Calculate x and y coordinates
		x := center.X + radius*float32(math.Cos(angle))
		y := center.Y + radius*float32(math.Sin(angle))

		points[i] = rl.Vector2{X: x, Y: y}
	}

	return points
}
