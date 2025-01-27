package data

import "gorm.io/gorm"

type Person struct {
	gorm.Model

	FirstName string
	LastName  string

	NickName string

	Age uint

	Morale float64
}
