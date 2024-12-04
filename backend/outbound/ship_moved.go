package outbound

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/coordination"
)

type ShipMovedReq struct {
	ShipID   uint               `json:"ship_id"`
	Location coordination.Point `json:"location"`
}

type ShipMovedResp struct{}

func (s *Sender) ShipMoved(shipId uint, location coordination.Point) error {
	slog.Info("ShipMoved", "id", shipId, "location", location)

	_, err := s.rpc.Send("ShipMoved", ShipMovedReq{
		ShipID:   shipId,
		Location: location,
	})
	if err != nil {
		return err
	}

	return nil
}
