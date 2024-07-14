package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/models/ship"
	"github.com/jhuggett/sea/outbound"
)

type MoveShipReq struct {
	ShipID uint    `json:"ship_id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type MoveShipResp struct {
	Success bool `json:"success"`
}

func MoveShip(send *outbound.Sender) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r MoveShipReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog.Info("MoveShip", "id", r.ShipID, "x", r.X, "y", r.Y)

		s, err := ship.Get(r.ShipID)
		if err != nil {
			return nil, err
		}

		err = s.Move(r.X, r.Y)
		if err != nil {
			return nil, err
		}

		_, err = send.ShipChangedTarget(outbound.ShipChangedTargetReq{
			ShipID: s.ID,
			X:      s.X,
			Y:      s.Y,
		})
		if err != nil {
			return nil, err
		}

		return MoveShipResp{
			Success: true,
		}, nil
	}
}
