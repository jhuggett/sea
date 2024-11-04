package models

import "gorm.io/gorm"

type Crew struct {
	gorm.Model

	ShipID uint
	Ship   *Ship `gorm:"foreignKey:ShipID"`
}
