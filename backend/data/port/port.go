package port

import (
	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/building"
	"github.com/jhuggett/sea/data/inventory"
	"github.com/jhuggett/sea/db"
)

type Port struct {
	Persistent data.Port
}

func (s *Port) Create() (uint, error) {
	// Create inventory
	i := data.Inventory{}

	err := db.Conn().Create(&i).Error

	if err != nil {
		return 0, err
	}

	s.Persistent.InventoryID = i.ID

	err = db.Conn().Create(&s.Persistent).Error
	if err != nil {
		return 0, err
	}

	// -- Create buildings

	building.Create(s.Persistent.ID, "Shipyard", string(building.TypeShipyard), 1, 1)

	building.Create(s.Persistent.ID, "Tavern", string(building.TypeTavern), 3, 2)

	building.Create(s.Persistent.ID, "Market", string(building.TypeMarket), 5, 6)

	// --

	return s.Persistent.ID, nil
}

func (s *Port) Save() error {
	return db.Conn().Save(s.Persistent).Error
}

func Get(id uint) (*Port, error) {
	var s data.Port
	err := db.Conn().First(&s, id).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: s,
	}, nil
}

func All(worldMapID uint) ([]*Port, error) {
	var portData []data.Port
	err := db.Conn().
		Preload("Point").
		Preload("Inventory").
		Preload("Inventory.Items").
		Joins("join continents on continents.id=ports.continent_id").
		Where("continents.world_map_id = ?", worldMapID).
		Find(&portData).Error
	if err != nil {
		return nil, err
	}

	ports := []*Port{}
	for _, p := range portData {
		ports = append(ports, &Port{
			Persistent: p,
		})
	}

	return ports, nil
}

func Find(point data.Point) (*Port, error) {
	var port data.Port
	err := db.Conn().Debug().
		Preload("Point").
		Preload("Inventory").
		Preload("Inventory.Items").
		Where("point_id = ?", point.ID).
		First(&port).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: port,
	}, nil
}

func (s *Port) Inventory() *inventory.Inventory {
	return &inventory.Inventory{Persistent: s.Persistent.Inventory}
}

func (s *Port) Fetch() (*Port, error) {
	var port data.Port
	err := db.Conn().Preload("Point").Preload("Inventory").Preload("Inventory.Items").First(&port, s.Persistent.ID).Error
	if err != nil {
		return nil, err
	}

	return &Port{
		Persistent: port,
	}, nil
}

func (s *Port) Continent() (*data.Continent, error) {
	var continent data.Continent
	err := db.Conn().Where("id = ?", s.Persistent.Point.ContinentID).First(&continent).Error
	if err != nil {
		return nil, err
	}

	return &continent, nil
}
