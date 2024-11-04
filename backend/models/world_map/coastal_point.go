package world_map

import (
	"gorm.io/gorm"
)

type CoastalPoint struct {
	gorm.Model

	ContinentID uint

	X int
	Y int
}

func (cp *CoastalPoint) Point() *Point {
	return &Point{
		X: cp.X,
		Y: cp.Y,
	}
}
