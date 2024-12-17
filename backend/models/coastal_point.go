package models

import (
	"github.com/jhuggett/sea/utils/coordination"
	"gorm.io/gorm"
)

type Point struct {
	gorm.Model

	ContinentID uint

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
