package ship

import (
	"github.com/jhuggett/sea/db"
	"gorm.io/gorm"
)

type Ship struct {
	gorm.Model
	X float64
	Y float64

	WorldMapID uint
}

func New() *Ship {
	return &Ship{}
}

func (s *Ship) Create() (uint, error) {
	err := db.Conn().Create(s).Error
	if err != nil {
		return 0, err
	}

	return s.ID, nil
}

func (s *Ship) Save() error {
	return db.Conn().Save(s).Error
}

func (s *Ship) Move(x, y float64) error {
	s.X = x
	s.Y = y

	err := s.Save()
	if err != nil {
		return err
	}

	return nil
}

func Get(id uint) (*Ship, error) {
	var s Ship
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &s, nil
}
