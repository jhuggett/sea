package inbound

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type HireCrewReq struct {
	Size int `json:"size"`
}

type HireCrewResp struct {
}

func HireCrew(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var r HireCrewReq
		if err := json.Unmarshal(req, &r); err != nil {
			return nil, err
		}

		slog.Info("HireCrew")

		ship, err := conn.Context().Ship()
		if err != nil {
			return nil, fmt.Errorf("failed to get ship: %w", err)
		}

		crew, err := ship.Crew()
		if err != nil {
			return nil, fmt.Errorf("failed to get crew: %w", err)
		}

		crew.Persistent.Size += r.Size

		if err := crew.Save(); err != nil {
			return nil, fmt.Errorf("failed to save crew: %w", err)
		}

		return HireCrewResp{}, nil
	}
}
