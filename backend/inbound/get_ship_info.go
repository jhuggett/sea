package inbound

import "log/slog"

type GetShipInfoReq struct {
	ShipID int `json:"ship_id"`
}

type GetShipInfoResp struct {
}

// This triggers the ShipChanged event to provide the ship data
func GetShipInfo(conn Connection, r GetShipInfoReq) (GetShipInfoResp, error) {
	resp := GetShipInfoResp{}

	err := conn.Sender().ShipChanged(uint(r.ShipID))
	if err != nil {
		slog.Error("Error sending ship changed", "err", err)
	}

	return resp, nil
}
