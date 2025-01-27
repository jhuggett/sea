package data

import "gorm.io/gorm"

type Producer struct {
	gorm.Model

	InventoryID uint
	Inventory   Inventory `gorm:"foreignKey:InventoryID"`

	PortID uint
	Port   Port `gorm:"foreignKey:PortID"`

	// Comma delimited list of products
	Products string
}
