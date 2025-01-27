package outbound

import (
	"github.com/jhuggett/sea/timeline"
)

type TimeChangedReq struct {
	CurrentTick    uint64 `json:"current_tick"`
	TicksPerSecond uint64 `json:"ticks_per_second"`

	CurrentDay  uint64 `json:"current_day"`
	CurrentYear uint64 `json:"current_year"`

	IsPaused bool `json:"is_paused"`
}

type TimeChangedResp struct{}

func (s *Sender) TimeChanged() error {

	currentTimeline := s.gameContext.Timeline
	currentTick := currentTimeline.CurrentTick()

	_, err := s.rpc.Send("TimeChanged", TimeChangedReq{
		CurrentTick:    currentTick,
		TicksPerSecond: currentTimeline.TicksPerCycle(),
		CurrentDay:     currentTick / timeline.Day,
		CurrentYear:    currentTick / timeline.Year,
		IsPaused:       !currentTimeline.IsRunning,
	})
	if err != nil {
		return err
	}

	return nil
}
