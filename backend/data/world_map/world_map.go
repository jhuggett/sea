package world_map

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/beefsack/go-astar"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/continent"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/name"
	"github.com/jhuggett/sea/utils/coordination"
	"github.com/ojrac/opensimplex-go"
	"gorm.io/gorm"
)

type WorldMap struct {
	Persistent data.WorldMap
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
	Octaves []Octave
}

type OctaveConfig struct {
	Scale  float64
	Weight float64
}

type Octave struct {
	Noise opensimplex.Noise

	OctaveConfig
}

func NewNoise(octaveConfigs []OctaveConfig) *Noise {
	octaves := []Octave{}
	for _, config := range octaveConfigs {
		octaves = append(octaves, Octave{
			Noise:        opensimplex.NewNormalized(rand.Int63()),
			OctaveConfig: config,
		})
	}

	return &Noise{
		Octaves: octaves,
	}
}

func (n *Noise) Sample(x, y float64) float64 {
	noise := 0.0
	weight := 0.0

	for _, octave := range n.Octaves {
		noise += octave.Noise.Eval2(x*octave.Scale, y*octave.Scale) * octave.Weight
		weight += octave.Weight
	}

	return noise / weight
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
	noise := NewNoise([]OctaveConfig{
		// {Scale: .05, Weight: .5},
		// {Scale: .1, Weight: 1},
		// {Scale: .2, Weight: .25},
		{Scale: .025, Weight: 1.5},
		{Scale: .095, Weight: .8},
		{Scale: .15, Weight: .4},
	})

	width, height := 100, 100

	visited := map[coordination.Point]bool{}
	groups := []map[coordination.Point]*Point{}

	isLand := func(point coordination.Point) bool {
		return noise.Sample(float64(point.X), float64(point.Y)) > 0.6
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			slog.Debug("Generating", "percent", float32(x+(y*width))/float32(width*height)*100.0, "x", x, "y", y)

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

							// if neighbor.X < 0 || neighbor.X >= width || neighbor.Y < 0 || neighbor.Y >= height {
							// 	continue
							// }

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
	slog.Debug("generateContinent - start")

	start := time.Now()

	continent := continent.Continent{
		Persistent: data.Continent{
			WorldMapID: w.Persistent.ID,

			Name: name.Generate(3),
		},
	}

	err := db.Conn().Create(&continent.Persistent).Error
	if err != nil {
		return fmt.Errorf("error creating continent: %w", err)
	}

	population := data.Population{
		WorldMapID:  w.Persistent.ID,
		ContinentID: continent.Persistent.ID,
		Size:        1000,
	}

	err = db.Conn().Create(&population).Error
	if err != nil {
		return fmt.Errorf("error creating population: %w", err)
	}

	pointsToCreate := []*data.Point{}

	for _, point := range points {
		pointsToCreate = append(pointsToCreate, &data.Point{
			ContinentID: continent.Persistent.ID,
			WorldMapID:  w.Persistent.ID,
			X:           point.X,
			Y:           point.Y,
			Elevation:   point.Elevation,
			Coastal:     point.IsCoastal,
		})

		// err := db.Conn().Create(&cp).Error
		// if err != nil {
		// 	return err
		// }
	}

	err = db.Conn().CreateInBatches(pointsToCreate, 1000).Error
	if err != nil {
		return fmt.Errorf("error creating points: %w", err)
	}

	slog.Debug("generateContinent - end", "duration (s)", time.Since(start).Seconds())

	return nil
}

func Get(id uint) (*WorldMap, error) {
	var w data.WorldMap
	err := db.Conn().Preload("Continents").First(&w, id).Error
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
		// should we return an error here? Should we try looking up the data.? Not sure.
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

func (w *WorldMap) HasLand(x int, y int) (*data.Point, error) {
	var point *data.Point = &data.Point{}
	err := db.Conn().Where("world_map_id = ?", w.Persistent.ID).Where("x = ?", x).Where("y = ?", y).First(point).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return point, nil
}

func (w *WorldMap) CoastalPoints() ([]*data.Point, error) {
	var points []*data.Point
	err := db.Conn().Where("world_map_id = ?", w.Persistent.ID).Where("coastal = ?", true).Find(&points).Error

	if err != nil {
		return nil, err
	}

	return points, nil
}
