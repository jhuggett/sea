package data

import "gorm.io/gorm"

type Session struct {
	gorm.Model

	ShipID    uint
	PlayerID  uint
	GameMapID uint
}
