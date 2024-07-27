package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/models/world_map"
)

type RegisterReq struct {
}

type RegisterResp struct {
	GameCtx game_context.Snapshot `json:"snapshot"`
}

func Register() InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r RegisterReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog.Info("Register")

		// create world map
		worldMap := world_map.New()
		worldMapID, err := worldMap.Create()
		if err != nil {
			return nil, err
		}

		// numberOfHightPoints := 20

		// minX := 0
		// maxX := 30
		// minY := 0
		// maxY := 30

		// for i := 0; i < numberOfHightPoints; i++ {
		// 	// randomize x and y and elevation
		// 	x := rand.Float64()*float64(maxX-minX) + float64(minX)
		// 	y := rand.Float64()*float64(maxY-minY) + float64(minY)
		// 	elevation := (rand.Float64()*50 + 50) / 100

		// 	err = worldMap.AddHighPoint(x, y, elevation)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }

		worldMap.GenerateCoasts()

		// create ship
		ship := ship.New()
		ship.WorldMapID = worldMapID
		shipID, err := ship.Create()
		if err != nil {
			return nil, err
		}

		return RegisterResp{
			GameCtx: game_context.Snapshot{
				ShipID:    shipID,
				GameMapID: worldMapID,
			},
		}, nil
	}
}
