package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
	"github.com/jhuggett/sea/utils/coordination"
)

type ShipMovedEvent struct {
	Location coordination.Point
}

var shipMovedRegistryMap = callback.NewRegistryMap[ShipMovedEvent]()

func (s *Ship) OnMovedDo(do func(ShipMovedEvent)) func() {
	return shipMovedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) Moved() {
	slog.Info("Ship moved event", "id", s.Persistent.ID)
	shipMovedRegistryMap.Invoke([]any{s.Persistent.ID}, ShipMovedEvent{
		Location: s.Location(),
	})
}
