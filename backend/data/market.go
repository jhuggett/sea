package data

import "gorm.io/gorm"

type Market struct {
	gorm.Model

	EconomyID uint
	Economy   Economy `gorm:"foreignKey:EconomyID"`

	ProductName string

	HistoricalSupplyAverage float64
	HistoricalDemandAverage float64

	HistoricalSurplus float64
	Surplus           float64

	SampleCount uint
}

/*

- what average amount has been wanted historically

- how much is needed to meet current demands

- how much do we expect will be wanted

++++

- what average amount has been produced historically

- how much is available

- how much we project will be produced

*/
