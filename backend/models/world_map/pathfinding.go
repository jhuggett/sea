package world_map

import (
	"log/slog"

	"github.com/beefsack/go-astar"
)

type ObstacleMap map[int]map[int]bool

func (o *ObstacleMap) AddObstacle(p *Point) {
	if _, ok := (*o)[p.X]; !ok {
		(*o)[p.X] = map[int]bool{}
	}
	(*o)[p.X][p.Y] = true
}

func (o *ObstacleMap) RemoveObstacle(p *Point) {
	if _, ok := (*o)[p.X]; !ok {
		return
	}
	delete((*o)[p.X], p.Y)
}

func (o *ObstacleMap) IsObstacle(p *Point) bool {
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

	obstacleMap ObstacleMap
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

func (w *WorldMap) PlotRoute(starting *Point, ending *Point, obstacles ObstacleMap) ([]Point, error) {

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

	slog.Debug("PlotRoute", "starting", starting, "ending", ending)

	path, distance, found := astar.Path(startingPoint, endingPoint)

	slog.Debug("Path", "distance", distance, "found", found)

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
