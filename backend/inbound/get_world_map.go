package inbound

import (
	"encoding/json"

	"github.com/jhuggett/sea/models/world_map"
)

type GetWorldMapReq struct {
}

type CoastalPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Continent struct {
	CoastalPoints []CoastalPoint `json:"coastal_points"`
}

type GetWorldMapResp struct {
	/// HighPoints []HighPoint `json:"high_points"`

	Continents []Continent `json:"continents"`
}

func GetWorldMap(conn connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {

		worldMap, err := world_map.Get(conn.Context().GameMapID)
		if err != nil {
			return nil, err
		}

		continents := []Continent{}

		for _, continent := range worldMap.Continents {
			c := Continent{
				CoastalPoints: []CoastalPoint{},
			}

			for _, coastalPoint := range continent.CoastalPoints {
				cp := CoastalPoint{
					X: coastalPoint.X,
					Y: coastalPoint.Y,
				}

				c.CoastalPoints = append(c.CoastalPoints, cp)
			}

			continents = append(continents, c)
		}

		// highPoints, err := worldMap.GetHighPoints()
		// if err != nil {
		// 	return nil, err
		// }

		// highPointsData := make([]HighPoint, 0, len(highPoints))

		// for _, highPoint := range highPoints {
		// 	highPointsData = append(highPointsData, HighPoint{
		// 		X:         highPoint.X,
		// 		Y:         highPoint.Y,
		// 		Elevation: highPoint.Elevation,
		// 	})
		// }

		return GetWorldMapResp{
			Continents: continents,
		}, nil
	}
}
