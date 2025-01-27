package inbound

import (
	"encoding/json"

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
}

type GetWorldMapResp struct {
	Continents []*Continent `json:"continents"`
}

func GetWorldMap(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {

		worldMap, err := world_map.Get(conn.Context().GameMapID())
		if err != nil {
			return nil, err
		}

		continents := []*Continent{}

		for _, continentData := range worldMap.Persistent.Continents {
			c := &Continent{
				Points: []Point{},
				Name:   continentData.Name,
			}

			continentModel := continent.Using(*continentData)

			coastalPoints, err := continentModel.CoastalPoints()
			if err != nil {
				return nil, err
			}

			center, _ := coordination.Sort(coastalPoints)

			for _, point := range coastalPoints {
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
}
