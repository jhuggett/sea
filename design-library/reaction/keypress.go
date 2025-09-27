package reaction

import "github.com/hajimehoshi/ebiten/v2"

const (
	KeyDown ReactionType = "KeyDown"
)

type KeyDownEvent struct {
	Key ebiten.Key
	*Event
}

func (e *KeyDownEvent) setEvent(event *Event) {
	e.Event = event
}

func NewKeyDownReaction(
	condition func(event *KeyDownEvent) bool,
	callback func(event *KeyDownEvent),
) Reaction {
	return NewReaction[*KeyDownEvent](
		KeyDown,
		condition,
		callback,
	)
}

func SpecificKeyDown(key ebiten.Key) func(event *KeyDownEvent) bool {
	return func(event *KeyDownEvent) bool {
		return event.Key == key
	}
}
