package ship

import "github.com/jhuggett/sea/utils/callback"

type AnchorLoweredLocation string

const (
	AnchorLoweredLocationUnknown AnchorLoweredLocation = "unknown"
	AnchorLoweredLocationPort    AnchorLoweredLocation = "port"
	AnchorLoweredLocationOpenSea AnchorLoweredLocation = "open_sea"
)

type AnchorLoweredEvent struct {
	Location AnchorLoweredLocation
}

var anchorLoweredRegistryMap = callback.NewRegistryMap[AnchorLoweredEvent]()

func (s *Ship) OnAnchorLoweredDo(do func(AnchorLoweredEvent)) func() {
	return anchorLoweredRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) AnchorLowered(e AnchorLoweredEvent) {
	anchorLoweredRegistryMap.Invoke([]any{s.Persistent.ID}, e)
}
