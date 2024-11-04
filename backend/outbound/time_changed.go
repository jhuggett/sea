package outbound

import "log/slog"

type TimeChangedReq struct {
	CurrentTick    uint64 `json:"current_tick"`
	TicksPerSecond uint64 `json:"ticks_per_second"`
}

type TimeChangedResp struct{}

func (s *Sender) TimeChanged(currentTick uint64, ticksPerSecond uint64) error {
	slog.Info("TimeChanged", "current_tick", currentTick, "ticks_per_second", ticksPerSecond)

	_, err := s.rpc.Send("TimeChanged", TimeChangedReq{
		CurrentTick:    currentTick,
		TicksPerSecond: ticksPerSecond,
	})
	if err != nil {
		return err
	}

	return nil
}
