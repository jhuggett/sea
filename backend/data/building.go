package data

import "gorm.io/gorm"

type Building struct {
	gorm.Model

	PortID uint

	Name string
	Type string

	X int
	Y int
}
