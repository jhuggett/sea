package inbound

import (
	"encoding/json"
	"log/slog"

	"github.com/jhuggett/sea/timeline"
)

type ControlTimeReq struct {
	SetTicksPerSecondTo *timeline.Tick `json:"set_ticks_per_second_to,omitempty"`
	SetTicksPerSecondBy *timeline.Tick `json:"set_ticks_per_second_by,omitempty"`

	Pause  bool `json:"pause,omitempty"`
	Resume bool `json:"resume,omitempty"`
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
		} else if reqObj.Pause {
			slog.Info("Pause")
			conn.Context().Timeline.Stop()
		} else if reqObj.Resume {
			slog.Info("Resume")
			conn.Context().Timeline.Start()
		}

		conn.Sender().TimeChanged()

		respObj := ControlTimeResp{}

		return respObj, nil
	}
}
