package data

import "gorm.io/gorm"

type WorldMap struct {
	gorm.Model

	Continents []*Continent `gorm:"foreignKey:WorldMapID"`
}
