package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/models/ship"
)

type GameContext struct {
	ShipID    uint
	GameMapID uint

	// Sign this for auth in future
}

type LoginReq struct {
	GameCtx GameContext `json:"context"`
}

type ShipInfo struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	ID uint    `json:"id"`
}
type LoginResp struct {
	Ship    ShipInfo `json:"ship"`
	Success bool     `json:"success"`
}

func Login(setGameContext func(gameCtx GameContext)) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r LoginReq
		if err := json.Unmarshal(req, &r); err != nil {
			slog.Error("Failed to unmarshal", "err", err)
			return nil, err
		}

		setGameContext(r.GameCtx)

		s, err := ship.Get(r.GameCtx.ShipID)
		if err != nil {
			slog.Error("Ship not found", "id", r.GameCtx.ShipID)
			return nil, err
		}

		slog.Info("Ship found", "id", s.ID)

		return LoginResp{
			Ship:    ShipInfo{ID: s.ID, X: s.X, Y: s.Y},
			Success: true,
		}, nil
	}
}
