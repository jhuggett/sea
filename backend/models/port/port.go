package port

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
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
	if err != nil {
		return 0, err
	}

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
	err := db.Conn().Preload("Point").Where("world_map_id = ?", worldMapID).Find(&persistedPortData).Error
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
	err := db.Conn().Preload("Point").Where("point_id = ?", point.ID).Find(&port).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: port,
	}, nil
}
