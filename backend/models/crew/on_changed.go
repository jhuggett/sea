package crew

import (
	"log/slog"

	"github.com/jhuggett/sea/utils/callback"
)

type OnChangedEventData struct {
	Crew Crew
}

var onChangedRegistryMap = callback.NewRegistryMap[OnChangedEventData]()

func (i *Crew) OnChangedDo(do func(OnChangedEventData)) func() {
	return onChangedRegistryMap.Register([]any{i.Persistent.ID}, do)
}

func (i *Crew) Changed() {
	slog.Info("Crew changed event", "id", i.Persistent.ID)

	i, err := i.Fetch()

	if err != nil {
		slog.Error("Failed to reload Crew", "id", i.Persistent.ID, "error", err)
		return
	}

	onChangedRegistryMap.Invoke([]any{i.Persistent.ID}, OnChangedEventData{
		Crew: *i,
	})
}
