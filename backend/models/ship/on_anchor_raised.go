package ship

import (
	"github.com/jhuggett/sea/utils/callback"
)

type AnchorRaisedEventData struct {
}

var anchorRaisedRegistryMap = callback.NewRegistryMap[AnchorRaisedEventData]()

func (s *Ship) OnAnchorRaisedDo(do func(AnchorRaisedEventData)) func() {
	return anchorRaisedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) AnchorRaised() {
	anchorRaisedRegistryMap.Invoke([]any{s.Persistent.ID}, AnchorRaisedEventData{})
}
