package ship

import "github.com/jhuggett/sea/utils/callback"

type AnchorLoweredLocation string

const (
	AnchorLoweredLocationUnknown AnchorLoweredLocation = "unknown"
	AnchorLoweredLocationPort    AnchorLoweredLocation = "port"
	AnchorLoweredLocationOpenSea AnchorLoweredLocation = "open_sea"
)

type AnchorLoweredEventData struct {
	Location AnchorLoweredLocation
}

var anchorLoweredRegistryMap = callback.NewRegistryMap[AnchorLoweredEventData]()

func (s *Ship) OnAnchorLoweredDo(do func(AnchorLoweredEventData)) func() {
	return anchorLoweredRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) AnchorLowered(data AnchorLoweredEventData) {
	anchorLoweredRegistryMap.Invoke([]any{s.Persistent.ID}, data)
}
