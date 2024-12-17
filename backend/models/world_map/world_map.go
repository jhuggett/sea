package world_map

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/beefsack/go-astar"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/continent"
	"github.com/jhuggett/sea/utils/coordination"
	"github.com/ojrac/opensimplex-go"
)

type WorldMap struct {
	Persistent models.WorldMap
}

func New() *WorldMap {
	return &WorldMap{}
}

func (w *WorldMap) Create() (uint, error) {
	err := db.Conn().Create(&w.Persistent).Error
	if err != nil {
		return 0, err
	}

	return w.Persistent.ID, nil
}

func (w *WorldMap) Save() error {
	return db.Conn().Save(w.Persistent).Error
}

// Uses Bresenham's Circle Algorithm.
func Circle(x, y, radius int) [][]int {
	points := [][]int{}

	x1, y1, err := -radius, 0, 2-2*radius
	for {
		points = append(
			points,
			[]int{x - x1, y + y1},
			[]int{x - y1, y - x1},
			[]int{x + x1, y - y1},
			[]int{x + y1, y + x1},
		)
		radius = err
		if radius > x1 {
			x1++
			err += x1*2 + 1
		}
		if radius <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}

	return points
}

type Noise struct {
	layers []opensimplex.Noise

	Octaves int
	Scale   float64
}

func NewNoise(octaves int, scale float64) *Noise {
	layers := []opensimplex.Noise{}
	for i := 0; i < octaves; i++ {
		noise := opensimplex.New(rand.Int63())
		layers = append(layers, noise)
	}

	return &Noise{
		layers:  layers,
		Octaves: octaves,
		Scale:   scale,
	}
}

func (n *Noise) Sample(x, y float64) float64 {
	xFloat := x * n.Scale
	yFloat := y * n.Scale
	noise := 0.0
	for _, layer := range n.layers {
		noise += layer.Eval2(xFloat, yFloat)
	}

	return noise / float64(n.Octaves)
}

type Point struct {
	X int
	Y int

	IsCoastal bool

	Elevation float64
}

func (p Point) Point() coordination.Point {
	return coordination.Point{
		X: p.X,
		Y: p.Y,
	}
}

func (w *WorldMap) Generate() error {
	noise := NewNoise(4, 1.0/16.0)

	width, height := 100, 100

	visited := map[coordination.Point]bool{}
	groups := []map[coordination.Point]*Point{}

	isLand := func(point coordination.Point) bool {
		return noise.Sample(float64(point.X), float64(point.Y)) > 0.1
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if _, ok := visited[coordination.Point{X: x, Y: y}]; ok {
				continue
			}

			if isLand(coordination.Point{X: x, Y: y}) {
				group := map[coordination.Point]*Point{}

				flood := []coordination.Point{{X: x, Y: y}}

				for len(flood) > 0 {
					nextFlood := []coordination.Point{}
					for _, point := range flood {
						group[point] = &Point{
							X:         point.X,
							Y:         point.Y,
							IsCoastal: false,
							Elevation: noise.Sample(float64(point.X), float64(point.Y)),
						}
						visited[point] = true

						for _, neighbor := range point.Adjacent() {
							if _, ok := visited[neighbor]; !ok && isLand(neighbor) {
								nextFlood = append(flood, neighbor)
							}
						}
					}
					flood = nextFlood
				}

				groups = append(groups, group)
			}
		}
	}

	for _, group := range groups {

		// need to get the coastal ring of the group

		for point, _ := range group {
			for _, neighbor := range point.Adjacent() {
				if _, ok := group[neighbor]; !ok {
					group[point].IsCoastal = true
				}
			}
		}

		err := w.generateContinent(group)
		if err != nil {
			return fmt.Errorf("error generating continent: %w", err)
		}
	}

	return nil
}

func (w *WorldMap) generateContinent(points map[coordination.Point]*Point) error {
	continent := continent.Continent{
		Persistent: models.Continent{
			WorldMapID: w.Persistent.ID,
		},
	}

	err := db.Conn().Create(&continent.Persistent).Error
	if err != nil {
		return err
	}

	for _, point := range points {
		cp := models.Point{
			ContinentID: continent.Persistent.ID,
			X:           point.X,
			Y:           point.Y,
			Elevation:   point.Elevation,
			Coastal:     point.IsCoastal,
		}

		err := db.Conn().Create(&cp).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(id uint) (*WorldMap, error) {
	var w models.WorldMap
	err := db.Conn().Preload("Continents.Points").First(&w, id).Error
	if err != nil {
		return nil, err
	}

	return &WorldMap{
		Persistent: w,
	}, nil
}

func (w *WorldMap) PlotRoute(starting coordination.Point, ending coordination.Point, obstacles coordination.ObstacleMap) ([]coordination.Point, error) {

	pointMap := coordination.LookUpMap{}

	startingPoint := pointMap.AddOrRetrievePoint(&coordination.AStartPoint{
		Point:       starting,
		ObstacleMap: obstacles,
		PointLookUp: pointMap,
	})

	endingPoint := pointMap.AddOrRetrievePoint(&coordination.AStartPoint{
		Point:       ending,
		ObstacleMap: obstacles,
		PointLookUp: pointMap,
	})

	slog.Debug("PlotRoute", "starting", starting, "ending", ending)

	path, distance, found := astar.Path(startingPoint, endingPoint)

	slog.Debug("Path", "distance", distance, "found", found)

	if !found {
		return nil, nil
	}

	points := []coordination.Point{}

	for _, p := range path {
		points = append(points, coordination.Point{
			X: p.(*coordination.AStartPoint).X,
			Y: p.(*coordination.AStartPoint).Y,
		})
	}

	// reverse points, the astar package returns the path in reverse order
	for i := 0; i < len(points)/2; i++ {
		j := len(points) - i - 1
		points[i], points[j] = points[j], points[i]
	}

	return points, nil
}

func (w *WorldMap) Continents() []*continent.Continent {
	if w.Persistent.Continents == nil {
		// should we return an error here? Should we try looking up the data? Not sure.
		return nil
	}

	continents := []*continent.Continent{}

	for _, c := range w.Persistent.Continents {
		continents = append(continents, &continent.Continent{
			Persistent: *c,
		})
	}

	return continents
}
