package producer

import (
	"strings"

	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/inventory"
)

type Producer struct {
	Persistent models.Producer
}

func Create(producer models.Producer, portID uint) (uint, error) {
	inventoryID, err := inventory.Create()
	if err != nil {
		return 0, err
	}

	producer.InventoryID = inventoryID

	producer.PortID = portID

	err = db.Conn().Create(&producer).Error
	if err != nil {
		return 0, err
	}

	return producer.ID, nil
}

func All(portID uint) ([]*Producer, error) {
	var persistedProducerData []models.Producer
	err := db.Conn().Where("port_id = ?", portID).Find(&persistedProducerData).Error
	if err != nil {
		return nil, err
	}

	producers := []*Producer{}

	for _, p := range persistedProducerData {
		producers = append(producers, &Producer{
			Persistent: p,
		})
	}

	return producers, nil
}

func (p *Producer) Products() []string {
	var products []string
	for _, product := range strings.Split(p.Persistent.Products, ",") {
		products = append(products, product)
	}
	return products
}
