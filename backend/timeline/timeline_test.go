package timeline

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimelineStartsAndStops(t *testing.T) {
	tl := New()

	triggerCycle := make(chan struct{})

	tl.cycler = func(cycle func(), stop chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case <-triggerCycle:
				cycle()
			}
		}
	}

	tl.Start()

	assert.Equal(t, uint64(0), tl.current)

	tl.Stop()

	assert.Equal(t, uint64(0), tl.current)

	tl.Start()

	triggerCycle <- struct{}{}

	assert.Equal(t, uint64(1), tl.current)

	tl.Stop()

	assert.Equal(t, uint64(1), tl.current)

	tl.Start()

	triggerCycle <- struct{}{}

	assert.Equal(t, uint64(2), tl.current)

	tl.Stop()

	assert.Equal(t, uint64(2), tl.current)
}

func TestTimelineRunsEvents(t *testing.T) {
	tl := New()

	triggerCycle := make(chan struct{})

	tl.cycler = func(cycle func(), stop chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case <-triggerCycle:
				cycle()
			}
		}
	}

	tl.Start()

	eventHappenedCount := 0

	tl.Do(func() uint64 {
		eventHappenedCount++
		return 0
	}, 1)

	assert.Equal(t, 1, tl.queue.Len())
	assert.Zero(t, eventHappenedCount)

	triggerCycle <- struct{}{}

	assert.Eventually(t, func() bool {
		return eventHappenedCount == 1
	}, time.Second, time.Second/4, "eventHappenedCount should be 1")

	assert.Eventually(t, func() bool {
		return tl.queue.Len() == 0
	}, time.Second, time.Second/4, "queue should be empty")

	tl.Do(func() uint64 {
		eventHappenedCount++
		return 1
	}, 0)

	assert.Equal(t, 1, tl.queue.Len())

	triggerCycle <- struct{}{}

	assert.Eventually(t, func() bool {
		return eventHappenedCount == 2
	}, time.Second, time.Second/4, "eventHappenedCount should be 2")

	assert.Eventually(t, func() bool {
		return tl.queue.Len() == 1
	}, time.Second, time.Second/4, "queue should still have 1 item")

	triggerCycle <- struct{}{}

	assert.Eventually(t, func() bool {
		return eventHappenedCount == 3
	}, time.Second, time.Second/4, "eventHappenedCount should be 3")

	assert.Eventually(t, func() bool {
		return tl.queue.Len() == 1
	}, time.Second, time.Second/4, "queue should still have 1 item")

	tl.Stop()

}
