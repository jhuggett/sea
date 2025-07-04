package inbound

import (
	"log/slog"

	"github.com/jhuggett/sea/data/ship"
)

type RepairShipReq struct {
	ShipID uint `json:"ship_id"`
}

type RepairShipResp struct {
}

func RepairShip(r RepairShipReq) (RepairShipResp, error) {
	slog.Info("Repairing ship")

	ship, err := ship.Get(r.ShipID)
	if err != nil {
		return RepairShipResp{}, err
	}

	ship.Persistent.StateOfRepair = 1
	err = ship.Save()
	if err != nil {
		return RepairShipResp{}, err
	}

	return RepairShipResp{}, nil
}
