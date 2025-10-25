package fleet

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
)

type Fleet struct {
	Persistent data.Fleet
}

func (s *Fleet) Create() (uint, error) {
	err := db.Conn().Create(&s.Persistent).Error
	if err != nil {
		return 0, err
	}

	return s.Persistent.ID, nil
}

func Get(id uint) (*Fleet, error) {
	var s data.Fleet
	err := db.Conn().Joins("Ships").First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Fleet{
		Persistent: s,
	}, nil
}

func GetFleetByShipID(id uint) (*Fleet, error) {
	var s data.Fleet
	err := db.Conn().
		Preload("Ships").
		Joins("join ships on ships.fleet_id=fleets.id").
		Where("ships.id = ?", id).
		First(&s).Error
	if err != nil {
		return nil, err
	}

	return &Fleet{
		Persistent: s,
	}, nil
}
