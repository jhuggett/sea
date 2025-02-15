package data

import (
	"github.com/jhuggett/sea/constructs/items"
	"gorm.io/gorm"
)

type Economy struct {
	gorm.Model

	PortID uint
	Port   Port `gorm:"foreignKey:PortID"`

	Markets []Market

	Populations []Population

	Industries []Industry
}

// The workforce of an economy is a percentage of the sum of all associated populations
// Industries takes shares of the workforce

func HowManyUnitsAreWanted(product items.ItemType) {

	// lookup market for product

	// get current demand per population

}
