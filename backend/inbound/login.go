package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/game_context"
	"github.com/jhuggett/sea/models/ship"
)

type LoginReq struct {
	Snapshot game_context.Snapshot `json:"snapshot"`
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

func Login(setGameContext func(snapshot game_context.Snapshot) Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r LoginReq
		if err := json.Unmarshal(req, &r); err != nil {
			slog.Error("Failed to unmarshal", "err", err)
			return nil, err
		}

		conn := setGameContext(r.Snapshot)

		ctx := conn.Context()

		s, err := ctx.Ship()
		if err != nil {
			slog.Error("Ship not found", "id", ctx.ShipID())
			return nil, err
		}

		slog.Info("Ship found", "id", s.ID)

		s.OnDockedDo(func(data ship.DockedEventData) {
			slog.Info("Ship docked", "id", s.ID)
			conn.Sender().ShipDocked(s.ID, data.Location, false)
		})

		s.OnUndockedDo(func(data ship.UnDockedEventData) {
			slog.Info("Ship undocked", "id", s.ID)
			conn.Sender().ShipDocked(s.ID, data.Location, true)
		})

		s.OnMovedDo(func(data ship.ShipMovedEventData) {
			slog.Info("Ship moved", "id", s.ID)
			conn.Sender().ShipMoved(s.ID, data.Location)
		})

		return LoginResp{
			Ship:    ShipInfo{ID: s.ID, X: s.X, Y: s.Y},
			Success: true,
		}, nil
	}
}
