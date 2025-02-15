package industry

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
)

type Industry struct {
	Persistent data.Industry
}

func (i *Industry) Create() (uint, error) {
	err := db.Conn().Create(&i.Persistent).Error
	return i.Persistent.ID, err
}
