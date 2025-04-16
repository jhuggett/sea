package doodad

import (
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Press struct {
	StartX, StartY int
	X, Y           int
	TimeStart      time.Time
}

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

func (r *CallbackRegistry[T]) Call(do func(T) (stop bool)) {
	for i := len(r.callbacks) - 1; i >= 0; i-- {
		cb := r.callbacks[i]
		if do(cb.Callback) {
			return
		}
	}
}

func (r *CallbackRegistry[T]) Register(callback T) func() {
	uuid := r.Add(callback)
	return func() {
		r.Remove(uuid)
	}
}

// Return indicates whether the event should continue to propagate
type OnMouseUpFunc func(x, y int) bool

type OnMouseDragFunc func(lastX, lastY, currentX, currentY int) bool

type OnMouseWheelFunc func(offset float64) bool

type OnMouseMoveFunc func(x, y int) bool

type Gesturer interface {
	OnMouseUp(OnMouseUpFunc) func()
	OnMouseDrag(OnMouseDragFunc) func()
	OnMouseWheel(OnMouseWheelFunc) func()
	OnMouseMove(OnMouseMoveFunc) func()
	Update()
}

type gesturer struct {
	OnMouseUpCallbacks    CallbackRegistry[OnMouseUpFunc]
	OnMouseDragCallbacks  CallbackRegistry[OnMouseDragFunc]
	OnMouseWheelCallbacks CallbackRegistry[OnMouseWheelFunc]
	OnMouseMoveCallbacks  CallbackRegistry[OnMouseMoveFunc]

	MouseX int
	MouseY int

	Press *Press
}

func NewGesturer() *gesturer {
	return &gesturer{
		OnMouseUpCallbacks:    CallbackRegistry[OnMouseUpFunc]{},
		OnMouseDragCallbacks:  CallbackRegistry[OnMouseDragFunc]{},
		OnMouseWheelCallbacks: CallbackRegistry[OnMouseWheelFunc]{},
		OnMouseMoveCallbacks:  CallbackRegistry[OnMouseMoveFunc]{},
	}
}

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

func (g *gesturer) Update() {
	x, y := ebiten.CursorPosition()

	if x != g.MouseX || y != g.MouseY {
		g.OnMouseMoveCallbacks.Call(func(ocf OnMouseMoveFunc) (stop bool) {
			return ocf(x, y)
		})
	}

	g.MouseX = x
	g.MouseY = y

	_, yoff := ebiten.Wheel()
	// if yoff > 0 {

	// 	// need a zoom callback

	// 	// w.Camera.ZoomFactor += .1
	// 	// slog.Debug("Zooming in", "zoom", w.Camera.ZoomFactor)
	// }
	// if yoff < 0 {
	// 	// w.Camera.ZoomFactor -= .1

	// 	// if w.Camera.ZoomFactor < 0.1 {
	// 	// 	w.Camera.ZoomFactor = 0.1
	// 	// }

	// 	// slog.Debug("Zooming out", "zoom", w.Camera.ZoomFactor)
	// }

	if yoff != 0 {
		g.OnMouseWheelCallbacks.Call(func(ocf OnMouseWheelFunc) (stop bool) {
			return ocf(yoff)
		})
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.Press == nil {
			g.Press = &Press{
				StartX:    x,
				StartY:    y,
				X:         x,
				Y:         y,
				TimeStart: time.Now(),
			}
		}

		// need drag callback

		if time.Since(g.Press.TimeStart) > 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-x)) > 25 || math.Abs(float64(g.Press.StartY-y)) > 25) {
			g.OnMouseDragCallbacks.Call(func(ocf OnMouseDragFunc) (stop bool) {
				return ocf(g.Press.X, g.Press.Y, x, y)
			})
			// 	w.Camera.Position[0] += float64(w.Press.X-x) / w.Camera.ZoomFactor
			// 	w.Camera.Position[1] += float64(w.Press.Y-y) / w.Camera.ZoomFactor
		}

		g.Press.X = x
		g.Press.Y = y

	} else {
		if g.Press != nil {
			if time.Since(g.Press.TimeStart) < 100*time.Millisecond || (math.Abs(float64(g.Press.StartX-g.Press.X)) < 8 && math.Abs(float64(g.Press.StartY-g.Press.Y)) < 8) {
				slog.Info("Click", "x", g.Press.X, "y", g.Press.Y)
				g.OnMouseUpCallbacks.Call(func(ocf OnMouseUpFunc) (stop bool) {
					return ocf(g.Press.X, g.Press.Y)
				})
			}
			g.Press = nil
		}
	}
}
