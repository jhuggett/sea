package data

import "gorm.io/gorm"

type Industry struct {
	gorm.Model

	// PopulationID uint
	// Population   Population `gorm:"foreignKey:PopulationID"`

	EconomyID uint
	Economy   Economy `gorm:"foreignKey:EconomyID"`

	Product        string
	ShareOfWorkers float64 // How mush of the workforce is dedicated to this industry
}
