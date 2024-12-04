package models

import "gorm.io/gorm"

// Need to move records up to models, then the model level methods use the types from models

type Ship struct {
	gorm.Model
	X float64
	Y float64

	CrewCapacity  uint
	CargoCapacity uint

	WorldMapID uint

	Coffers float64

	IsDocked bool
}
