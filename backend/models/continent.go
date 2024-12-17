package models

import "gorm.io/gorm"

type Continent struct {
	gorm.Model

	WorldMapID uint

	Points []*Point `gorm:"foreignKey:ContinentID"`
}
