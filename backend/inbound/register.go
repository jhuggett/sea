package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jhuggett/sea/data"
	"github.com/jhuggett/sea/data/continent"
	"github.com/jhuggett/sea/data/crew"
	"github.com/jhuggett/sea/data/fleet"
	"github.com/jhuggett/sea/data/population"
	"github.com/jhuggett/sea/data/port"
	"github.com/jhuggett/sea/data/ship"
	"github.com/jhuggett/sea/data/world_map"
	"github.com/jhuggett/sea/db"
	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/name"
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
			Persistent: data.WorldMap{},
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

		for _, continentData := range worldMap.Persistent.Continents {

			continentModel := continent.Using(*continentData)

			popData, err := continentModel.Population()
			if err != nil {
				return nil, err
			}

			popModel := population.Using(popData)

			// create ports
			port := &port.Port{
				Persistent: data.Port{
					Name:         "Port of " + continentModel.Persistent.Name,
					PopulationID: popModel.Persistent.ID,
				},
			}
			port.Persistent.WorldMapID = worldMapID

			coastalPoints, err := continentModel.CoastalPoints()
			if err != nil {
				return nil, err
			}

			port.Persistent.PointID = coastalPoints[0].ID
			_, err = port.Create()
			if err != nil {
				return nil, err
			}

			port, err = port.Fetch()

			// err = port.Inventory().AddItem(data..Item{
			// 	Name:   string(constructs.PieceOfEight),
			// 	Amount: 1000,
			// })
			// if err != nil {
			// 	return nil, fmt.Errorf("failed to add item: %w", err)
			// }

			// for i := 0; i < 3; i++ {
			// 	_, err := producer.Create(data..Producer{
			// 		Products: strings.Join([]string{
			// 			string(constructs.Grain),
			// 		}, ""),
			// 	}, portID)
			// 	if err != nil {
			// 		return nil, fmt.Errorf("failed to create producer: %w", err)
			// 	}
			// }

		}

		// create fleet
		fleetData := data.Fleet{}
		fleetModel := fleet.Fleet{
			Persistent: fleetData,
		}
		fleetID, err := fleetModel.Create()
		if err != nil {
			return nil, err
		}

		// create ship
		shipData := data.Ship{
			WorldMapID: worldMapID,

			MinimumSafeManning: 2,
			MaximumSafeManning: 10,

			StateOfRepair: 1.0,

			BaseSpeed: 1.0,

			RecommendedMaxCargoWeightCapacity: 300,
			MaxCargoSpaceCapacity:             50,

			Name: "The Ship",

			FleetID: fleetID,
		}

		ship := &ship.Ship{
			Persistent: shipData,
		}

		shipID, err := ship.Create()
		if err != nil {
			return nil, err
		}

		fleetModel.Persistent.CapitalShipID = shipID
		err = db.Conn().Save(&fleetModel.Persistent).Error
		if err != nil {
			return nil, fmt.Errorf("failed to save fleet: %w", err)
		}

		// err = ship.Inventory().AddItem(data..Item{
		// 	Name:   string(constructs.PieceOfEight),
		// 	Amount: 1000,
		// })
		// if err != nil {
		// 	return nil, err
		// }

		// create crew
		crew := crew.Crew{
			Persistent: data.Crew{
				ShipID: shipID,
			},
		}
		_, err = crew.Create()
		if err != nil {
			return nil, err
		}

		// create player
		player := &data.Person{
			FirstName: name.Generate(2),
			LastName:  name.Generate(4),
			NickName:  name.GenerateNickName(),
		}
		err = db.Conn().Create(player).Error
		if err != nil {
			return nil, err
		}

		termsData := &data.EmploymentTerms{
			Title:     "Captain",
			StartDate: 0,
		}

		err = db.Conn().Create(termsData).Error
		if err != nil {
			return nil, err
		}

		contractData := &data.Contract{
			OfferorID:   crew.Persistent.ID,
			OfferorKind: data.CrewContractParty,

			OffereeKind: data.PersonContractParty,
			OffereeID:   player.ID,

			TermsKind: data.EmploymentContractTerms,
			TermsID:   termsData.ID,
		}

		err = db.Conn().Create(contractData).Error
		if err != nil {
			return nil, err
		}

		deedData := data.Deed{
			OwnerID:   player.ID,
			OwnerType: data.DeedOwnerPerson,

			ObjectID:   shipID,
			ObjectType: data.DeedObjectShip,
		}

		err = db.Conn().Create(&deedData).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create deed: %w", err)
		}

		return RegisterResp{
			GameCtx: game_context.Snapshot{
				ShipID:    shipID,
				GameMapID: worldMapID,
			},
		}, nil
	}
}
