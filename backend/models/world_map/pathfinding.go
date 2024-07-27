package world_map

import (
	"errors"
	"log/slog"

	"github.com/beefsack/go-astar"
	"github.com/soniakeys/raycast"
)

type obstacleMap map[int]map[int]bool

func (o *obstacleMap) AddObstacle(p *Point) {
	if _, ok := (*o)[p.X]; !ok {
		(*o)[p.X] = map[int]bool{}
	}
	(*o)[p.X][p.Y] = true
}

func (o *obstacleMap) IsObstacle(p *Point) bool {
	if _, ok := (*o)[p.X]; !ok {
		return false
	}
	if _, ok := (*o)[p.X][p.Y]; !ok {
		return false
	}
	return true
}

type lookUpMap map[int]map[int]*AStartPoint

func (l *lookUpMap) AddOrRetrievePoint(p *AStartPoint) *AStartPoint {
	if _, ok := (*l)[p.X][p.Y]; ok {
		return (*l)[p.X][p.Y]
	}

	if _, ok := (*l)[p.X]; !ok {
		(*l)[p.X] = map[int]*AStartPoint{}
	}
	(*l)[p.X][p.Y] = p

	return p
}

type AStartPoint struct {
	Point

	obstacleMap obstacleMap
	pointLookUp lookUpMap
}

func (p *AStartPoint) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}

	north := p.North()
	if !p.obstacleMap.IsObstacle(north) {
		neighbors = append(neighbors, p.pointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       *north,
			obstacleMap: p.obstacleMap,
			pointLookUp: p.pointLookUp,
		}))
	}

	south := p.South()
	if !p.obstacleMap.IsObstacle(south) {
		neighbors = append(neighbors, p.pointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       *south,
			obstacleMap: p.obstacleMap,
			pointLookUp: p.pointLookUp,
		}))
	}

	east := p.East()
	if !p.obstacleMap.IsObstacle(east) {
		neighbors = append(neighbors, p.pointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       *east,
			obstacleMap: p.obstacleMap,
			pointLookUp: p.pointLookUp,
		}))
	}

	west := p.West()
	if !p.obstacleMap.IsObstacle(west) {
		neighbors = append(neighbors, p.pointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       *west,
			obstacleMap: p.obstacleMap,
			pointLookUp: p.pointLookUp,
		}))
	}

	return neighbors
}

func (p *AStartPoint) PathNeighborCost(to astar.Pather) float64 {
	return 0.5
}

func (p *AStartPoint) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*AStartPoint)
	absX := toT.X - p.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - p.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

func (w *WorldMap) PlotRoute(starting *Point, ending *Point) ([]Point, error) {
	obstacles := obstacleMap{}

	if starting.SameAs(ending) {
		slog.Info("Starting and ending points are the same")
		return []Point{}, nil
	}

	slog.Info("Plotting Route", "starting", *starting, "ending", *ending)

	endingXY := raycast.XY{
		X: float64(ending.X),
		Y: float64(ending.Y),
	}

	for _, continent := range w.Continents {
		poly := raycast.Poly{}

		Sort(continent.CoastalPoints)

		for _, coastalPoint := range continent.CoastalPoints {

			slog.Info("Coastal Point", "x", coastalPoint.X, "y", coastalPoint.Y)
			if ending.SameAs(coastalPoint.Point()) {
				slog.Info("ending point is a coastal point")
				return nil, errors.New("ending point is a coastal point")
			}
			poly = append(poly, raycast.XY{
				X: float64(coastalPoint.X),
				Y: float64(coastalPoint.Y),
			})
			obstacles.AddObstacle(&Point{
				X: coastalPoint.X,
				Y: coastalPoint.Y,
			})
		}
		if endingXY.In(poly) {
			slog.Info("Ending point is in a continent")
			return nil, errors.New("ending point is in a continent")
		}
	}

	pointMap := lookUpMap{}

	startingPoint := pointMap.AddOrRetrievePoint(&AStartPoint{
		Point:       *starting,
		obstacleMap: obstacles,
		pointLookUp: pointMap,
	})

	endingPoint := pointMap.AddOrRetrievePoint(&AStartPoint{
		Point:       *ending,
		obstacleMap: obstacles,
		pointLookUp: pointMap,
	})

	slog.Info("PlotRoute", "starting", starting, "ending", ending)

	path, distance, found := astar.Path(startingPoint, endingPoint)

	slog.Info("Path", "path", path, "distance", distance, "found", found)

	if !found {
		return nil, nil
	}

	points := []Point{}

	for _, p := range path {
		points = append(points, Point{
			X: p.(*AStartPoint).X,
			Y: p.(*AStartPoint).Y,
		})
	}

	// reverse points, the astar package returns the path in reverse order
	for i := 0; i < len(points)/2; i++ {
		j := len(points) - i - 1
		points[i], points[j] = points[j], points[i]
	}

	return points, nil
}
