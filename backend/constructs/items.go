package constructs

type Item struct {
	Name string

	// How much space this item takes up
	SpacePerItem float32

	// How much weight each item has
	WeightPerItem float32

	RationValue float32

	Composition map[ItemType]float32
}

type ItemType string

const (
	PieceOfEight ItemType = "piece_of_eight"
	Fish         ItemType = "fish"
	Grain        ItemType = "grain"
	Silver       ItemType = "silver"
)

/*

Base items define their weight relative to how much space they take up. (e.g. silver)

Composite items do not define their weight. They define their composition of base items and their space
per item which allows the weight to be calculated based on the composition and space.

*/

var itemConstructs = map[ItemType]Item{
	Silver: {
		Name:          "Silver",
		SpacePerItem:  1,
		WeightPerItem: 2,
		RationValue:   0,
	},
	PieceOfEight: {
		Name:         "Piece of Eight",
		SpacePerItem: 0.01,
		RationValue:  0,
		Composition: map[ItemType]float32{
			Silver: 1,
		},
	},
	Fish: {
		Name:          "Fish",
		SpacePerItem:  1,
		WeightPerItem: 1,
		RationValue:   1,
	},
	Grain: {
		Name:          "Grain",
		SpacePerItem:  1,
		WeightPerItem: 1,
		RationValue:   .5,
	},
}

func LookupItem(name string) Item {

	itemType := ItemType(name)

	return itemConstructs[itemType]
}

func (i Item) Weight() float32 {
	if i.Composition == nil || len(i.Composition) == 0 {
		return i.WeightPerItem
	}

	var weight float32
	for itemType, amount := range i.Composition {
		weight += itemConstructs[itemType].Weight() * (amount * i.SpacePerItem)
	}

	return weight
}
