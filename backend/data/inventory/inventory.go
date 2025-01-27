package inventory

import (
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/constructs/items"
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/item"
	"github.com/jhuggett/sea/db"
)

type Inventory struct {
	Persistent data.Inventory
}

func (i *Inventory) Items() []item.Item {
	var items []item.Item
	for _, it := range i.Persistent.Items {
		items = append(items, item.Item{Persistent: *it})
	}
	return items
}

func (i *Inventory) AddItem(it data.Item) error {
	it.InventoryID = i.Persistent.ID

	defer i.Changed()

	existingTypeOfSameName := i.FindItem(it.Name)

	if existingTypeOfSameName != nil {
		existingTypeOfSameName.Persistent.Amount += it.Amount

		return existingTypeOfSameName.Save()
	}

	_, err := item.Create(it)

	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (i *Inventory) RemoveItem(it data.Item) error {
	defer i.Changed()

	existingTypeOfSameName := i.FindItem(it.Name)

	if existingTypeOfSameName == nil {
		return fmt.Errorf("item not found")
	}

	if existingTypeOfSameName.Persistent.Amount < it.Amount {
		return fmt.Errorf("not enough items")
	}

	existingTypeOfSameName.Persistent.Amount -= it.Amount

	if existingTypeOfSameName.Persistent.Amount == 0 {
		return existingTypeOfSameName.Delete()
	}

	return existingTypeOfSameName.Save()
}

func (i *Inventory) FindItem(name string) *item.Item {

	var itemData *data.Item

	err := db.Conn().Where("inventory_id = ? AND name = ?", i.Persistent.ID, name).First(&itemData).Error

	if err != nil {
		slog.Debug("Error finding item", "err", err)
		return nil
	}

	return &item.Item{Persistent: *itemData}

	// return nil
	// for _, it := range i.Persistent.Items {
	// 	if it.Name == name {
	// 		return &item.Item{Persistent: *it}
	// 	}
	// }
	// return nil
}

// Returns a new Inventory struct using the pointer receiver's ID
func (i *Inventory) Fetch() (*Inventory, error) {
	var inv data.Inventory
	err := db.Conn().Preload("Items").First(&inv, i.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return &Inventory{Persistent: inv}, nil
}

func Fetch(id uint) (*Inventory, error) {
	var inv data.Inventory
	x := db.Conn().Debug().Preload("Items").First(&inv, id)

	err := x.Error

	if err != nil {
		return nil, err
	}

	return &Inventory{Persistent: inv}, nil
}

func Create() (uint, error) {
	inv := data.Inventory{}
	err := db.Conn().Create(&inv).Error
	if err != nil {
		return 0, err
	}

	return inv.ID, nil
}

func (i *Inventory) Rations() ([]item.Item, error) {
	var itemData []data.Item
	err := db.Conn().Where("inventory_id = ? AND marked_as_ration = ?", i.Persistent.ID, true).Find(&itemData).Error
	if err != nil {
		return nil, err
	}

	items := make([]item.Item, len(itemData))
	for i, it := range itemData {
		items[i] = item.Item{Persistent: it}
	}

	return items, nil
}

func (i *Inventory) TotalWeight() float32 {
	var total float32
	for _, it := range i.Persistent.Items {

		itemConstruct := items.LookupItem(it.Name)

		total += itemConstruct.Weight() * it.Amount
	}
	return total
}

func (i *Inventory) OccupiedSpace() float32 {
	var total float32
	for _, it := range i.Persistent.Items {

		itemConstruct := items.LookupItem(it.Name)

		total += itemConstruct.SpacePerItem * it.Amount
	}
	return total
}
