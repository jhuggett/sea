package coordination

import (
	"github.com/beefsack/go-astar"
)

type ObstacleMap map[int]map[int]bool

func (o *ObstacleMap) AddObstacle(p Point) {
	if _, ok := (*o)[p.X]; !ok {
		(*o)[p.X] = map[int]bool{}
	}
	(*o)[p.X][p.Y] = true
}

func (o *ObstacleMap) RemoveObstacle(p Point) {
	if _, ok := (*o)[p.X]; !ok {
		return
	}
	delete((*o)[p.X], p.Y)
}

func (o *ObstacleMap) IsObstacle(p Point) bool {
	if _, ok := (*o)[p.X]; !ok {
		return false
	}
	if _, ok := (*o)[p.X][p.Y]; !ok {
		return false
	}
	return true
}

type LookUpMap map[int]map[int]*AStartPoint

func (l *LookUpMap) AddOrRetrievePoint(p *AStartPoint) *AStartPoint {
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

	ObstacleMap ObstacleMap
	PointLookUp LookUpMap
}

func (p *AStartPoint) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}

	north := p.North()
	if !p.ObstacleMap.IsObstacle(north) {
		neighbors = append(neighbors, p.PointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       north,
			ObstacleMap: p.ObstacleMap,
			PointLookUp: p.PointLookUp,
		}))
	}

	south := p.South()
	if !p.ObstacleMap.IsObstacle(south) {
		neighbors = append(neighbors, p.PointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       south,
			ObstacleMap: p.ObstacleMap,
			PointLookUp: p.PointLookUp,
		}))
	}

	east := p.East()
	if !p.ObstacleMap.IsObstacle(east) {
		neighbors = append(neighbors, p.PointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       east,
			ObstacleMap: p.ObstacleMap,
			PointLookUp: p.PointLookUp,
		}))
	}

	west := p.West()
	if !p.ObstacleMap.IsObstacle(west) {
		neighbors = append(neighbors, p.PointLookUp.AddOrRetrievePoint(&AStartPoint{
			Point:       west,
			ObstacleMap: p.ObstacleMap,
			PointLookUp: p.PointLookUp,
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
