package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
	"github.com/jhuggett/sea/utils/coordination"
)

type UnDockedEvent struct {
	Location coordination.Point
}

var undockedRegistryMap = callback.NewRegistryMap[UnDockedEvent]()

func (s *Ship) OnUndockedDo(do func(UnDockedEvent)) func() {
	return undockedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) Undocked() {
	s.Persistent.IsDocked = false
	err := s.Save()
	if err != nil {
		slog.Error("Failed to save ship", "id", s.Persistent.ID, "error", err)
		return
	}
	undockedRegistryMap.Invoke([]any{s.Persistent.ID}, UnDockedEvent{
		Location: s.Location(),
	})
}
