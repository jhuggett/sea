package session

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
)

type Session struct {
	Persistent data.Session
}

func All() ([]data.Session, error) {
	var sessionData []data.Session
	err := db.Conn().Find(&sessionData).Error
	if err != nil {
		return nil, err
	}
	return sessionData, nil
}
