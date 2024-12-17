package coordination

import (
	"log/slog"
	"math"
	"sort"
)

// In radians
func getAngle(center Point, point Point) float64 {
	dx := float64(point.X - center.X)
	dy := float64(point.Y - center.Y)

	angle := math.Atan2(dy, dx)

	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}

// Euclidean distance
func getDistance(a Point, b Point) float64 {
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

func comparePoints(center Point, a Point, b Point, clockwise bool, nearest bool) bool {
	angleA := getAngle(center, a)
	angleB := getAngle(center, b)

	distanceA := getDistance(center, a)
	distanceB := getDistance(center, b)

	if nearest {
		if distanceA < distanceB {
			return true
		} else if distanceA > distanceB {
			return false
		}
	}

	if clockwise {
		if angleA > angleB {
			return true
		} else if angleA < angleB {
			return false
		}
	} else {
		if angleA < angleB {
			return true
		} else if angleA > angleB {
			return false
		}
	}

	if distanceA < distanceB {
		return true
	}

	return false
}

func Sort2[T Pointable](points []T) []T {

	pointMap := map[Point]T{}

	for _, point := range points {
		pointMap[point.Point()] = point
	}

	startingPoint := points[0].Point()

	previousPoint := startingPoint

	currentPoint := startingPoint

	sortedPoints := []T{}

	// the scan should always be over water

	for {
		sortedPoints = append(sortedPoints, pointMap[currentPoint])

		adjacentPoints := currentPoint.AdjacentDiagonalAndCardinal()

		sort.SliceStable(adjacentPoints, func(i, j int) bool {
			return comparePoints(currentPoint, adjacentPoints[i], adjacentPoints[j], true, true)
		})

		previousPointIndex := 0

		for i, point := range adjacentPoints {
			if point.SameAs(previousPoint) {
				previousPointIndex = i
				break
			}
		}

		if !currentPoint.SameAs(startingPoint) {
			resortedAdjacentPoints := []Point{}

			currentIndex := previousPointIndex + 1

			for len(resortedAdjacentPoints) < len(adjacentPoints) {
				if currentIndex >= len(adjacentPoints) {
					currentIndex = 0
				}

				resortedAdjacentPoints = append(resortedAdjacentPoints, adjacentPoints[currentIndex])

				currentIndex++
			}

			adjacentPoints = resortedAdjacentPoints
		}

		for _, point := range adjacentPoints {
			if _, ok := pointMap[point]; ok {
				previousPoint = currentPoint
				currentPoint = point
				break
			}
		}

		if currentPoint.SameAs(startingPoint) {
			slog.Info("Reached starting point", "currentPoint", currentPoint)
			break
		}
	}

	if len(sortedPoints) != len(points) {
		slog.Error("Did not sort all points", "sorted", len(sortedPoints), "points", len(points))
	}

	return sortedPoints
}

// Returns the center of the points and the points sorted in clockwise order
func Sort[T Pointable](points []T) (Point, []T) {

	center := Point{X: 0, Y: 0}

	for _, point := range points {
		center.X += point.Point().X
		center.Y += point.Point().Y
	}

	center.X /= len(points)
	center.Y /= len(points)

	return center, points

	pointMap := map[Point]T{}

	for _, point := range points {
		pointMap[point.Point()] = point
	}

	start := points[0].Point()

	var previousPoints map[Point]bool = map[Point]bool{}

	isIncludedInPrevious := func(p Point) bool {
		return previousPoints[p]
	}

	pushPrevious := func(p Point) {
		previousPoints[p] = true
	}

	current := start

	reachedStart := false

	sorted := []T{}

	visited := map[Point]bool{}

	for !reachedStart {
		if len(sorted) >= len(points) {
			slog.Warn("Breaking out of loop", "sorted", len(sorted), "points", len(points))
			break
		}

		if _, ok := visited[current]; ok {
			slog.Warn("Already visited", "current", current)
			break
		}

		visited[current] = true

		sorted = append(sorted, pointMap[current])

		possiblePoints := func(points []Point) []Point {
			p := []Point{}

			for _, neighbor := range points {
				if _, ok := pointMap[neighbor]; ok && !isIncludedInPrevious(neighbor) {
					p = append(p, neighbor)
				}
			}

			return p
		}

		foundNeighbor := false

		adjacentOptions := possiblePoints(current.AdjacentDiagonalAndCardinal())

		selectNext := func(options []Point) {
			pushPrevious(current)
			current = options[0]
			foundNeighbor = true
		}

		if len(adjacentOptions) > 0 {
			sort.SliceStable(adjacentOptions, func(i, j int) bool {
				return comparePoints(current, adjacentOptions[i], adjacentOptions[j], true, false)
			})

			selectNext(adjacentOptions)
		}

		if !foundNeighbor {
			if len(sorted) >= len(points) {
				slog.Info("Reached start", "sorted", len(sorted), "points", len(points))
				reachedStart = true
			} else {
				slog.Warn("Did not find neighbor", "current", current, "sorted", len(sorted), "points", len(points))
				return center, sorted
			}
		}
	}

	return center, sorted
}
