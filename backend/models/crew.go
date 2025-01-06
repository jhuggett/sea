package models

import "gorm.io/gorm"

type Crew struct {
	gorm.Model

	Size int
	// Ship determines the crew capacity

	Wage    float64 // per crew member
	Rations float64 // per crew member

	Morale float64

	ShipID uint
	Ship   *Ship `gorm:"foreignKey:ShipID"`
}
