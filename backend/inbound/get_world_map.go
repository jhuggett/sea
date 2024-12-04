package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/utils/coordination"
)

type GetWorldMapReq struct {
}

type CoastalPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Continent struct {
	CoastalPoints []CoastalPoint     `json:"coastal_points"`
	Center        coordination.Point `json:"center"`
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
				CoastalPoints: []CoastalPoint{},
			}

			center, s := coordination.Sort(continent.CoastalPoints)

			slog.Info("Center", "x", center.X, "y", center.Y, "coastal_points", s)

			for _, coastalPoint := range s {
				cp := CoastalPoint{
					X: coastalPoint.X,
					Y: coastalPoint.Y,
				}

				c.CoastalPoints = append(c.CoastalPoints, cp)
			}

			c.Center = center

			continents = append(continents, c)
		}

		return GetWorldMapResp{
			Continents: continents,
		}, nil
	}
}
