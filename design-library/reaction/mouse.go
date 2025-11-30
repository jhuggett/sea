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
	*Event
}

func (e MouseDownEvent) XY() (int, int) {
	return e.X, e.Y
}

func (e *MouseDownEvent) setEvent(event *Event) {
	e.Event = event
}

func NewMouseDownReaction(
	condition func(event *MouseDownEvent) bool,
	callback func(event *MouseDownEvent),
) Reaction {
	return NewReaction[*MouseDownEvent](
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
	*Event
}

func (e *MouseUpEvent) setEvent(event *Event) {
	e.Event = event
}

func (e MouseUpEvent) XY() (int, int) {
	return e.X, e.Y
}

func NewMouseUpReaction(
	condition func(event *MouseUpEvent) bool,
	callback func(event *MouseUpEvent),
) Reaction {
	return NewReaction[*MouseUpEvent](
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
	*Event
}

func (e *MouseMovedEvent) setEvent(event *Event) {
	e.Event = event
}

func (e MouseMovedEvent) XY() (int, int) {
	return e.X, e.Y
}

type PositionedEvent interface {
	XY() (int, int)
}

func NewMouseMovedReaction(
	condition func(event *MouseMovedEvent) bool,
	callback func(event *MouseMovedEvent),
) Reaction {
	return NewReaction[*MouseMovedEvent](
		MouseMoved,
		condition,
		callback,
	)
}

// Mouse Drag

const MouseDrag ReactionType = "MouseDrag"

type OnMouseDragEvent struct {
	OrignX, OrignY int
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
	Button         ebiten.MouseButton
	*Event
}

func (o *OnMouseDragEvent) XY() (int, int) {
	return o.OrignX, o.OrignY
}

func (e *OnMouseDragEvent) setEvent(event *Event) {
	e.Event = event
}

func NewMouseDragReaction(
	condition func(event *OnMouseDragEvent) bool,
	callback func(event *OnMouseDragEvent),
) Reaction {
	return NewReaction[*OnMouseDragEvent](
		MouseDrag,
		condition,
		callback,
	)
}

// Mouse Wheel

const MouseWheel ReactionType = "MouseWheel"

type MouseWheelEvent struct {
	OriginX, OriginY int
	YOffset          float64
	*Event
}

func (m *MouseWheelEvent) XY() (int, int) {
	return m.OriginX, m.OriginY
}

func (e *MouseWheelEvent) setEvent(event *Event) {
	e.Event = event
}

func NewMouseWheelReaction(
	condition func(event *MouseWheelEvent) bool,
	callback func(event *MouseWheelEvent),
) Reaction {
	return NewReaction[*MouseWheelEvent](
		MouseWheel,
		condition,
		callback,
	)
}
