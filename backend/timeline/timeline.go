package timeline

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/jhuggett/sea/utils/callback"
	"github.com/jhuggett/sea/utils/priority_queue"
)

type Timeline struct {
	current Tick

	ticksPerCycle Tick

	queue *priority_queue.PriorityQueue[Event]

	stop   chan struct{}
	cycler func(cycle func(), stop chan struct{})

	IsRunning bool
}

func New() *Timeline {
	return &Timeline{
		current:       0,
		ticksPerCycle: 1,

		cycler: func(cycle func(), stop chan struct{}) {
			for {
				select {
				case <-stop:
					return
				case <-time.After(time.Second):
					cycle()
				}
			}
		},

		queue: priority_queue.New(func(a, b Event) bool {
			return (a).Target < (b).Target
		}),
	}
}

// Stop the timeline.
func (t *Timeline) Stop() {
	if !t.IsRunning {
		return
	}
	slog.Debug("Stopping timeline", "current", t.current)
	close(t.stop)
	t.IsRunning = false
}

// Start the timeline (will resume if previously stopped).
func (t *Timeline) Start() {
	if t.IsRunning {
		return
	}
	t.IsRunning = true
	t.stop = make(chan struct{})

	slog.Debug("Starting timeline", "current", t.current)

	go t.cycler(func() {
		slog.Debug("Cycling", "current", t.current)
		t.current += t.ticksPerCycle
		t.processQueue()
		slog.Debug("Finished cycling", "current", t.current)
	}, t.stop)
}

// Should return how many ticks to wait until it is invoke again. 0 means it will not be invoked again.
type EventDo func() Tick

type Event struct {
	Target   Tick    // The tick when the event should be invoked.
	Enqueued Tick    // The tick when the event was enqueued.
	Do       EventDo // The function to invoke.

	uid string
}

func (e Event) SameAs(other priority_queue.Compareable) bool {
	return e.uid == other.(Event).uid
}

func (e *Event) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("Event{Target: %d, Enqueued: %d}", e.Target, e.Enqueued))
}

func (t *Timeline) processQueue() {
	slog.Debug("Processing queue", "current", t.current, "queue", t.queue.Len())
	event := t.queue.PopIt()

	for event != nil {
		slog.Debug("Processing event", "event", event)
		if event.Target > t.current {
			t.queue.PushIt(*event)
			slog.Debug("Re-enqueued event; passed target", "event", event)
			break
		}

		inTicks := event.Do()
		if inTicks > 0 {
			slog.Debug("Re-enqueued event", "event", event)
			t.queue.PushIt(Event{
				Target:   event.Target + inTicks,
				Do:       event.Do,
				Enqueued: event.Target,
				uid:      event.uid,
			})
		}
		event = t.queue.PopIt()
	}

	slog.Debug("Finished processing queue", "current", t.current, "queue", t.queue.Len())
}

func generateUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Failed to generate UID", "error", err)
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// Do will invoke the do function after inTicks.
func (t *Timeline) Do(do EventDo, afterTicks Tick) func() {
	slog.Debug("Enqueuing event", "current", t.current, "afterTicks", afterTicks)

	e := Event{
		Target:   t.current + afterTicks,
		Do:       do,
		Enqueued: t.current,
		uid:      generateUID(),
	}

	t.queue.PushIt(e)

	return func() {
		slog.Debug("Removing event", "current", t.current, "afterTicks", afterTicks)

		t.queue.RemoveIt(e)
	}
}

// TODO: add actual identifier for timeline
func (t *Timeline) id() int {
	return 0
}

func (t *Timeline) TicksPerCycle() Tick {
	return t.ticksPerCycle
}

func (t *Timeline) CurrentTick() Tick {
	return t.current
}

type TicksPerCycleChangedEvent struct {
	PrevTicksPerCycle Tick
	NewTicksPerCycle  Tick
	CurrentTick       Tick
}

var onTicksPerCycleChanged = callback.NewRegistryMap[TicksPerCycleChangedEvent]()

func (t *Timeline) SetTicksPerCycle(ticksPerCycle Tick) {

	Event := TicksPerCycleChangedEvent{
		PrevTicksPerCycle: t.ticksPerCycle,
		NewTicksPerCycle:  ticksPerCycle,
		CurrentTick:       t.current,
	}

	t.ticksPerCycle = ticksPerCycle

	onTicksPerCycleChanged.Invoke([]any{t.id()}, Event)
}

func (t *Timeline) OnTicksPerCycleChangedDo(do func(TicksPerCycleChangedEvent)) func() {
	return onTicksPerCycleChanged.Register([]any{t.id()}, do)
}

type Tick uint64

const (
	Day   Tick = 2
	Week  Tick = 7 * Day
	Month Tick = 30 * Day
	Year  Tick = 12 * Month
)
