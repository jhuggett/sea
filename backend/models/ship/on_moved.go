package ship

import (
	"log/slog"

	"github.com/jhuggett/sea/models/world_map"
	"github.com/jhuggett/sea/utils/callback"
)

type ShipMovedEventData struct {
	Location world_map.Point
}

var shipMovedRegistryMap = callback.NewRegistryMap[ShipMovedEventData]()

func (s *Ship) OnMovedDo(do func(ShipMovedEventData)) func() {
	return shipMovedRegistryMap.Register([]any{s.ID}, do)
}

func (s *Ship) Moved() {
	slog.Info("Ship moved event", "id", s.ID)
	shipMovedRegistryMap.Invoke([]any{s.ID}, ShipMovedEventData{
		Location: s.Location(),
	})
}
