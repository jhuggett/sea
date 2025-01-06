package port

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/inventory"
)

type Port struct {
	Persistent models.Port
}

func New() *Port {
	return &Port{}
}

func (s *Port) Create() (uint, error) {
	// Create inventory
	i := models.Inventory{}

	err := db.Conn().Create(&i).Error

	if err != nil {
		return 0, err
	}

	s.Persistent.InventoryID = i.ID

	err = db.Conn().Create(&s.Persistent).Error

	return s.Persistent.ID, nil
}

func (s *Port) Save() error {
	return db.Conn().Save(s.Persistent).Error
}

func Get(id uint) (*Port, error) {
	var s models.Port
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: s,
	}, nil
}

func All(worldMapID uint) ([]*Port, error) {
	var persistedPortData []models.Port
	err := db.Conn().Preload("Point").Preload("Inventory").Preload("Inventory.Items").Where("world_map_id = ?", worldMapID).Find(&persistedPortData).Error
	if err != nil {
		return nil, err
	}

	ports := []*Port{}
	for _, p := range persistedPortData {
		ports = append(ports, &Port{
			Persistent: p,
		})
	}

	return ports, nil
}

func Find(point models.Point) (*Port, error) {
	var port models.Port
	err := db.Conn().Debug().Preload("Point").Preload("Inventory").Preload("Inventory.Items").Where("point_id = ?", point.ID).First(&port).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: port,
	}, nil
}

func (s *Port) Inventory() *inventory.Inventory {
	return &inventory.Inventory{Persistent: s.Persistent.Inventory}
}

func (s *Port) Fetch() (*Port, error) {
	var port models.Port
	err := db.Conn().Preload("Point").Preload("Inventory").Preload("Inventory.Items").First(&port, s.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: port,
	}, nil
}
