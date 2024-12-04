package ship

import (
	"fmt"

	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/crew"
	"github.com/jhuggett/sea/utils/coordination"
)

type Ship struct {
	Persistent models.Ship
}

func New() *Ship {
	return &Ship{}
}

func (s *Ship) Create() (uint, error) {
	err := db.Conn().Create(&s.Persistent).Error
	if err != nil {
		return 0, err
	}

	return s.Persistent.ID, nil
}

func (s *Ship) Save() error {
	return db.Conn().Save(s.Persistent).Error
}

func (s *Ship) Move(x, y float64) error {
	s.Persistent.X = x
	s.Persistent.Y = y

	err := s.Save()
	if err != nil {
		return err
	}

	return nil
}

func Get(id uint) (*Ship, error) {
	var s models.Ship
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Ship{
		Persistent: s,
	}, nil
}

func (s *Ship) Location() coordination.Point {
	return coordination.Point{X: int(s.Persistent.X), Y: int(s.Persistent.Y)}
}

func (s *Ship) Crew() (*crew.Crew, error) {
	return crew.Where(models.Crew{ShipID: s.Persistent.ID})
}

func (s *Ship) SubtractFromCoffers(amount float64) error {
	if s.Persistent.Coffers < amount {
		return fmt.Errorf("insufficient funds")
	}

	s.Persistent.Coffers -= amount
	return s.Save()
}
