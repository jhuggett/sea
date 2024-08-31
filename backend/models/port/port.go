package port

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models/world_map"
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	CoastalPointID uint
	WorldMapID     uint

	CoastalPoint *world_map.CoastalPoint `gorm:"foreignKey:CoastalPointID"`
}

func New() *Port {
	return &Port{}
}

func (s *Port) Create() (uint, error) {
	err := db.Conn().Create(s).Error
	if err != nil {
		return 0, err
	}

	return s.ID, nil
}

func (s *Port) Save() error {
	return db.Conn().Save(s).Error
}

func Get(id uint) (*Port, error) {
	var s Port
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func All(worldMapID uint) ([]*Port, error) {
	var ports []*Port
	err := db.Conn().Preload("CoastalPoint").Where("world_map_id = ?", worldMapID).Find(&ports).Error
	if err != nil {
		return nil, err
	}

	return ports, nil
}
