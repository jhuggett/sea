package building

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/db"
)

type Type string

const (
	TypeShipyard Type = "shipyard"
	TypeTavern   Type = "tavern"
	TypeMarket   Type = "market"
)

type Building struct {
	Persistent data.Building
}

func Create(portID uint, name string, buildingType string, x, y int) (uint, error) {
	persistent := data.Building{
		PortID: portID,
		Name:   name,
		Type:   buildingType,
		X:      x,
		Y:      y,
	}

	err := db.Conn().Create(&persistent).Error
	if err != nil {
		return 0, err
	}

	return persistent.ID, nil
}

func Where(portId uint) ([]Building, error) {
	var buildingsData []data.Building

	err := db.Conn().Where("port_id = ?", portId).Find(&buildingsData).Error
	if err != nil {
		return nil, err
	}

	var buildings []Building

	for _, buildingData := range buildingsData {
		buildings = append(buildings, Building{
			Persistent: buildingData,
		})
	}

	return buildings, nil
}
