package ships

// Crew size range for best operation of the ship
// MinimumSafeManning uint
// MaximumSafeManning uint

// // Recommended cargo capacity is how much cargo weight the ship can carry without slowing down
// RecommendedMaxCargoWeightCapacity uint

// // Max cargo capacity is how much cargo the ship can carry until there just isn't any more space
// MaxCargoSpaceCapacity uint

// StateOfRepair float64

// // Squares moved per day
// BaseSpeed float64

// WorldMapID uint
// WorldMap   *WorldMap `gorm:"foreignKey:WorldMapID"`

// IsDocked bool

// InventoryID uint
// Inventory   *Inventory

// Name string

type Type string

const (
	TypeFishingBoat Type = "fishing_boat"
	TypeSloop       Type = "sloop"
	TypeFrig        Type = "frigate"
)

// Number of Masts
// Number of Sails
// Number of Long Cannons
// Number of Short Cannons
// Number of Swivel Cannons
// Number of Mortars
// Number of Decks

// Figurehead

// Amount of space and how it's divided up (e.g. 60% cargo, 20% crew, 5% for food (kitchen etc), blacksmith, carpenter, medical, captains quarters, entertainment, etc)

// Color of sails
// Color of hull
// Color of trim
// Color of figurehead
// Color can be paint or the wood, or material (e.g. red dyed sails, gold painted trim, and a bronze figurehead)

// Overall workmanship quality (poor, average, good, excellent, masterwork)
// State of repair (cannon holes, torn sails, rotting wood, leaks, burn damage, etc)
// Age

// IDEA!!!!!!!!! BELOW

// Careening
// Need to clean bottom of ship every few months to remove barnacles and other growth that slows the ship down considerably (copper plating the hull stops this, but is very very expensive)
// Can take weeks to do, and the ship is out of commission during this time and exposed to attack (or take it to a port and pay for it to be done faster, assuming the port welcomes you)

// Earthquake and Tsunamis that damage ports
