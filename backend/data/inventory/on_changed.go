package inventory

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
)

type OnChangedEvent struct {
	Inventory Inventory
}

var onChangedRegistryMap = callback.NewRegistryMap[OnChangedEvent]()

func (i *Inventory) OnChangedDo(do func(OnChangedEvent)) func() {
	return onChangedRegistryMap.Register([]any{i.Persistent.ID}, do)
}

func (i *Inventory) Changed() {
	i, err := i.Fetch()

	if err != nil {
		slog.Error("Failed to reload inventory", "id", i.Persistent.ID, "error", err)
		return
	}

	onChangedRegistryMap.Invoke([]any{i.Persistent.ID}, OnChangedEvent{
		Inventory: *i,
	})
}
