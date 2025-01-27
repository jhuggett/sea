package data

import "gorm.io/gorm"

type EmploymentTerms struct {
	gorm.Model

	Title string

	StartDate uint
	EndDate   uint

	RationsPerDay uint
}
