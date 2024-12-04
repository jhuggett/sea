package crew

import (
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/models"
)

type Crew struct {
	Persistent models.Crew
}

func New(data models.Crew) *Crew {
	return &Crew{
		Persistent: data,
	}
}

func (s *Crew) Create() (uint, error) {
	err := db.Conn().Create(&s.Persistent).Error
	if err != nil {
		return 0, err
	}

	return s.Persistent.ID, nil
}

func (s *Crew) Save() error {
	return db.Conn().Save(s.Persistent).Error
}

func Get(id uint) (*Crew, error) {
	var s models.Crew
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Crew{
		Persistent: s,
	}, nil
}

func Where(data models.Crew) (*Crew, error) {
	var s models.Crew
	err := db.Conn().Where(&data).First(&s).Error
	if err != nil {
		return nil, err
	}

	return &Crew{
		Persistent: s,
	}, nil
}
