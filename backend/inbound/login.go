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

		slog.Info("Ship found", "id", s.Persistent.ID)

		s.OnDockedDo(func(data ship.DockedEventData) {
			slog.Info("Ship docked", "id", s.Persistent.ID)
			conn.Sender().ShipDocked(s.Persistent.ID, data.Location, false)
		})

		s.OnUndockedDo(func(data ship.UnDockedEventData) {
			slog.Info("Ship undocked", "id", s.Persistent.ID)
			conn.Sender().ShipDocked(s.Persistent.ID, data.Location, true)
		})

		// s.OnMovedDo(func(data ship.ShipMovedEventData) {
		// 	slog.Info("Ship moved", "id", s.Persistent.ID)
		// 	conn.Sender().ShipMoved(s.Persistent.ID, data.Location)
		// })

		return LoginResp{
			Ship:    ShipInfo{ID: s.Persistent.ID, X: s.Persistent.X, Y: s.Persistent.Y},
			Success: true,
		}, nil
	}
}
