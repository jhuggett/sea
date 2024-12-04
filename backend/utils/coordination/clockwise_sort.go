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

func comparePoints(center Point, a Point, b Point) bool {
	angleA := getAngle(center, a)
	angleB := getAngle(center, b)

	distanceA := getDistance(center, a)
	distanceB := getDistance(center, b)

	if distanceA < distanceB {
		return true
	} else if distanceA > distanceB {
		return false
	}

	if angleA < angleB {
		return true
	} else if angleA > angleB {
		return false
	}

	if distanceA < distanceB {
		return true
	}

	// if distanceA < distanceB {
	// 	if angleA < angleB {
	// 		return true
	// 	} else if angleA > angleB {
	// 		return false
	// 	}
	// } else if distanceA > distanceB {
	// 	if angleA > angleB {
	// 		return true
	// 	} else if angleA < angleB {
	// 		return false
	// 	}
	// }

	return false
}

func Sort2[T Pointable](points []T) (Point, []T) {

	return Point{}, points
}

// Returns the center of the points and the points sorted in clockwise order
func Sort[T Pointable](points []T) (Point, []T) {

	// return Sort2(points)

	center := Point{X: 0, Y: 0}

	for _, point := range points {
		center.X += point.Point().X
		center.Y += point.Point().Y
	}

	center.X /= len(points)
	center.Y /= len(points)

	//	return center, points

	// sort.SliceStable(points, func(i, j int) bool {
	// 	return comparePoints(center, points[i].Point(), points[j].Point())
	// })

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
		// delete(pointMap, current)

		//		foundNeighbor := false

		// possibleNeighbors := []Point{}

		possiblePoints := func(points []Point) []Point {
			p := []Point{}

			for _, neighbor := range points {
				if _, ok := pointMap[neighbor]; ok && !isIncludedInPrevious(neighbor) {
					p = append(p, neighbor)
				}
			}

			return p
		}

		// for _, neighbor := range current.AdjacentDiagonalAndCardinal() {
		// 	if _, ok := pointMap[neighbor]; ok && !neighbor.SameAs(previous) {
		// 		// previous = current
		// 		// current = neighbor
		// 		// foundNeighbor = true

		// 		possibleNeighbors = append(possibleNeighbors, neighbor)
		// 	}
		// }

		// if !foundNeighbor {
		// 	for _, neighbor := range current.AdjacentDiagonal() {
		// 		if _, ok := pointMap[neighbor]; ok && !neighbor.SameAs(previous) {
		// 			// previous = current
		// 			// current = neighbor
		// 			// foundNeighbor = true

		// 			possibleNeighbors = append(possibleNeighbors, neighbor)
		// 		}
		// 	}
		// }
		foundNeighbor := false

		adjacentOptions := possiblePoints(current.AdjacentDiagonalAndCardinal())

		selectNext := func(options []Point) {
			pushPrevious(current)
			current = options[0]
			foundNeighbor = true
		}

		if len(adjacentOptions) > 0 {
			sort.SliceStable(adjacentOptions, func(i, j int) bool {
				return comparePoints(current, adjacentOptions[i], adjacentOptions[j])
			})

			selectNext(adjacentOptions)
		}
		// } else {
		// 	diagonalOptions := possiblePoints(current.AdjacentDiagonal())

		// 	if len(diagonalOptions) > 0 {
		// 		sort.SliceStable(diagonalOptions, func(i, j int) bool {
		// 			return comparePoints(current, diagonalOptions[i], diagonalOptions[j])
		// 		})

		// 		slog.Info("Diagonal options", "current", current, "diagonalOptions", diagonalOptions)

		// 		selectNext(diagonalOptions)
		// 	}
		// }

		// if len(possibleNeighbors) > 0 {
		// 	sort.SliceStable(possibleNeighbors, func(i, j int) bool {
		// 		return comparePoints(current, possibleNeighbors[i], possibleNeighbors[j])
		// 	})

		// 	previous = current
		// 	current = possibleNeighbors[0]
		// 	foundNeighbor = true
		// }

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
