package outbound

import (
	"github.com/jhuggett/sea/timeline"
)

type TimeChangedReq struct {
	CurrentTick    timeline.Tick `json:"current_tick"`
	TicksPerSecond timeline.Tick `json:"ticks_per_second"`

	CurrentDay  timeline.Tick `json:"current_day"`
	CurrentYear timeline.Tick `json:"current_year"`

	IsPaused bool `json:"is_paused"`
}

type TimeChangedResp struct{}

func (s *Sender) TimeChanged() error {

	currentTimeline := s.gameContext.Timeline
	currentTick := currentTimeline.CurrentTick()

	_, err := s.Receiver.OnTimeChanged(TimeChangedReq{
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
