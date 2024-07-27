package world_map

import (
	"github.com/jhuggett/sea/db"
	"gorm.io/gorm"
)

type WorldMap struct {
	gorm.Model

	Continents []*Continent `gorm:"foreignKey:WorldMapID"`
}

func New() *WorldMap {
	return &WorldMap{}
}

func (w *WorldMap) Create() (uint, error) {
	err := db.Conn().Create(w).Error
	if err != nil {
		return 0, err
	}

	return w.ID, nil
}

func (w *WorldMap) Save() error {
	return db.Conn().Save(w).Error
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

func (w *WorldMap) generateContinent(x, y, r int) error {
	circlePoints := Circle(x, y, r)

	continent := Continent{
		WorldMapID: w.ID,
	}

	err := db.Conn().Create(&continent).Error
	if err != nil {
		return err
	}

	for _, point := range circlePoints {
		cp := CoastalPoint{
			ContinentID: continent.ID,
			X:           point[0],
			Y:           point[1],
		}

		err := db.Conn().Create(&cp).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WorldMap) GenerateCoasts() error {

	continentsToGenerate := []struct {
		x, y, r int
	}{
		{x: 0, y: 20, r: 3},
		{x: -10, y: -2, r: 10},
		{x: 10, y: -2, r: 6},
	}

	for _, continent := range continentsToGenerate {
		err := w.generateContinent(continent.x, continent.y, continent.r)
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(id uint) (*WorldMap, error) {
	var w WorldMap
	err := db.Conn().Preload("Continents.CoastalPoints").First(&w, id).Error
	if err != nil {
		return nil, err
	}

	return &w, nil
}
