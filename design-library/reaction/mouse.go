package reaction

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MouseDown ReactionType = "MouseDown"
)

type MouseDownEvent struct {
	X, Y   int
	Button ebiten.MouseButton
}

func NewMouseDownReaction(
	condition func(event MouseDownEvent) bool,
	callback func(event MouseDownEvent),
) Reaction {
	return NewReaction[MouseDownEvent](
		MouseDown,
		condition,
		callback,
	)
}

const (
	MouseUp ReactionType = "MouseUp"
)

type MouseUpEvent struct {
	X, Y   int
	Button ebiten.MouseButton
}

func (e MouseUpEvent) XY() (int, int) {
	return e.X, e.Y
}

func NewMouseUpReaction(
	condition func(event MouseUpEvent) bool,
	callback func(event MouseUpEvent),
) Reaction {
	return NewReaction[MouseUpEvent](
		MouseUp,
		condition,
		callback,
	)
}

const (
	MouseMoved ReactionType = "MouseMoved"
)

type MouseMovedEvent struct {
	X, Y int
}

func (e MouseMovedEvent) XY() (int, int) {
	return e.X, e.Y
}

type PositionedEvent interface {
	XY() (int, int)
}

func NewMouseMovedReaction(
	condition func(event MouseMovedEvent) bool,
	callback func(event MouseMovedEvent),
) Reaction {
	return NewReaction[MouseMovedEvent](
		MouseMoved,
		condition,
		callback,
	)
}

// Mouse Drag

const MouseDrag ReactionType = "MouseDrag"

type OnMouseDragEvent struct {
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
	Button         ebiten.MouseButton
}

func NewMouseDragReaction(
	condition func(event OnMouseDragEvent) bool,
	callback func(event OnMouseDragEvent),
) Reaction {
	return NewReaction[OnMouseDragEvent](
		MouseDrag,
		condition,
		callback,
	)
}

// Mouse Wheel

const MouseWheel ReactionType = "MouseWheel"

type MouseWheelEvent struct {
	YOffset float64
}

func NewMouseWheelReaction(
	condition func(event MouseWheelEvent) bool,
	callback func(event MouseWheelEvent),
) Reaction {
	return NewReaction[MouseWheelEvent](
		MouseWheel,
		condition,
		callback,
	)
}
