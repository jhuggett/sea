package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
	"github.com/jhuggett/sea/utils/coordination"
)

type UnDockedEventData struct {
	Location coordination.Point
}

var undockedRegistryMap = callback.NewRegistryMap[UnDockedEventData]()

func (s *Ship) OnUndockedDo(do func(UnDockedEventData)) func() {
	return undockedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) Undocked() {
	s.Persistent.IsDocked = false
	err := s.Save()
	if err != nil {
		slog.Error("Failed to save ship", "id", s.Persistent.ID, "error", err)
		return
	}
	undockedRegistryMap.Invoke([]any{s.Persistent.ID}, UnDockedEventData{
		Location: s.Location(),
	})
}
