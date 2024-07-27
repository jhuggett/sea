package world_map

import (
	"math"
	"sort"
)

func getAngle(center Point, point Point) float64 {
	x := point.X - center.X
	y := point.Y - center.Y

	angle := math.Atan2(float64(y), float64(x))

	if angle <= 0 {
		return angle + math.Pi*2
	}

	return angle
}

func getDistance(a Point, b Point) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

func comparePoints(center Point, a Point, b Point) bool {
	angleA := getAngle(center, a)
	angleB := getAngle(center, b)

	if angleA < angleB {
		return true
	}

	distanceA := getDistance(center, a)
	distanceB := getDistance(center, b)

	if angleA == angleB && distanceA < distanceB {
		return true
	}

	return false
}

func Sort[T Pointable](points []T) (center Point, sorted []T) {
	center = Point{X: 0, Y: 0}

	for _, point := range points {
		center.X += point.Point().X
		center.Y += point.Point().Y
	}

	center.X /= len(points)
	center.Y /= len(points)

	sort.SliceStable(points, func(i, j int) bool {
		return comparePoints(center, *points[i].Point(), *points[j].Point())
	})

	return center, points
}
