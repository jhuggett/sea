package outbound

import (
	"log/slog"
)

type ShipChangedTargetReq struct {
	ShipID uint    `json:"ship_id"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

type ShipChangedTargetResp struct{}

func (s *Sender) ShipChangedTarget(req ShipChangedTargetReq) (interface{}, error) {
	slog.Info("ShipChangedTarget", "id", req.ShipID, "x", req.X, "y", req.Y)

	_, err := s.rpc.Send("ShipChangedTarget", req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
