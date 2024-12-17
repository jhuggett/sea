package inventory

import (
	"fmt"

	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/item"
)

type Inventory struct {
	Persistent models.Inventory
}

func (i *Inventory) Items() []item.Item {
	var items []item.Item
	for _, it := range i.Persistent.Items {
		items = append(items, item.Item{Persistent: *it})
	}
	return items
}

func (i *Inventory) AddItem(it models.Item) error {
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

func (i *Inventory) RemoveItem(it models.Item) error {
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
	for _, it := range i.Persistent.Items {
		if it.Name == name {
			return &item.Item{Persistent: *it}
		}
	}
	return nil
}

// Returns a new Inventory struct using the pointer receiver's ID
func (i *Inventory) Fetch() (*Inventory, error) {
	var inv models.Inventory
	err := db.Conn().Preload("Items").First(&inv, i.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return &Inventory{Persistent: inv}, nil
}
