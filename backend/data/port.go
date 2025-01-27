package data

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	WorldMapID uint
	WorldMap   WorldMap `gorm:"foreignKey:WorldMapID"`

	PointID uint
	Point   *Point `gorm:"foreignKey:PointID"`

	// Not sure if it makes sense to have it here, but it's here for now
	InventoryID uint
	Inventory   Inventory `gorm:"foreignKey:InventoryID"`

	PopulationID uint
	Population   Population `gorm:"foreignKey:PopulationID"`

	Name string
}
