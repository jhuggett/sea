package inbound

import (
	"encoding/json"
	"log/slog"
)

type ControlTimeReq struct {
	SetTicksPerSecondTo *uint64 `json:"set_ticks_per_second_to,omitempty"`
	SetTicksPerSecondBy *uint64 `json:"set_ticks_per_second_by,omitempty"`
}

type ControlTimeResp struct {
}

func ControlTime(conn Connection) InboundFunc {
	return func(req json.RawMessage) (interface{}, error) {
		var reqObj ControlTimeReq
		if err := json.Unmarshal(req, &reqObj); err != nil {
			return nil, err
		}

		if reqObj.SetTicksPerSecondTo != nil {
			slog.Info("SetTicksPerSecondTo", "to", *reqObj.SetTicksPerSecondTo)
			conn.Context().Timeline.SetTicksPerCycle(*reqObj.SetTicksPerSecondTo)
		} else if reqObj.SetTicksPerSecondBy != nil {
			slog.Info("SetTicksPerSecondBy", "by", *reqObj.SetTicksPerSecondBy)
			conn.Context().Timeline.SetTicksPerCycle(conn.Context().Timeline.TicksPerCycle() + *reqObj.SetTicksPerSecondBy)
		}

		respObj := ControlTimeResp{}

		return respObj, nil
	}
}
