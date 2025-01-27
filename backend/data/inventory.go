package data

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model

	Items []*Item
}
