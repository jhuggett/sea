package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/continent"
	"github.com/jhuggett/sea/models/crew"
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
		worldMap := &world_map.WorldMap{
			Persistent: models.WorldMap{},
		}
		worldMapID, err := worldMap.Create()
		if err != nil {
			return nil, err
		}

		worldMap.Generate()

		worldMap, err = world_map.Get(worldMapID) // reload
		if err != nil {
			return nil, err
		}

		for _, c := range worldMap.Persistent.Continents {
			// create ports
			port := port.New()
			port.Persistent.WorldMapID = worldMapID

			continent := continent.Continent{Persistent: *c}

			port.Persistent.PointID = continent.GetCoastalPoints()[0].ID
			_, err = port.Create()
			if err != nil {
				return nil, err
			}
		}

		// create ship
		ship := ship.New()
		// ship.Persistent.Coffers = 1000
		ship.Persistent.WorldMapID = worldMapID
		shipID, err := ship.Create()
		if err != nil {
			return nil, err
		}

		err = ship.Inventory().AddItem(models.Item{
			Name:   string(models.ItemTypePieceOfEight),
			Amount: 1000,
		})
		if err != nil {
			return nil, err
		}

		// create crew
		crew := crew.New(models.Crew{
			ShipID:  shipID,
			Size:    1,
			Wage:    1,
			Rations: 1,
			Morale:  1,
		})
		_, err = crew.Create()
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
