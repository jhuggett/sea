package ship

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/crew"
	"github.com/jhuggett/sea/models/inventory"
	"github.com/jhuggett/sea/utils/coordination"
)

type Ship struct {
	Persistent models.Ship
}

func (s *Ship) Create() (uint, error) {
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

	// reload

	err = db.Conn().Preload("Inventory").Preload("Inventory.Items").First(&s.Persistent, s.Persistent.ID).Error
	if err != nil {
		return 0, err
	}

	return s.Persistent.ID, nil
}

func (s *Ship) Save() error {
	return db.Conn().Save(&s.Persistent).Error
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
	err := db.Conn().Preload("Inventory").Preload("Inventory.Items").First(&s, id).Error
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

func (s *Ship) Inventory() *inventory.Inventory {

	if s.Persistent.Inventory != nil {
		return &inventory.Inventory{Persistent: *s.Persistent.Inventory}
	}

	panic("Inventory not loaded for ship")

	return nil
}

func (s *Ship) SailingSpeed() (float64, error) {

	crew, err := s.Crew()
	if err != nil {
		return 0, err
	}

	speed := s.Persistent.BaseSpeed

	if uint(crew.Persistent.Size) < s.Persistent.MinimumSafeManning {
		speed = speed * (1 - (float64(crew.Persistent.Size) / float64(s.Persistent.MinimumSafeManning)))
	} else if uint(crew.Persistent.Size) > s.Persistent.MaximumSafeManning {
		// TOOD: add speed reduction for overmanning
		speed = speed
	}

	speed = speed * s.Persistent.StateOfRepair

	// add cargo speed reduction

	return speed, nil
}

func (s *Ship) Fetch() (*Ship, error) {
	err := db.Conn().Preload("Inventory").Preload("Inventory.Items").First(&s.Persistent, s.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
