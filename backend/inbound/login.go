package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/models/ship"
)

type LoginReq struct {
	ShipID uint `json:"ship_id"`
}

type LoginResp struct {
	ShipID  uint `json:"ship_id,omitempty"`
	Success bool `json:"success"`
}

func Login() InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r LoginReq
		if err := json.Unmarshal(req, &r); err != nil {
			slog.Error("Failed to unmarshal", "err", err)
			return nil, err
		}

		slog.Info("Login", "id", r.ShipID)

		if r.ShipID == 0 {
			slog.Info("Creating new ship")
			s := ship.New()
			id, err := s.Create()
			if err != nil {
				return nil, err
			}
			slog.Info("Ship created", "id", id)
			return LoginResp{
				ShipID:  id,
				Success: true,
			}, nil
		}

		s, err := ship.Get(r.ShipID)
		if err != nil {
			slog.Error("Ship not found", "id", r.ShipID)
			return nil, err
		}

		slog.Info("Ship found", "id", s.ID)

		return LoginResp{
			ShipID:  s.ID,
			Success: true,
		}, nil
	}
}
