package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/data/ship"
)

type RepairShipReq struct {
	ShipID uint `json:"ship_id"`
}

type RepairShipResp struct {
}

func RepairShip(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		// _, err := conn.Context().RepairShip()
		// if err != nil {
		// 	return nil, err
		// }

		slog.Info("Repairing ship")

		var reqObj RepairShipReq
		if err := json.Unmarshal(req, &reqObj); err != nil {
			return nil, err
		}

		ship, err := ship.Get(reqObj.ShipID)
		if err != nil {
			return nil, err
		}

		ship.Persistent.StateOfRepair = 1
		err = ship.Save()
		if err != nil {
			return nil, err
		}

		return RepairShipResp{}, nil
	}
}
