package world_map

import "gorm.io/gorm"

type Continent struct {
	gorm.Model

	WorldMapID uint

	CoastalPoints []*CoastalPoint `gorm:"foreignKey:ContinentID"`
}
