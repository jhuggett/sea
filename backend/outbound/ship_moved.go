package outbound

import (
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
)

type ShipMovedReq struct {
	ShipID   uint            `json:"ship_id"`
	Location world_map.Point `json:"location"`
}

type ShipMovedResp struct{}

func (s *Sender) ShipMoved(shipId uint, location world_map.Point) error {
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
