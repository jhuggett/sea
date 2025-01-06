package models

import "gorm.io/gorm"

type Continent struct {
	gorm.Model

	WorldMapID uint
	WorldMap   WorldMap `gorm:"foreignKey:WorldMapID"`

	Points []*Point `gorm:"foreignKey:ContinentID"`
}
