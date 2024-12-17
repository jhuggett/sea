package inventory

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
)

type OnChangedEventData struct {
	Inventory Inventory
}

var onChangedRegistryMap = callback.NewRegistryMap[OnChangedEventData]()

func (i *Inventory) OnChangedDo(do func(OnChangedEventData)) func() {
	return onChangedRegistryMap.Register([]any{i.Persistent.ID}, do)
}

func (i *Inventory) Changed() {
	slog.Info("Inventory changed event", "id", i.Persistent.ID)

	i, err := i.Fetch()

	if err != nil {
		slog.Error("Failed to reload inventory", "id", i.Persistent.ID, "error", err)
		return
	}

	onChangedRegistryMap.Invoke([]any{i.Persistent.ID}, OnChangedEventData{
		Inventory: *i,
	})
}
