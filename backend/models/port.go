package models

import (
	"github.com/jhuggett/sea/models/world_map"
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	CoastalPointID uint
	WorldMapID     uint

	CoastalPoint *world_map.CoastalPoint `gorm:"foreignKey:CoastalPointID"`
}
