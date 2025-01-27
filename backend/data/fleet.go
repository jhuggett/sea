package data

import "gorm.io/gorm"

type Fleet struct {
	gorm.Model

	Ships []Ship

	CapitalShipID uint
	CapitalShip   *Ship `gorm:"foreignKey:CapitalShipID"`
}
