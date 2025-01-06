package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jhuggett/sea/constructs"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/models"
	"github.com/jhuggett/sea/models/continent"
	"github.com/jhuggett/sea/models/crew"
	"github.com/jhuggett/sea/models/port"
	"github.com/jhuggett/sea/models/producer"
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

			coastalPoints, err := continent.CoastalPoints()
			if err != nil {
				return nil, err
			}

			port.Persistent.PointID = coastalPoints[0].ID
			portID, err := port.Create()
			if err != nil {
				return nil, err
			}

			port, err = port.Fetch()

			err = port.Inventory().AddItem(models.Item{
				Name:   string(constructs.PieceOfEight),
				Amount: 1000,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to add item: %w", err)
			}

			for i := 0; i < 3; i++ {
				_, err := producer.Create(models.Producer{
					Products: strings.Join([]string{
						string(constructs.Grain),
					}, ""),
				}, portID)
				if err != nil {
					return nil, fmt.Errorf("failed to create producer: %w", err)
				}
			}

		}

		// create ship
		shipData := models.Ship{
			WorldMapID: worldMapID,

			MinimumSafeManning: 2,
			MaximumSafeManning: 10,

			StateOfRepair: 1.0,

			BaseSpeed: 1.0,

			RecommendedMaxCargoWeightCapacity: 300,
			MaxCargoSpaceCapacity:             50,
		}

		ship := &ship.Ship{
			Persistent: shipData,
		}

		shipID, err := ship.Create()
		if err != nil {
			return nil, err
		}

		err = ship.Inventory().AddItem(models.Item{
			Name:   string(constructs.PieceOfEight),
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
