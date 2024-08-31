package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/models/port"
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

		worldMap.GenerateCoasts()

		worldMap, err = world_map.Get(worldMapID) // reload
		if err != nil {
			return nil, err
		}

		// create ports
		port := port.New()
		port.WorldMapID = worldMapID
		port.CoastalPointID = worldMap.Continents[0].CoastalPoints[0].ID
		_, err = port.Create()
		if err != nil {
			return nil, err
		}

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
