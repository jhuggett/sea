package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/utils/callback"
)

type DockedEventData struct {
	Location world_map.Point
}

var dockedRegistryMap = callback.NewRegistryMap[DockedEventData]()

func (s *Ship) OnDockedDo(do func(DockedEventData)) func() {
	return dockedRegistryMap.Register([]any{s.ID}, do)
}

func (s *Ship) Docked() {
	s.IsDocked = true
	err := s.Save()
	if err != nil {
		slog.Error("Failed to save ship", "id", s.ID, "error", err)
		return
	}
	dockedRegistryMap.Invoke([]any{s.ID}, DockedEventData{
		Location: s.Location(),
	})
}
