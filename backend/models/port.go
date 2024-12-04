package models

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	CoastalPointID uint
	WorldMapID     uint

	CoastalPoint *CoastalPoint `gorm:"foreignKey:CoastalPointID"`
}
