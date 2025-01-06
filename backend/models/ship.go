package models

import "gorm.io/gorm"

// Need to move records up to models, then the model level methods use the types from models

type Ship struct {
	gorm.Model
	X float64
	Y float64

	// Crew size range for best operation of the ship
	MinimumSafeManning uint
	MaximumSafeManning uint

	// Recommended cargo capacity is how much cargo weight the ship can carry without slowing down
	RecommendedMaxCargoWeightCapacity uint

	// Max cargo capacity is how much cargo the ship can carry until there just isn't any more space
	MaxCargoSpaceCapacity uint

	StateOfRepair float64

	// Squares moved per day
	BaseSpeed float64

	WorldMapID uint
	WorldMap   *WorldMap `gorm:"foreignKey:WorldMapID"`

	IsDocked bool

	InventoryID uint
	Inventory   *Inventory
}
