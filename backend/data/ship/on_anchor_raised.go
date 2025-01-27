package ship

import (
	"github.com/jhuggett/sea/utils/callback"
)

type AnchorRaisedEvent struct {
}

var anchorRaisedRegistryMap = callback.NewRegistryMap[AnchorRaisedEvent]()

func (s *Ship) OnAnchorRaisedDo(do func(AnchorRaisedEvent)) func() {
	return anchorRaisedRegistryMap.Register([]any{s.Persistent.ID}, do)
}

func (s *Ship) AnchorRaised() {
	anchorRaisedRegistryMap.Invoke([]any{s.Persistent.ID}, AnchorRaisedEvent{})
}
