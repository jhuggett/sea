package inbound

import (
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

func ControlTime(conn Connection, req ControlTimeReq) (ControlTimeResp, error) {

	if req.SetTicksPerSecondTo != nil {
		slog.Info("SetTicksPerSecondTo", "to", *req.SetTicksPerSecondTo)
		conn.Context().Timeline.SetTicksPerCycle(*req.SetTicksPerSecondTo)
	} else if req.SetTicksPerSecondBy != nil {
		slog.Info("SetTicksPerSecondBy", "by", *req.SetTicksPerSecondBy)
		conn.Context().Timeline.SetTicksPerCycle(conn.Context().Timeline.TicksPerCycle() + *req.SetTicksPerSecondBy)
	} else if req.Pause {
		slog.Info("Pause")
		conn.Context().Timeline.Stop()
	} else if req.Resume {
		slog.Info("Resume")
		conn.Context().Timeline.Start()
	}

	conn.Sender().TimeChanged()

	respObj := ControlTimeResp{}

	return respObj, nil

}
