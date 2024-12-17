package item

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
)

type Item struct {
	Persistent models.Item
}

func Create(item models.Item) (uint, error) {
	err := db.Conn().Create(&item).Error
	if err != nil {
		return 0, err
	}

	return item.ID, nil
}

func (i *Item) Save() error {
	return db.Conn().Save(&i.Persistent).Error
}

func (i *Item) Delete() error {
	return db.Conn().Delete(&i.Persistent).Error
}
