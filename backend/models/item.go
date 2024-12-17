package models

import "gorm.io/gorm"

type ItemType string

const (
	ItemTypePieceOfEight ItemType = "piece_of_eight"
	ItemTypeFish         ItemType = "fish"
)

type Item struct {
	gorm.Model

	Name   string
	Amount float32

	InventoryID uint
	Inventory   *Inventory `gorm:"foreignKey:InventoryID"`
}
