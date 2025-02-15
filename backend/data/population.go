package data

import "gorm.io/gorm"

type Population struct {
	gorm.Model

	// WorldMapID uint
	// WorldMap   WorldMap `gorm:"foreignKey:WorldMapID"`

	// ContinentID uint
	// Continent   Continent `gorm:"foreignKey:ContinentID"`

	EconomyID uint
	Economy   Economy `gorm:"foreignKey:EconomyID"`

	Size uint
}

/*

Need to be able to associate value to item constructs.

Universal concept of value.
E.g. for a population, value of 1 piece of eight could be 1.00.
And the value of a fish could be 0.10. Meaning 1 piece of eight is worth 10 fish.


How does money get valued then? Demand for commerce?



Needs:
- Clothing
- Food
  - Spices
  - Beverages
- Building materials
- Tools
- Weapons
- Medicine
- Books

We need some form of import export representation. Which given that continents produce different
types of things at different rates, their inter trading should be determined by distance (traveling salesmen problem),
resulting in a graph of trade routes.



Variety of goods should also be valued. (Maybe fads too? Like how tulips were once worth a lot in the Netherlands)
*/
