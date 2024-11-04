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
			conn.Context().Timeline.SetTicksPerSecond(*reqObj.SetTicksPerSecondTo)
		} else if reqObj.SetTicksPerSecondBy != nil {
			slog.Info("SetTicksPerSecondBy", "by", *reqObj.SetTicksPerSecondBy)
			conn.Context().Timeline.SetTicksPerSecond(conn.Context().Timeline.TicksPerSecond() + *reqObj.SetTicksPerSecondBy)
		}

		respObj := ControlTimeResp{}

		return respObj, nil
	}
}
