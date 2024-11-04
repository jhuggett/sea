package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/utils/callback"
)

type UnDockedEventData struct {
	Location world_map.Point
}

var undockedRegistryMap = callback.NewRegistryMap[UnDockedEventData]()

func (s *Ship) OnUndockedDo(do func(UnDockedEventData)) func() {
	return undockedRegistryMap.Register([]any{s.ID}, do)
}

func (s *Ship) Undocked() {
	s.IsDocked = false
	err := s.Save()
	if err != nil {
		slog.Error("Failed to save ship", "id", s.ID, "error", err)
		return
	}
	undockedRegistryMap.Invoke([]any{s.ID}, UnDockedEventData{
		Location: s.Location(),
	})
}
