package inbound

import (
	"log/slog"

	"github.com/jhuggett/sea/data/ship"
)

type GetInventoryInfoReq struct {
	ShipID int `json:"ship_id"`
}

type GetInventoryInfoResp struct {
}

// This triggers the ShipChanged event to provide the ship data
func GetInventoryInfo(conn Connection, r GetInventoryInfoReq) (GetInventoryInfoResp, error) {
	resp := GetInventoryInfoResp{}

	ship, err := ship.Get(uint(r.ShipID))
	if err != nil {
		slog.Error("Error getting ship", "err", err)
		return resp, err
	}

	err = conn.Sender().ShipInventoryChanged(uint(r.ShipID), *ship.Inventory())
	if err != nil {
		slog.Error("Error sending ship inventory changed", "err", err)
	}

	return resp, nil
}
