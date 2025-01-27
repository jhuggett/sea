package data

import (
	"github.com/jhuggett/sea/utils/coordination"
	"gorm.io/gorm"
)

type Point struct {
	gorm.Model

	WorldMapID uint
	WorldMap   WorldMap `gorm:"foreignKey:WorldMapID"`

	ContinentID uint
	Continent   Continent `gorm:"foreignKey:ContinentID"`

	X int
	Y int

	Coastal   bool
	Elevation float64
}

func (cp *Point) Point() coordination.Point {
	return coordination.Point{
		X: cp.X,
		Y: cp.Y,
	}
}
