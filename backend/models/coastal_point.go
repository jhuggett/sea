package models

import (
	"github.com/jhuggett/sea/utils/coordination"
	"gorm.io/gorm"
)

type CoastalPoint struct {
	gorm.Model

	ContinentID uint

	X int
	Y int
}

func (cp *CoastalPoint) Point() coordination.Point {
	return coordination.Point{
		X: cp.X,
		Y: cp.Y,
	}
}
