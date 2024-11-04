package timeline

import (
	"log/slog"
	"time"

	"github.com/jhuggett/sea/utils/callback"
)

var eventStore = map[string]Event{}

func StoreEvent(event Event, id string) {
	eventStore[id] = event
}

func GetEvent(id string) (Event, bool) {
	event, ok := eventStore[id]
	return event, ok
}

type Timeline struct {
	current uint64

	stop chan struct{}

	ticksPerSecond uint64

	continualEvents []ContinualEvent
}

func New() *Timeline {
	return &Timeline{
		current:        0,
		stop:           make(chan struct{}),
		ticksPerSecond: 1,
	}
}

type Event interface {
}

type ContinualEvent interface {
	Do(ticks uint64) (stop bool)
}

// type TargetedEvent interface {
// }

var OnTicksPerSecondChanged = callback.NewRegistry[struct {
	Old       uint64
	New       uint64
	TickCount uint64
}]()

func (t *Timeline) SetTicksPerSecond(ticksPerSecond uint64) {
	t.ticksPerSecond = ticksPerSecond

	OnTicksPerSecondChanged.Invoke(struct {
		Old       uint64
		New       uint64
		TickCount uint64
	}{
		Old:       t.ticksPerSecond,
		New:       ticksPerSecond,
		TickCount: t.current,
	})
}

func (t *Timeline) TicksPerSecond() uint64 {
	return t.ticksPerSecond
}

func (t *Timeline) Tick() uint64 {
	t.current += t.ticksPerSecond
	return t.ticksPerSecond
}

func (t *Timeline) RegisterContinualEvent(event ContinualEvent) {
	t.continualEvents = append(t.continualEvents, event)
}

func (t *Timeline) run() {
	for {
		select {
		case <-t.stop:
			slog.Info("Timeline stopped")
			// should save here in future etc
			return
		case <-time.After(1 * time.Second):
			elapsedTicks := t.Tick()
			t.processContinualEvents(elapsedTicks)
		}
	}
}

func (t *Timeline) Stop() {
	close(t.stop)
}

func (t *Timeline) Start() {
	go t.run()
}

func (t *Timeline) processContinualEvents(ticks uint64) {

	// seems like continual events aren't working...
	// eg. ship moving only seems to move once and then hangs or something

	var nextContinualEvents []ContinualEvent
	for _, event := range t.continualEvents {
		if !event.Do(ticks) {
			// t.continualEvents = append(t.continualEvents[:i], t.continualEvents[i+1:]...)
			nextContinualEvents = append(nextContinualEvents, event)
		}
	}
	t.continualEvents = nextContinualEvents
}

/*
Eg

Move ship
every tick, move ship based on speed


Events
- Continual events (pass number of ticks as time passes until it says to stop)
- Targeted events (trigger when total ticks reaches the target tick. Could be repeatable)
*/
