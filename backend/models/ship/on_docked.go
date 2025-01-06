package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
	"github.com/jhuggett/sea/utils/coordination"
)

type DockedEventData struct {
	Location coordination.Point
}

var dockedRegistryMap = callback.NewRegistryMap[DockedEventData]()

func (s *Ship) OnDockedDo(do func(DockedEventData)) func() {
	return dockedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) Docked(location coordination.Point) {
	slog.Info("Ship docked", "id", s.Persistent.ID)
	s.Persistent.IsDocked = true
	err := s.Save()
	if err != nil {
		slog.Error("Failed to save ship", "id", s.Persistent.ID, "error", err)
		return
	}
	dockedRegistryMap.Invoke([]any{s.Persistent.ID}, DockedEventData{
		Location: location,
	})
}
