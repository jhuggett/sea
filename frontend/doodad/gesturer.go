package doodad

import (
	"errors"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type CallbackRegistryRecord[T any] struct {
	Callback T
	UUID     string
}

type CallbackRegistry[T any] struct {
	callbacks []CallbackRegistryRecord[T]
}

func (r *CallbackRegistry[T]) Add(callback T) string {
	uuid := uuid.NewString()
	r.callbacks = append(r.callbacks, CallbackRegistryRecord[T]{Callback: callback, UUID: uuid})
	return uuid
}

func (r *CallbackRegistry[T]) Remove(uuid string) {
	for i, cb := range r.callbacks {
		if cb.UUID == uuid {
			r.callbacks = append(r.callbacks[:i], r.callbacks[i+1:]...)
			break
		}
	}
}

var ErrStopPropagation = errors.New("stop propagation")
var ErrUnregister = errors.New("unregister")

func (r *CallbackRegistry[T]) InvokeEndToStart(do func(T) error) {
	for i := len(r.callbacks) - 1; i >= 0; i-- {
		cb := r.callbacks[i]
		err := do(cb.Callback)
		if err != nil {
			if errors.Is(err, ErrStopPropagation) {
				return
			}
			if errors.Is(err, ErrUnregister) {
				r.Remove(cb.UUID)
			}
		}
	}
}

func (r *CallbackRegistry[T]) Register(callback T) func() {
	uuid := r.Add(callback)
	return func() {
		r.Remove(uuid)
	}
}

type Press struct {
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
	Button         ebiten.MouseButton
}

type MouseUpEvent struct {
	X, Y   int
	Button ebiten.MouseButton
}

// Return indicates whether the event should continue to propagate
type OnMouseUpFunc func(MouseUpEvent) error

type OnMouseDragFunc func(lastX, lastY, currentX, currentY int) error

type OnMouseWheelFunc func(offset float64) error

type OnMouseMoveFunc func(x, y int) error

type Gesturer interface {
	OnMouseUp(OnMouseUpFunc) func()
	OnMouseDrag(OnMouseDragFunc) func()
	OnMouseWheel(OnMouseWheelFunc) func()
	OnMouseMove(OnMouseMoveFunc) func()
	Update()

	Parent() Gesturer
	CreateChild() Gesturer
}

type gesturer struct {
	OnMouseUpCallbacks    CallbackRegistry[OnMouseUpFunc]
	OnMouseDragCallbacks  CallbackRegistry[OnMouseDragFunc]
	OnMouseWheelCallbacks CallbackRegistry[OnMouseWheelFunc]
	OnMouseMoveCallbacks  CallbackRegistry[OnMouseMoveFunc]

	MouseX int
	MouseY int

	Press *Press

	parent Gesturer
}

func NewGesturer() *gesturer {
	return &gesturer{
		OnMouseUpCallbacks:    CallbackRegistry[OnMouseUpFunc]{},
		OnMouseDragCallbacks:  CallbackRegistry[OnMouseDragFunc]{},
		OnMouseWheelCallbacks: CallbackRegistry[OnMouseWheelFunc]{},
		OnMouseMoveCallbacks:  CallbackRegistry[OnMouseMoveFunc]{},
	}
}

func (g *gesturer) Parent() Gesturer {
	if g.parent == nil {
		return nil
	}
	return g.parent
}

func (g *gesturer) CreateChild() Gesturer {
	child := NewGesturer()
	child.parent = g
	return child
}

// Register callbacks

func (g *gesturer) OnMouseUp(callback OnMouseUpFunc) func() {
	return g.OnMouseUpCallbacks.Register(callback)
}

func (g *gesturer) OnMouseDrag(callback OnMouseDragFunc) func() {
	return g.OnMouseDragCallbacks.Register(callback)
}

func (g *gesturer) OnMouseWheel(callback OnMouseWheelFunc) func() {
	return g.OnMouseWheelCallbacks.Register(callback)
}

func (g *gesturer) OnMouseMove(callback OnMouseMoveFunc) func() {
	return g.OnMouseMoveCallbacks.Register(callback)
}

// Update logic

func (g *gesturer) Update() {
	x, y := ebiten.CursorPosition()

	if x != g.MouseX || y != g.MouseY {
		g.OnMouseMoveCallbacks.InvokeEndToStart(func(ommf OnMouseMoveFunc) error {
			return ommf(x, y)
		})
	}

	g.MouseX = x
	g.MouseY = y

	_, yoff := ebiten.Wheel()
	if yoff != 0 {
		g.OnMouseWheelCallbacks.InvokeEndToStart(func(omwf OnMouseWheelFunc) error {
			return omwf(yoff)
		})
	}

	var pressedMouseButton ebiten.MouseButton = -1
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		pressedMouseButton = ebiten.MouseButtonLeft
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		pressedMouseButton = ebiten.MouseButtonRight
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		pressedMouseButton = ebiten.MouseButtonMiddle
	}

	if pressedMouseButton != -1 {
		if g.Press == nil {
			g.Press = &Press{
				StartX:    x,
				StartY:    y,
				X:         x,
				Y:         y,
				TimeStart: time.Now(),
				Button:    pressedMouseButton,
			}
		}

		if time.Since(g.Press.TimeStart) > 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-x)) > 25 || math.Abs(float64(g.Press.StartY-y)) > 25) {
			g.OnMouseDragCallbacks.InvokeEndToStart(func(omdf OnMouseDragFunc) error {
				return omdf(g.Press.X, g.Press.Y, x, y)
			})
		}

		g.Press.X = x
		g.Press.Y = y

	} else {
		if g.Press != nil {
			if time.Since(g.Press.TimeStart) < 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-g.Press.X)) < 8 && math.Abs(float64(g.Press.StartY-g.Press.Y)) < 8) {
				slog.Info("Click", "x", g.Press.X, "y", g.Press.Y)
				g.OnMouseUpCallbacks.InvokeEndToStart(func(omuf OnMouseUpFunc) error {
					return omuf(MouseUpEvent{
						X:      g.Press.X,
						Y:      g.Press.Y,
						Button: g.Press.Button,
					})
				})
			}
			g.Press = nil
		}
	}
}
