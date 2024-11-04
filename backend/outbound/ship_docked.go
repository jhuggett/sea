package outbound

import (
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
)

type ShipDockedReq struct {
	ShipID   uint            `json:"ship_id,omitempty"`
	Location world_map.Point `json:"location,omitempty"`
	Undocked bool            `json:"undocked,omitempty"`
}

type ShipDockedResp struct{}

func (s *Sender) ShipDocked(shipId uint, location world_map.Point, undocked bool) error {
	slog.Info("ShipDocked", "id", shipId, "location", location)

	_, err := s.rpc.Send("ShipDocked", ShipDockedReq{
		ShipID:   shipId,
		Location: location,
		Undocked: undocked,
	})
	if err != nil {
		return err
	}

	return nil
}
