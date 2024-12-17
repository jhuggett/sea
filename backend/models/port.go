package models

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	PointID    uint
	WorldMapID uint

	Point *Point `gorm:"foreignKey:PointID"`

	// Not sure if it makes sense to have it here, but it's here for now
	InventoryID uint
	Inventory   Inventory `gorm:"foreignKey:InventoryID"`
}
