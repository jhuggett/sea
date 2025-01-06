package outbound

import (
	"log/slog"

	"github.com/jhuggett/sea/timeline"
)

type TimeChangedReq struct {
	CurrentTick    uint64 `json:"current_tick"`
	TicksPerSecond uint64 `json:"ticks_per_second"`

	CurrentDay  uint64 `json:"current_day"`
	CurrentYear uint64 `json:"current_year"`
}

type TimeChangedResp struct{}

func (s *Sender) TimeChanged(currentTick uint64, ticksPerSecond uint64) error {
	slog.Info("TimeChanged", "current_tick", currentTick, "ticks_per_second", ticksPerSecond)

	_, err := s.rpc.Send("TimeChanged", TimeChangedReq{
		CurrentTick:    currentTick,
		TicksPerSecond: ticksPerSecond,
		CurrentDay:     currentTick / timeline.Day,
		CurrentYear:    currentTick / timeline.Year,
	})
	if err != nil {
		return err
	}

	return nil
}
