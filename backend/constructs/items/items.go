package items

import "github.com/jhuggett/sea/timeline"

type Uses string

const (
	// EdibleGrain Uses = "edible_grain"
	// EdibleMeat  Uses = "edible_meat"

	Food Uses = "food"

	Fuel Uses = "fuel"

	Currency Uses = "currency"
)

type Item struct {
	Name string

	// How much space this item takes up
	SpacePerItem float32

	// How much weight each item has
	WeightPerItem float32

	RationValue float32

	Composition map[ItemType]float32

	// E.g. how long it takes for a crop to grow
	BaseTimeToProduce float64 // in days

	// E.g. how much time it takes to do the work necessary to grow a crop
	WorkTimeToProduce float64 // in days

	/*
		In other words BaseTimeToProduce can't be sped up by adding more workers, but WorkTimeToProduce can.
	*/

	// Purposes the item provides. The demand of these usages determine it's value
	Uses map[Uses]float64
}

type ItemType string

const (
	PieceOfEight ItemType = "piece_of_eight"
	Fish         ItemType = "fish"
	Grain        ItemType = "grain"
	Silver       ItemType = "silver"

	Wood ItemType = "wood"
)

/*

Base items define their weight relative to how much space they take up. (e.g. silver)

Composite items do not define their weight. They define their composition of base items and their space
per item which allows the weight to be calculated based on the composition and space.

*/

var itemConstructs = map[ItemType]Item{
	Wood: {
		Name:          "Wood",
		SpacePerItem:  1,
		WeightPerItem: 1,
		RationValue:   0,

		BaseTimeToProduce: float64(timeline.Day) / 3,
		WorkTimeToProduce: float64(timeline.Week * 2),

		Uses: map[Uses]float64{
			Fuel: 1,
		},
	},
	Silver: {
		Name:          "Silver",
		SpacePerItem:  1,
		WeightPerItem: 2,
		RationValue:   0,

		BaseTimeToProduce: 0,
		WorkTimeToProduce: float64(timeline.Week * 2),

		Uses: map[Uses]float64{
			Currency: 1,
		},
	},
	PieceOfEight: {
		Name:         "Piece of Eight",
		SpacePerItem: 0.01,
		RationValue:  0,
		Composition: map[ItemType]float32{
			Silver: 1,
		},

		BaseTimeToProduce: float64(timeline.Day) / 1,
		WorkTimeToProduce: float64(timeline.Day) / 100,

		Uses: map[Uses]float64{
			Currency: 1,
		},
	},
	Fish: {
		Name:          "Fish",
		SpacePerItem:  1,
		WeightPerItem: 1,
		RationValue:   1,

		BaseTimeToProduce: float64(timeline.Day) / 2,
		WorkTimeToProduce: 0,

		Uses: map[Uses]float64{
			Food: 1,
		},
	},
	Grain: {
		Name:          "Grain",
		SpacePerItem:  1,
		WeightPerItem: 1,
		RationValue:   .5,

		BaseTimeToProduce: float64(timeline.Day) / 3,
		WorkTimeToProduce: float64(timeline.Week * 2),

		Uses: map[Uses]float64{
			Food: 1,
		},
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
