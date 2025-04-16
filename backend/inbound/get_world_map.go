package inbound

import (
	"github.com/jhuggett/sea/data/continent"
	"github.com/jhuggett/sea/data/world_map"
	"github.com/jhuggett/sea/utils/coordination"
)

type GetWorldMapReq struct {
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`

	Coastal   bool    `json:"coastal"`
	Elevation float64 `json:"elevation"`
}

type Continent struct {
	Points []Point            `json:"points"`
	Center coordination.Point `json:"center"`

	Name string `json:"name"`

	ID uint `json:"id"`
}

type GetWorldMapResp struct {
	Continents []*Continent `json:"continents"`
}

func GetWorldMap(r GetWorldMapReq, gameMapID uint) (GetWorldMapResp, error) {
	worldMap, err := world_map.Get(gameMapID)
	if err != nil {
		return GetWorldMapResp{}, err
	}

	continents := []*Continent{}

	for _, continentData := range worldMap.Persistent.Continents {
		c := &Continent{
			Points: []Point{},
			Name:   continentData.Name,
			ID:     continentData.ID,
		}

		continentModel := continent.Using(*continentData)

		_, err := continentModel.LoadPoints()
		if err != nil {
			return GetWorldMapResp{}, err
		}

		coastalPoints, err := continentModel.CoastalPoints()
		if err != nil {
			return GetWorldMapResp{}, err
		}

		center, _ := coordination.Sort(coastalPoints)

		for _, point := range continentModel.Persistent.Points {
			c.Points = append(c.Points, Point{
				X:         point.X,
				Y:         point.Y,
				Coastal:   point.Coastal,
				Elevation: point.Elevation,
			})
		}

		c.Center = center

		continents = append(continents, c)
	}

	return GetWorldMapResp{
		Continents: continents,
	}, nil
}
