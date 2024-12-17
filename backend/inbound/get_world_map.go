package inbound

import (
	"encoding/json"

	"github.com/jhuggett/sea/models/world_map"
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

		for _, continent := range worldMap.Persistent.Continents {
			c := &Continent{
				Points: []Point{},
			}

			center, s := coordination.Sort(continent.Points)

			for _, point := range s {
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
