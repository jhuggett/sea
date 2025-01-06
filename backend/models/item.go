package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model

	Name string

	// Number of items (potentially partial)
	Amount float32

	InventoryID uint
	Inventory   *Inventory `gorm:"foreignKey:InventoryID"`

	PerishDate *uint

	MarkedAsRation bool
}

/*

Need something.

Items could have a % composition of other base materials

*/
