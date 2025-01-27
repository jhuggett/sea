package data

import "gorm.io/gorm"

type Industry struct {
	gorm.Model

	PopulationID uint
	Population   Population `gorm:"foreignKey:PopulationID"`

	Product string
	Workers uint
}
